package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	v1alpha1 "github.com/wychl/README/go/demo/gen/go/user/v1alpha1"
	v1beta1 "github.com/wychl/README/go/demo/gen/go/user/v1beta1"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	v1alpha1Client(conn)
	v1betav1Client(conn)
}

func v1alpha1Client(conn *grpc.ClientConn) {
	clietn := v1alpha1.NewUserApiClient(conn)
	resp, err := clietn.Create(context.Background(), &v1alpha1.CreateRequest{
		Name:   "alice",
		Age:    18,
		Gender: v1alpha1.Gender_GENDER_FEMAN,
	})
	if err != nil {
		log.Fatalf("create user err:%v", err)
		return
	}

	data, _ := json.Marshal(resp)
	fmt.Println(string(data))

}

func v1betav1Client(conn *grpc.ClientConn) {
	clietnV1beta1 := v1beta1.NewUserApiClient(conn)
	respV1beta1, err := clietnV1beta1.Create(context.Background(), &v1beta1.CreateRequest{
		Name:   "alice",
		Age:    18,
		Gender: v1beta1.Gender_GENDER_FEMAN,
	})
	if err != nil {
		log.Fatalf("create user err:%v", err)
		return
	}

	data, _ := json.Marshal(respV1beta1)
	fmt.Println(string(data))
}
