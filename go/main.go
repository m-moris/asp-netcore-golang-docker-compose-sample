package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/m-moris/asp-netcore-golang-docker-compose-sample/go/lib/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port int = 50051

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	name := "Empty"
	if in.GetName() != "" {
		name = in.GetName()
	}
	log.Printf("Received: %v", name)
	return &pb.HelloReply{
		Message: "Hello " + name + "-san",
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func main() {

	if xport, err := strconv.Atoi(os.Getenv("GRPC_PORT")); err == nil {
		port = xport
	}

	addr := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRpc Server start : port = %v", port)

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
