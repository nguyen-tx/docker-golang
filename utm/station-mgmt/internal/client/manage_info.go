package client

import (
	"context"
	"fmt"

	mipb "github.com/utm/station-mgmt/pkg/pb/manage_info/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ManageInfoClient struct {
	conn   *grpc.ClientConn
	client mipb.ManageInfoServiceClient
}

func NewManageInfoClient(addr string) (*ManageInfoClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("manage info grpc dial: %w", err)
	}
	return &ManageInfoClient{
		conn:   conn,
		client: mipb.NewManageInfoServiceClient(conn),
	}, nil
}

func (c *ManageInfoClient) GetSummary(ctx context.Context) (*mipb.ManageInfoSummary, error) {
	return c.client.GetSummary(ctx, &mipb.GetSummaryRequest{})
}

func (c *ManageInfoClient) Close() { c.conn.Close() }
