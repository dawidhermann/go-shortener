package shortener

import (
	"context"
	"log"

	"github.com/dawidhermann/shortener-url/internal/db/db"
	pb "github.com/dawidhermann/shortener-url/internal/protobuf"
)

type Server struct {
	pb.UnimplementedShortenerClientServer
	kvStore *db.KVStore
}

func NewServer(kvStore *db.KVStore) *Server {
	return &Server{
		kvStore: kvStore,
	}
}

// Creates shortened url and saves it in kv store
func (s *Server) CreateShortenedUrl(ctx context.Context, req *pb.CreateShortenedUrlRequest) (*pb.CreateShortenedUrlResponse, error) {
	shortenedUrl := shortenUrl()
	err := s.kvStore.SaveUrl(shortenedUrl, req.Url)
	if err != nil {
		return &pb.CreateShortenedUrlResponse{}, err
	}
	log.Printf("Created url: %v", shortenedUrl)
	return &pb.CreateShortenedUrlResponse{Id: shortenedUrl}, nil
}

// Deletes shortened url from kv store
func (s *Server) DeleteShortenedUrl(ctx context.Context, req *pb.DeleteShortenedUrlRequest) (*pb.DeleteShortenedUrlResponse, error) {
	err := s.kvStore.DeleteUrl(req.Key)
	if err != nil {
		return &pb.DeleteShortenedUrlResponse{}, err
	}
	log.Printf("Deleted url: %v", req.Key)
	return &pb.DeleteShortenedUrlResponse{Key: req.Key}, nil
}
