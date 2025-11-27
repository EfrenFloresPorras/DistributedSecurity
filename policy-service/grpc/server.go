package grpcserver

import (
	"context"
	"log"
	"net"

	policypb "DistributedSecurity/pkg/proto/policypb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PolicyServer struct {
	policypb.UnimplementedPolicyServiceServer
}

func (s *PolicyServer) CheckPolicy(ctx context.Context, req *policypb.CheckRequest) (*policypb.CheckResponse, error) {
	log.Printf("[Policy] Checking user=%s action=%s", req.Username, req.Action)
	if req.Username == "blocked" {
		return &policypb.CheckResponse{
			Allowed: false,
			Reason:  "user is blocked",
		}, nil
	}
	return &policypb.CheckResponse{
		Allowed: true,
		Reason:  "ok",
	}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	policypb.RegisterPolicyServiceServer(s, &PolicyServer{})

	// ✅ reflection for grpcurl
	reflection.Register(s)

	log.Println("[Policy] gRPC server running on :8082")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
