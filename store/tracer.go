package store

import (
	"io"
	"log"
)

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
		service,
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil
	}

	return closer
}
