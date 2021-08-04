package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/golog"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "172.16.100.130:4000",
	})

	abf := client.Get(context.Background(), "abf")
	golog.Errorf("Abf: %+v", abf)

	info := client.Info(context.Background())

	nodes := client.ClusterNodes(context.Background())
	args := nodes.Args()
	golog.Errorf("Nodes: %+v, \n%+v, Info:%+v", nodes, args, info)
}
