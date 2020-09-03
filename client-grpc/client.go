package client_grpc

import (
	pb "bench-hash/server-grpc/hasher_server"
	"context"
	"google.golang.org/grpc"
)

type HashClient struct {
	conn   *grpc.ClientConn
	client pb.HasherClient
}

func NewHashClient(addr string) (*HashClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &HashClient{
		conn:   conn,
		client: pb.NewHasherClient(conn),
	}, nil
}

func (c HashClient) Hash(input string) (string, error) {
	resp, err := c.client.Hash(context.TODO(), &pb.HashRequest{Input: input})
	if err != nil {
		return "", err
	}

	return resp.Output, nil
}

func (c HashClient) Close() {
	c.conn.Close()
}
