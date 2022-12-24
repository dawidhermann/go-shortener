package shortener

import (
	"context"
	"github.com/dawidhermann/shortener-url/internal/db/db"
	pb "github.com/dawidhermann/shortener-url/internal/protobuf"
)

type Server struct {
	pb.UnimplementedShortenerClientServer
}

func (s *Server) ShortenUrl(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	shortenedUrl := shortenUrl()
	err := db.SaveUrl(shortenedUrl, req.Url)
	return &pb.ShortenResponse{Id: shortenedUrl}, err
}
