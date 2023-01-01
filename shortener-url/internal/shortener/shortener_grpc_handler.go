package shortener

import (
	"context"
	"github.com/dawidhermann/shortener-url/internal/db/db"
	pb "github.com/dawidhermann/shortener-url/internal/protobuf"
	"log"
)

type Server struct {
	pb.UnimplementedShortenerClientServer
}

func (s *Server) CreateShortenedUrl(ctx context.Context, req *pb.CreateShortenedUrlRequest) (*pb.CreateShortenedUrlResponse, error) {
	shortenedUrl := shortenUrl()
	err := db.SaveUrl(shortenedUrl, req.Url)
	if err != nil {
		return &pb.CreateShortenedUrlResponse{}, err
	}
	log.Printf("Created url: %v", shortenedUrl)
	return &pb.CreateShortenedUrlResponse{Id: shortenedUrl}, nil
}

func (s *Server) DeleteShortenedUrl(ctx context.Context, req *pb.DeleteShortenedUrlRequest) (*pb.DeleteShortenedUrlResponse, error) {
	err := db.DeleteUrl(req.Key)
	if err != nil {
		return &pb.DeleteShortenedUrlResponse{}, err
	}
	log.Printf("Deleted url: %v", req.Key)
	return &pb.DeleteShortenedUrlResponse{Key: req.Key}, nil
}
