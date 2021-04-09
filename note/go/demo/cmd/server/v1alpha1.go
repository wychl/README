package main

import (
	"context"

	v1alpha1 "github.com/wychl/README/go/demo/gen/go/user/v1alpha1"
)

type v1alpha1Server struct{}

func newV1alpha1Server() *v1alpha1Server {
	return &v1alpha1Server{}
}

func (s *v1alpha1Server) Create(ctx context.Context, request *v1alpha1.CreateRequest) (*v1alpha1.CreateResponse, error) {
	return &v1alpha1.CreateResponse{
		Name:    request.GetName(),
		Age:     request.GetAge(),
		Gender:  request.GetGender(),
		Version: "v1alpha1",
	}, nil
}
