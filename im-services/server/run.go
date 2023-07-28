package server

import (
	"im-services/internal/config"
	grpcMessage "im-services/server/grpc/message"
	"log"
	"net"

	"google.golang.org/grpc"
)

var RpcServer = grpc.NewServer()

func StartGrpc() {
	if config.Conf.Server.ServiceOpen {

		var message grpcMessage.ImGrpcMessage

		grpcMessage.RegisterImMessageServer(RpcServer, message)

		listener, err := net.Listen("tcp", ":8002")
		if err != nil {
			log.Fatal("grpc服务启动失败", err)
		}
		_ = RpcServer.Serve(listener)
	}
}
