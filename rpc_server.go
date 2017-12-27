package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Do struct{}

func (r *Do) Deed(item *Item, num *int64) error {

	tracer, closer := TracerInit("S" + strconv.Itoa(rand.Intn(2)+item.Level))
	if closer != nil {
		defer closer.Close()
	}

	var span opentracing.Span
	callerContext, err := tracer.Extract(opentracing.TextMap, item.Trace)
	if callerContext != nil {
		span = tracer.StartSpan("Deed", ext.RPCServerOption(callerContext))
	} else {
		span = tracer.StartSpan("Deed")
	}
	defer span.Finish()

	if time.Now().UnixNano()%2 == 0 {
		fmt.Println("Nested")
		item.Level++
		*num, err = nextCall(tracer, span, *item)
		return err
	}

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	s := Store{}
	*num, err = s.Add(tracer, ctx, item)
	return err
}

func serveRPC() {
	do := new(Do)
	rpc.Register(do)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	rpc.Accept(listener)

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
