package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"

	"github.com/kavehmz/opentracer-example/store"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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
	parent := opentracing.StartSpan("Add Item")
	defer parent.Finish()

	client, err := conn(parent, "ADDSRV")
	defer client.Close()

	item := store.Item{Title: "test", Url: "url"}
	var reply int64
	callSpan := opentracing.StartSpan("call", opentracing.ChildOf(parent.Context()))
	defer callSpan.Finish()

	item.Trace = make(map[string]string)
	opentracing.GlobalTracer().Inject(callSpan.Context(), opentracing.TextMap, item.Trace)

	err = client.Call("Add.Add", item, &reply)
	if err != nil {
		log.Panic("add error:", err)
	}
	fmt.Fprintf(w, "added item number: %d\n", reply)
}

func rm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "rm!")
}

func ls(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ls!")
}

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get!")
}

func tracerInit(service string) io.Closer {
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
		service,
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil
	}

	return closer
}

func main() {
	defer tracerInit("Store").Close()

	http.HandleFunc("/add", add)
	http.HandleFunc("/rm", rm)
	http.HandleFunc("/ls", ls)
	http.HandleFunc("/get", get)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
