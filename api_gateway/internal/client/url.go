package client

import (
	"context"
	"log/slog"

	"api_gateway/errs"
	"api_gateway/pkg/proto/url"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name UrlClient
type UrlClient interface {
	FollowUrl(ctx context.Context, shortUrl string) (string, error)
	ShortenUrl(ctx context.Context, longUrl string) (string, error)
}

type grpcUrlClient struct {
	logger        *slog.Logger
	urlGrpcClient url.UrlClient
}

func NewGrpcUrlClient(
	logger *slog.Logger,
	urlGrpcClient url.UrlClient,

) UrlClient {
	return &grpcUrlClient{
		logger:        logger,
		urlGrpcClient: urlGrpcClient,
	}
}

func (u *grpcUrlClient) FollowUrl(ctx context.Context, shortUrl string) (string, error) {
	longURLResp, err := u.urlGrpcClient.FollowUrl(ctx, &url.ShortUrlRequest{
		ShortUrl: shortUrl,
	})

	if err != nil {
		u.logger.Error(err.Error())
		st, ok := status.FromError(err)
		if !ok || st.Code() == codes.Internal {
			return "", errs.ErrInternal
		}

		if st.Code() == codes.NotFound {
			return "", errs.ErrNotFound
		}
		if st.Code() == codes.InvalidArgument {
			return "", errs.ErrInvalidArgument
		}

		return "", errs.ErrInternal
	}

	return longURLResp.LongUrl, nil
}

func (u *grpcUrlClient) ShortenUrl(ctx context.Context, longUrl string) (string, error) {
	shortURLResp, err := u.urlGrpcClient.ShortenUrl(context.Background(), &url.LongUrlRequest{
		LongUrl: longUrl,
	})

	if err != nil {
		u.logger.Error(err.Error())
		st, ok := status.FromError(err)
		if !ok || st.Code() == codes.Internal {
			return "", errs.ErrInternal
		}
		if st.Code() == codes.InvalidArgument {
			return "", errs.ErrInvalidArgument
		}

		return "", errs.ErrInternal
	}

	return shortURLResp.ShortUrl, nil
}
