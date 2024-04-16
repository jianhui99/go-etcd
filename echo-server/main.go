package main

import (
	"etcddemo/echo"
	"etcddemo/echo-server/server"
	"etcddemo/etcd"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port        = flag.Int("port", 50051, "port to listen on")
	serviceName = flag.String("service", "echo-service", "service name")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	echo.RegisterEchoServer(grpcServer, &server.EchoServer{})
	err = etcd.CustomServiceRegister(*serviceName, fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
