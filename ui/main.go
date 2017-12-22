package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/kavehmz/opentracer-example/store"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func conn(parent opentracing.Span, service string) (*rpc.Client, error) {
	sp := opentracing.StartSpan("conncet to rpc server", opentracing.ChildOf(parent.Context()))
	defer sp.Finish()

	client, err := rpc.Dial("tcp", os.Getenv(service))
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

	args := store.Item{Title: "test", Url: "url", Ctx: parent.Context()}
	var reply int64
	callSpan := opentracing.StartSpan("call", opentracing.ChildOf(parent.Context()))
	defer callSpan.Finish()
	err = (client).Call("Add.Add", args, &reply)
	if err != nil {
		log.Panic("arith error:", err)
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
	fmt.Fprintf(w, "ls!")
}

func tracerInit(serviceName string) (opentracing.Tracer, io.Closer) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, closer, err := cfg.New(
		serviceName,
	)
	if err != nil {

	}
	return tracer, closer
}

func main() {
	tracer, closer := tracerInit("prime-example")
	defer closer.Close()
	opentracing.InitGlobalTracer(tracer)

	http.HandleFunc("/add", add)
	http.HandleFunc("/rm", rm)
	http.HandleFunc("/ls", ls)
	http.HandleFunc("/get", get)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
