package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	authpb "DistributedSecurity/pkg/proto/authpb"
	policypb "DistributedSecurity/pkg/proto/policypb"
	threatlogpb "DistributedSecurity/pkg/proto/threatlogpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	policyAddr    string
	threatlogAddr string
}

func NewAuthServer(policyAddr, threatlogAddr string) *AuthServer {
	return &AuthServer{
		policyAddr:    policyAddr,
		threatlogAddr: threatlogAddr,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	log.Printf("[Auth] Login attempt user=%s", req.Username)

	conn, err := grpc.Dial(s.policyAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to policy service: %v", err)
	}
	defer conn.Close()

	policyClient := policypb.NewPolicyServiceClient(conn)
	policyResp, err := policyClient.CheckPolicy(ctx, &policypb.CheckRequest{
		Username: req.Username,
		Action:   "login",
	})
	if err != nil {
		return nil, fmt.Errorf("policy check failed: %v", err)
	}

	if !policyResp.Allowed {
		s.logToThreatLog(ctx, req.Username, "unauthorized_login")
		return &authpb.LoginResponse{
			Valid:   false,
			Message: "Access denied: " + policyResp.Reason,
		}, nil
	}

	return &authpb.LoginResponse{
		Token:   "fake-jwt-token",
		Valid:   true,
		Message: "Login successful",
	}, nil
}

func (s *AuthServer) logToThreatLog(ctx context.Context, username, event string) {
	conn, err := grpc.Dial(s.threatlogAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("[Auth] Could not connect to ThreatLog: %v", err)
		return
	}
	defer conn.Close()

	threatClient := threatlogpb.NewThreatLogServiceClient(conn)
	_, err = threatClient.LogEvent(ctx, &threatlogpb.LogRequest{
		Username:  username,
		EventType: event,
		Timestamp: "now",
	})
	if err != nil {
		log.Printf("[Auth] Failed to log event in ThreatLog: %v", err)
	}
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, NewAuthServer("policy-svc:8082", "threatlog-svc:8081"))

	// ✅ reflection for grpcurl
	reflection.Register(s)

	log.Println("[Auth] gRPC server running on :8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
