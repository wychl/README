package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"

	"github.com/wychl/README/jaeger/demo/lib/tracing"
)

func main() {
	tracer, closer := tracing.Init("service1")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	http.HandleFunc("/service1", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("service1-opname", ext.RPCServerOption(spanCtx))
		defer span.Finish()
		helloStr := "hello service1"
		span.LogFields(
			log.String("event", "service1-handle"),
			log.String("value", helloStr),
		)
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		//time.Sleep(time.Duration(2)*time.Second)
		CallService2(ctx)
		w.Write([]byte(helloStr))
	})
	panic(http.ListenAndServe(":10081", nil))
}

func CallService2(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "callservice2-opname")
	defer span.Finish()
	url := "http://localhost:10082/service2"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	client := http.Client{}
	//span.LogEvent("Before call service2")
	//time.Sleep(time.Duration(2)*time.Second)
	resp, err := client.Do(req)
	//span.LogEvent("After call service2")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	respStr := string(body)
	span.LogFields(
		log.String("event", "service2-call"),
		log.String("value", respStr),
	)
	fmt.Println(respStr)
}
