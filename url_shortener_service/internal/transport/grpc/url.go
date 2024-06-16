package grpc

import (
	"context"
	"errors"
	"log/slog"

	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/service"
	url "CoolUrlShortener/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UrlServer struct {
	logger     *slog.Logger
	urlService service.URLService
	url.UnimplementedUrlServer
}

func NewUrlServer(
	logger *slog.Logger,
	urlService service.URLService,
) *UrlServer {
	return &UrlServer{
		logger:     logger,
		urlService: urlService,
	}
}

func (s *UrlServer) ShortenUrl(ctx context.Context, req *url.LongUrlRequest) (*url.UrlDataResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	shortURL, err := s.urlService.SaveURL(ctx, req.LongUrl)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &url.UrlDataResponse{
		LongUrl:  req.LongUrl,
		ShortUrl: shortURL,
	}, nil
}

func (s *UrlServer) FollowUrl(ctx context.Context, req *url.ShortUrlRequest) (*url.LongUrlResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	longUrl, err := s.urlService.GetLongURL(ctx, req.ShortUrl)
	if err != nil {
		s.logger.Error(err.Error())
		if errors.Is(err, errs.ErrNoURL) {
			return nil, status.Error(codes.NotFound, "short url not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &url.LongUrlResponse{
		LongUrl: longUrl,
	}, nil
}
