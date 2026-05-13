package client

import (
	"context"
	"fmt"

	rtkpb "github.com/utm/station-mgmt/pkg/pb/rtk_station/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RTKStationClient struct {
	conn   *grpc.ClientConn
	client rtkpb.RTKStationServiceClient
}

func NewRTKStationClient(addr string) (*RTKStationClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("rtk station grpc dial: %w", err)
	}
	return &RTKStationClient{
		conn:   conn,
		client: rtkpb.NewRTKStationServiceClient(conn),
	}, nil
}

func (c *RTKStationClient) GetSummary(ctx context.Context) (*rtkpb.StationSummary, error) {
	return c.client.GetSummary(ctx, &rtkpb.GetSummaryRequest{})
}

func (c *RTKStationClient) Close() { c.conn.Close() }
