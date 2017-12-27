package main

import (
	"log"
	"net/rpc"

	opentracing "github.com/opentracing/opentracing-go"
)

func nextCall(tracer opentracing.Tracer, parent opentracing.Span) int64 {
	client, err := conn(tracer, parent, "RPC")
	defer client.Close()

	span := tracer.StartSpan("Call RPC", opentracing.ChildOf(parent.Context()))
	defer span.Finish()

	var reply int64
	item := Item{Title: "title", Url: "url", Trace: make(map[string]string)}
	tracer.Inject(span.Context(), opentracing.TextMap, item.Trace)

	err = client.Call("Do.Deed", item, &reply)
	if err != nil {
		log.Panic("Do.Deed error:", err)
	}
	return reply
}

func conn(tracer opentracing.Tracer, parent opentracing.Span, service string) (*rpc.Client, error) {
	sp := opentracing.StartSpan("Conncet to rpc server", opentracing.ChildOf(parent.Context()))
	defer sp.Finish()

	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Panic("dialing service:", err)
	}

	return client, err
}
