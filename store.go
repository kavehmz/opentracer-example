package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	opentracing "github.com/opentracing/opentracing-go"
)

type Store struct {
}

type Item struct {
	Title string
	Url   string
	Trace TextMapCarrier
	Level int
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

func (s *Store) Add(tracer opentracing.Tracer, ctx context.Context, item *Item) (int64, error) {
	parent := opentracing.SpanFromContext(ctx)
	if parent == nil {
		fmt.Println("Span not found")
	} else {
		sp := tracer.StartSpan("Save To Redis", opentracing.ChildOf(parent.Context()))
		defer sp.Finish()
	}
	if os.Getenv("REDISURL") != "" {

		c, err := redis.Dial("tcp", os.Getenv("REDISURL"))
		if err != nil {
			return 0, err
		}
		defer c.Close()
		c.Do("SET", item.Title, item.Url)
		return redis.Int64(c.Do("INCR", "items"))
	}

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(3000)))
	return time.Now().Unix(), nil
}
