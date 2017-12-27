package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/kavehmz/opentracer-example/store"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Add struct{}

func (r *Add) Add(item *store.Item, num *int64) error {
	var span opentracing.Span
	callerContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, item.Trace)
	if callerContext != nil {
		span = opentracing.StartSpan("add operation", ext.RPCServerOption(callerContext))
	} else {
		span = opentracing.StartSpan("add operation")
	}
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	s := store.Store{}
	*num, err = s.Add(ctx, item)
	return err
}

func main() {

	defer store.TracerInit("Add Service").Close()

	add := new(Add)
	rpc.Register(add)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}
