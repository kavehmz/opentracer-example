# Test run
```bash
$ go test -bench=. opentracing_test.go 
goos: darwin
goarch: amd64
BenchmarkStartSpan-4   	100000000	        17.8 ns/op
BenchmarkDial-4        	    1000	   1444179 ns/op
BenchmarkSprintf-4     	20000000	       108 ns/op
PASS
ok  	command-line-arguments	5.699s
```