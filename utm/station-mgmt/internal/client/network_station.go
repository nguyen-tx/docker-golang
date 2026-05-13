package client

import (
	"context"
	"fmt"

	nspb "github.com/utm/station-mgmt/pkg/pb/network_station/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NetworkStationClient struct {
	conn   *grpc.ClientConn
	client nspb.NetworkStationServiceClient
}

func NewNetworkStationClient(addr string) (*NetworkStationClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("network station grpc dial: %w", err)
	}
	return &NetworkStationClient{
		conn:   conn,
		client: nspb.NewNetworkStationServiceClient(conn),
	}, nil
}

func (c *NetworkStationClient) GetSummary(ctx context.Context) (*nspb.StationSummary, error) {
	return c.client.GetSummary(ctx, &nspb.GetSummaryRequest{})
}

func (c *NetworkStationClient) Close() { c.conn.Close() }
