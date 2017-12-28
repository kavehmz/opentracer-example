# opentracer-example
Example for travelaudience usage. This is to serve only as an example for how we might use opentracing as an intrumentation standard, Ubers' Jaeger for providing the service along our RPC calls.

# Dependencies

```bash
$ dep ensure -v
```

# Optional Redis

This test will issue redis `SET` is `REDISURL` environment variable is set. Otherwise it will sleep randomly for about 3ms.

# Run
```bash
go run *.go
```

# References

Uber's blog about Jaeger and opentracing:

https://medium.com/opentracing/take-opentracing-for-a-hotrod-ride-f6e3141f7941

Go opentracing API and how to use it:

https://github.com/opentracing/opentracing-go

Notice: opentracing itself and what we see in the code is independent of implemetation. We only need to know the API. Jaeger or Zipkin are just implementation details which will only change some details at Initilization time.
