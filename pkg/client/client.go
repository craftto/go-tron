package client

import (
	"context"
	"time"

	"github.com/craftto/go-tron/pkg/proto/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GrpcClient struct {
	GrpcURL     string
	Conn        *grpc.ClientConn
	Client      api.WalletClient
	grpcTimeout time.Duration
	opts        []grpc.DialOption
	apiKey      string
}

func NewGrpcClient(grpcURL string, opts ...grpc.DialOption) (*GrpcClient, error) {
	conn, err := grpc.Dial(grpcURL, opts...)
	if err != nil {
		return nil, err
	}

	return &GrpcClient{
		GrpcURL:     grpcURL,
		Conn:        conn,
		Client:      api.NewWalletClient(conn),
		grpcTimeout: 5 * time.Second,
		opts:        opts,
	}, nil
}

func (g *GrpcClient) SetAPIKey(apiKey string) error {
	g.apiKey = apiKey
	return nil
}

func (g *GrpcClient) GetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), g.grpcTimeout)

	if len(g.apiKey) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", g.apiKey)
	}

	return ctx, cancel
}

func (g *GrpcClient) Close() {
	if g.Conn != nil {
		g.Conn.Close()
	}
}

// GetMessageBytes return grpc message from bytes
func GetMessageBytes(m []byte) *api.BytesMessage {
	return &api.BytesMessage{
		Value: m,
	}
}

// GetPaginatedMessage return grpc message number
func GetPaginatedMessage(offset int64, limit int64) *api.PaginatedMessage {
	return &api.PaginatedMessage{
		Offset: offset,
		Limit:  limit,
	}
}
