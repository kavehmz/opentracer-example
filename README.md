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
