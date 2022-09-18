package main

import "github.com/grpc-ecosystem/grpc-gateway/runtime"

// +ioc:autowire=true
// +ioc:autowire:type=singleton
type Mux struct {
	*runtime.ServeMux
}

func New() *Mux {
	return &Mux{}
}
