package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	//从环境变量读取端口密码
	port := os.Getenv("PORT")
	fmt.Println("port:", port)

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":8888"), nil))
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}
