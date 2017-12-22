package store

import (
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/garyburd/redigo/redis"
)

type Store struct {
}

type Item struct {
	Title string
	Url   string
	Ctx   opentracing.SpanContext
}

func (s *Store) storage() string {
	return os.Getenv("REDIS")
}

func (s *Store) Add(item *Item) (int64, error) {
	c, err := redis.Dial("tcp", s.storage())
	if err != nil {
		return 0, err
	}
	defer c.Close()
	c.Do("SET", item.Title, item.Url)
	return redis.Int64(c.Do("INCR", "items"))
}
