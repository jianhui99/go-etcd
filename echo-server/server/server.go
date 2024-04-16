package server

import (
	"context"
	"etcddemo/echo"
	"fmt"
)

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func (EchoServer) UnaryEcho(ctx context.Context, m *echo.EchoMessage) (*echo.EchoMessage, error) {
	fmt.Println("server received in: ", m)
	return &echo.EchoMessage{
		Message: "hello client",
	}, nil
}
