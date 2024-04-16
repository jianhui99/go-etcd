package main

import (
	"context"
	"etcddemo/echo"
	"etcddemo/etcd"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	serviceName = flag.String("service", "echo-service", "service name")
)

func main() {
	flag.Parse()
	etcd.CustomLoadService(*serviceName)
	addr := etcd.CustomServiceDiscovery(*serviceName)

	log.Println("service addr:", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := echo.NewEchoClient(conn)
	CallUnaryEcho(c)
}

func CallUnaryEcho(c echo.EchoClient) {
	ctx := context.Background()
	in := echo.EchoMessage{
		Message: "client say hi",
	}
	res, err := c.UnaryEcho(ctx, &in)
	if err != nil {
		log.Fatalf("could not echo: %v", err)
	}
	fmt.Println("client received response:", res.Message)
}
