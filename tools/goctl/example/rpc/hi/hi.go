package main

import (
	"flag"
	"fmt"

	"go-kit/tools/goctl/example/rpc/hi/internal/config"
	eventServer "go-kit/tools/goctl/example/rpc/hi/internal/server/event"
	greetServer "go-kit/tools/goctl/example/rpc/hi/internal/server/greet"
	"go-kit/tools/goctl/example/rpc/hi/internal/svc"
	"go-kit/tools/goctl/example/rpc/hi/pb/hi"

	"go-kit/core/conf"
	"go-kit/core/service"
	"go-kit/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/hi.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		hi.RegisterGreetServer(grpcServer, greetServer.NewGreetServer(ctx))
		hi.RegisterEventServer(grpcServer, eventServer.NewEventServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
