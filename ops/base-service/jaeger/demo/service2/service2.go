package main

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/wychl/README/jaeger/demo/lib/tracing"
)

func main() {
	tracer, closer := tracing.Init("service2")
	defer closer.Close()
	http.HandleFunc("/service2", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("service2-opname", ext.RPCServerOption(spanCtx))
		defer span.Finish()
		helloStr := "hello service2"
		println(helloStr)
		w.Write([]byte(helloStr))
	})
	panic(http.ListenAndServe(":10082", nil))
}
