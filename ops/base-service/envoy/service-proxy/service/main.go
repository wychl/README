package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.Handle("/service/", http.HandlerFunc(service))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func service(w http.ResponseWriter, r *http.Request) {
	service := strings.TrimPrefix(r.URL.Path, "/service/")
	output := fmt.Sprintf("Hello from behind Envoy (service %s)! hostname %s \n", service, os.Getenv("HOSTNAME"))
	fmt.Fprint(w, output)
}
