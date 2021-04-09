package main

import (
	"context"

	v1beta1 "github.com/wychl/README/go/demo/gen/go/user/v1beta1"
)

type v1beta1Server struct{}

func newV1beta1Server() *v1beta1Server {
	return &v1beta1Server{}
}

func (s *v1beta1Server) Create(ctx context.Context, request *v1beta1.CreateRequest) (*v1beta1.CreateResponse, error) {
	return &v1beta1.CreateResponse{
		Name:    request.GetName(),
		Age:     request.GetAge(),
		Gender:  request.GetGender(),
		Version: "v1beta1",
	}, nil
}
