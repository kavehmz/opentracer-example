package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"

	"github.com/kavehmz/opentracer-example/store"
	opentracing "github.com/opentracing/opentracing-go"
)

func conn(parent opentracing.Span, service string) (*rpc.Client, error) {
	sp := opentracing.StartSpan("conncet to rpc server", opentracing.ChildOf(parent.Context()))
	defer sp.Finish()

	client, err := rpc.DialHTTP("tcp", os.Getenv(service))
	if err != nil {
		log.Panic("dialing add service:", err)
	}

	return client, err
}

func add(w http.ResponseWriter, r *http.Request) {
	parent := opentracing.StartSpan("Frontend Add Item")
	defer parent.Finish()

	client, err := conn(parent, "ADDSRV")
	defer client.Close()

	callSpan := opentracing.StartSpan("call", opentracing.ChildOf(parent.Context()))
	defer callSpan.Finish()

	var reply int64
	item := store.Item{Title: "test", Url: "url", Trace: make(map[string]string)}
	opentracing.GlobalTracer().Inject(callSpan.Context(), opentracing.TextMap, item.Trace)

	err = client.Call("Add.Add", item, &reply)
	if err != nil {
		log.Panic("add error:", err)
	}
	fmt.Fprintf(w, "added item number: %d\n", reply)
}

func rm(w http.ResponseWriter, r *http.Request) {
	parent := opentracing.StartSpan("Frontend Remove Item")
	defer parent.Finish()

	client, err := conn(parent, "REMOVESRV")
	defer client.Close()

	callSpan := opentracing.StartSpan("call", opentracing.ChildOf(parent.Context()))
	defer callSpan.Finish()

	var reply int64
	item := store.Item{Title: "test", Url: "url", Trace: make(map[string]string)}
	opentracing.GlobalTracer().Inject(callSpan.Context(), opentracing.TextMap, item.Trace)

	err = client.Call("Remove.Remove", item, &reply)
	if err != nil {
		log.Panic("remove error:", err)
	}
	fmt.Fprintf(w, "remove item number: %d\n", reply)
}

func main() {
	defer store.TracerInit("Store Front").Close()

	http.HandleFunc("/add", add)
	http.HandleFunc("/remove", rm)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
