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
	tracer, closer := tracing.Init("hello")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("hello-opname", ext.RPCServerOption(spanCtx))
		span.SetTag("hello-tag-key", "hello-tag-value")
		defer span.Finish()

		helloStr := "hello Jaeger"
		span.LogFields(
			log.String("event", "hello-handle"),
			log.String("value", helloStr),
		)
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		CallService1(ctx)
		CallService2(ctx)
		w.Write([]byte(helloStr))
	})
	panic(http.ListenAndServe(":10080", nil))
}

func CallService1(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "callservice1-opname")
	defer span.Finish()
	url := "http://localhost:10081/service1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	respStr := string(body)
	span.LogFields(
		log.String("event", "service1-call"),
		log.String("value", respStr),
	)
	fmt.Println(respStr)
}

func CallService2(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "callservice2-opname")
	defer span.Finish()
	url := "http://localhost:10082/service2"
	//for i:=0; i<32; i++ {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	client := http.Client{}
	resp, err := client.Do(req)
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
	//}
}
