package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/kavehmz/opentracer-example/store"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Add struct {
}

func (r *Add) Add(item *store.Item, num *int64) error {
	var serverSpan opentracing.Span
	wireContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, item.Trace)
	if err != nil {
		panic(err)
	}

	serverSpan = opentracing.StartSpan("add operation", ext.RPCServerOption(wireContext))
	defer serverSpan.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), serverSpan)

	s := store.Store{}
	*num, err = s.Add(ctx, item)
	return err
}

func tracerInit() io.Closer {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	closer, err := cfg.InitGlobalTracer(
		"Add Service",
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil
	}

	return closer
}

func main() {
	defer tracerInit().Close()

	add := new(Add)
	rpc.Register(add)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}

func checkError(err error) {
	if err != nil {
		log.Panic("Fatal error ", err.Error())
	}
}
