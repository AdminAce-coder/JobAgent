package main

import (
	"flag"
	"fmt"

	"github.com/AdminAce-coder/jobAgent/internal/config"
	"github.com/AdminAce-coder/jobAgent/internal/server"
	"github.com/AdminAce-coder/jobAgent/internal/svc"
	"github.com/AdminAce-coder/jobAgent/pb/jobAgent"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/jobagent.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		jobAgent.RegisterJobAgentServer(grpcServer, server.NewJobAgentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
