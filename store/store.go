package store

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Store struct {
}

type Item struct {
	Title string
	Url   string
	Trace TextMapCarrier
}

type TextMapCarrier map[string]string

func (c TextMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range c {
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (c TextMapCarrier) Set(key, val string) {
	fmt.Println("SET", key, val)
	c[key] = val
}

func (s *Store) storage() string {
	return os.Getenv("REDIS")
}

func (s *Store) Add(item *Item) (int64, error) {
	return time.Now().Unix(), nil
	// c, err := redis.Dial("tcp", s.storage())
	// if err != nil {
	// 	return 0, err
	// }
	// defer c.Close()
	// c.Do("SET", item.Title, item.Url)
	// return redis.Int64(c.Do("INCR", "items"))
}

func TracerInit(service string) io.Closer {
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
		"service",
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil
	}

	return closer
}