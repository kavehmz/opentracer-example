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

type Remove struct{}

func (r *Remove) Remove(item *store.Item, num *int64) error {
	var span opentracing.Span
	callerContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, item.Trace)
	if callerContext != nil {
		span = opentracing.StartSpan("remove operation", ext.RPCServerOption(callerContext))
	} else {
		span = opentracing.StartSpan("remove operation")
	}
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	s := store.Store{}
	*num, err = s.Remove(ctx, item)
	return err
}

func main() {
	defer store.TracerInit("Remove Service").Close()

	remove := new(Remove)
	rpc.Register(remove)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1235")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}
