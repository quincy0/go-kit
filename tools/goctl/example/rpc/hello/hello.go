package main

import (
	"flag"
	"fmt"

	"go-kit/tools/goctl/example/rpc/hello/internal/config"
	greetServer "go-kit/tools/goctl/example/rpc/hello/internal/server/greet"
	"go-kit/tools/goctl/example/rpc/hello/internal/svc"
	"go-kit/tools/goctl/example/rpc/hello/pb/hello"

	"go-kit/core/conf"
	"go-kit/core/service"
	"go-kit/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/hello.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		hello.RegisterGreetServer(grpcServer, greetServer.NewGreetServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
