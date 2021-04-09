package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/argoproj/argo-cd/util/swagger"
)

var (
	swaggerPath = flag.String("swagger", "", "swagger json file path")
	httpPort    = flag.Int("port", 9999, "http port")
)

func main() {
	flag.Parse()

	if *swaggerPath == "" {
		log.Fatalf("swagger file path not set\n")
		return
	}

	data, err := ioutil.ReadFile(*swaggerPath)
	if err != nil {
		log.Fatalf("read swagger file error:%v\n", err)
		return
	}

	r := http.NewServeMux()
	swagger.ServeSwaggerUI(r, string(data), "/")
	http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), r)
}
