package grpcserver

import (
	"context"
	"log"
	"net"

	threatlogpb "DistributedSecurity/pkg/proto/threatlogpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ThreatLogServer struct {
	threatlogpb.UnimplementedThreatLogServiceServer
	events []string
}

func (s *ThreatLogServer) LogEvent(ctx context.Context, req *threatlogpb.LogRequest) (*threatlogpb.LogResponse, error) {
	log.Printf("[ThreatLog] Event received: user=%s type=%s", req.Username, req.EventType)
	s.events = append(s.events, req.Username+" - "+req.EventType)
	return &threatlogpb.LogResponse{
		Success: true,
		Message: "Event logged successfully",
	}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	threatlogpb.RegisterThreatLogServiceServer(s, &ThreatLogServer{})

	// ✅ reflection for grpcurl
	reflection.Register(s)

	log.Println("[ThreatLog] gRPC server running on :8081")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
