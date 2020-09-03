package main

import (
	"bench-hash/hasher"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "bench-hash/server-grpc/hasher_server"
)

type hasherServer struct{}

func (h hasherServer) Hash(_ context.Context, req *pb.HashRequest) (*pb.HashResponse, error) {
	resp := &pb.HashResponse{
		Output: hasher.String(req.Input),
	}
	return resp, nil
}

func main() {
	port := 3033
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal("Unable to start TCP listener", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHasherServer(grpcServer, &hasherServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("unable to start GRPC server", err)
	}
}
