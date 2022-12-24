package main

import (
	"fmt"
	"github.com/dawidhermann/shortener-url/internal/shortener"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

import pb "github.com/dawidhermann/shortener-url/internal/protobuf"

func main() {
	//url := shortener.ShortenUrl()
	//fmt.Println(url)
	//api.StartServer()
	//dbInstance := db.NewConnection("shortener_user", "P0sTgr3sP4SS", "postgres:5432", "shortener_db")
	//db.CreateUser("dhermann", "secretpass", "email@example.com", dbInstance)
	port := os.Getenv("SERVICE_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
