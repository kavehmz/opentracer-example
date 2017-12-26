# Test run

go test -bench=. opentracing_test.go 
goos: darwin
goarch: amd64
BenchmarkOpenTracing-4   	100000000	        17.8 ns/op
BenchmarkDial-4          	    1000	   1304226 ns/op
BenchmarkSprintf-4       	20000000	       106 ns/op
PASS
ok  	command-line-arguments	5.559s