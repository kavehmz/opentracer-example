package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Do struct{}

func (r *Do) Deed(item *Item, num *int64) error {
	var span opentracing.Span
	callerContext, err := rpcTracer.Extract(opentracing.TextMap, item.Trace)
	if callerContext != nil {
		span = rpcTracer.StartSpan("Deed", ext.RPCServerOption(callerContext))
	} else {
		span = rpcTracer.StartSpan("Deed")
	}
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	s := Store{}
	*num, err = s.Add(ctx, item)
	return err
}

var rpcTracer opentracing.Tracer

func serveRPC() {
	var closer io.Closer
	rpcTracer, closer = TracerInit("Do Service")
	if closer != nil {
		defer closer.Close()
	}

	do := new(Do)
	rpc.Register(do)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
