package main

import (
	"fmt"
	pb "github.com/dawidhermann/shortener-url/internal/protobuf"
	"github.com/dawidhermann/shortener-url/internal/shortener"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_SERVER_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterShortenerClientServer(grpcServer, &shortener.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
