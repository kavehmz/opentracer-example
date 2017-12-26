package main

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	opentracing "github.com/opentracing/opentracing-go"
)

func BenchmarkOpenTracing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		span := opentracing.StartSpan("add operation")
		span.LogEvent("test event")
		span.Finish()
	}
}

func BenchmarkDial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c, err := redis.Dial("tcp", "localhost:6379")
		if err == nil {
			c.Close()
		}
	}
}

func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if x := fmt.Sprintf("%d", 42); x != "42" {
			b.Fatalf("Unexpected string: %s", x)
		}
	}
}
