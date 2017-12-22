package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/opentracing/opentracing-go"

	"github.com/kavehmz/opentracer-example/store"
)

type Add struct {
	spanContext opentracing.SpanContext
}

func (r *Add) Add(item *store.Item, num *int64) error {
	s := store.Store{}
	var err error
	*num, err = s.Add(item)
	return err
}

func main() {
	add := new(Add)
	rpc.Register(add)

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
		log.Panic("Fatal error ", err.Error())
	}
}
