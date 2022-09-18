package main

import (
	"fmt"
	"github.com/alibaba/ioc-golang"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
type App struct {
	Mux *Mux `singleton:"main.Mux"`
	// Ttp *grpc.Server `singleton:"grpc.Server"`
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton
type OtherMux struct {
	Mux *Mux `singleton:"main.Mux"`
}

func main() {
	// 加载所有结构
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	// 获取结构
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}

	other, err := GetOtherMuxSingleton()
	if err != nil {
		panic(err)
	}

	fmt.Printf("App mux: %T, %v, Other mux: %T, %v", app.Mux, app.Mux, other.Mux, other.Mux)
}
