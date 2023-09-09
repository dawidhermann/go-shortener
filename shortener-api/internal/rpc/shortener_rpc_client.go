// rpc client
package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/dawidhermann/shortener-api/internal/config"
	pb "github.com/dawidhermann/shortener-api/internal/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnRpc struct {
	conn *grpc.ClientConn
}

// Connect to rpc service
func Connect(rpcConfig config.GrpcConfig) ConnRpc {
	grpcServerAddr := fmt.Sprintf("%s:%d", rpcConfig.GrpcServerHost, rpcConfig.GrpcServerPort)
	connPtr, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return ConnRpc{conn: connPtr}
}

// close rpc connection
func (rpcConn ConnRpc) Close() {
	rpcConn.Close()
}

// create shortened url using rpc connection
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

// delete shortened url using rpc connection
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
