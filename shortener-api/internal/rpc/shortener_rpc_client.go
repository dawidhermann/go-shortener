package rpc

import (
	"context"
	"fmt"
	"github.com/dawidhermann/shortener-api/config"
	pb "github.com/dawidhermann/shortener-api/internal/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ConnRpc struct {
	conn *grpc.ClientConn
}

func Connect(rpcConfig config.RpcConfig) ConnRpc {
	grpcServerAddr := fmt.Sprintf("%s:%s", rpcConfig.ServerHost, rpcConfig.ServerPort)
	connPtr, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return ConnRpc{conn: connPtr}
}

func (rpcConn ConnRpc) Close() {
	rpcConn.Close()
}

func (rpcConn ConnRpc) CreateShortenUrl(targetUrl string) (string, error) {
	client := pb.NewShortenerClientClient(rpcConn.conn)
	resp, err := client.CreateShortenedUrl(context.Background(), &pb.CreateShortenedUrlRequest{
		Url: targetUrl,
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (rpcConn ConnRpc) DeleteShortenedUrl(urlKey string) error {
	client := pb.NewShortenerClientClient(rpcConn.conn)
	_, err := client.DeleteShortenedUrl(context.Background(), &pb.DeleteShortenedUrlRequest{
		Key: urlKey,
	})
	if err != nil {
		return err
	}
	return nil
}
