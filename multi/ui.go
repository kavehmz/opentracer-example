package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
)

var uiTracer opentracing.Tracer

func serve(w http.ResponseWriter, r *http.Request) {
	span := uiTracer.StartSpan("Frontend Service")
	defer span.Finish()

	reply := nextCall(uiTracer, span)

	fmt.Fprintf(w, "served: %d\n", reply)
}

func serveUI() {
	var closer io.Closer
	uiTracer, closer = TracerInit("Store Front")
	if closer != nil {
		defer closer.Close()
	}

	http.HandleFunc("/serve", serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
