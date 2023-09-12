package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dawidhermann/shortener-url/internal/config"
	"github.com/dawidhermann/shortener-url/internal/db/db"
	pb "github.com/dawidhermann/shortener-url/internal/protobuf"
	"github.com/dawidhermann/shortener-url/internal/shortener"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.GetAppConfiguration()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	kvStore := db.New(cfg.Store)
	grpcServer := grpc.NewServer()
	pb.RegisterShortenerClientServer(grpcServer, shortener.NewServer(kvStore))
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
