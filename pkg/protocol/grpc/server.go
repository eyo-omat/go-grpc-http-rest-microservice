package grpc

import (
	"log"
	"os/signal"
	"os"
	"net"
	"context"
	
	"google.golang.org/grpc"


	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/api/v1"
)

// RunServer runs the gRPC service to publich the ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	//register service
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig  is a ^C, handle it
			log.Println("shutting down gRPC server...")
			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}

