package scrp

import (
	"context"
	"fmt"

	scrppb "github.com/utm/backend/pkg/pb/scrp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client scrppb.SCRPServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("scrp grpc dial: %w", err)
	}
	return &Client{
		conn:   conn,
		client: scrppb.NewSCRPServiceClient(conn),
	}, nil
}

func (c *Client) ResolveConflict(ctx context.Context, req *scrppb.ResolveConflictRequest) (*scrppb.ResolveConflictResponse, error) {
	return c.client.ResolveConflict(ctx, req)
}

func (c *Client) Close() { c.conn.Close() }
