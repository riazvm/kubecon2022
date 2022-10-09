package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"

	pb "example.com/go-edgemgmt-grpc/edgemgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50051"
)

type EdgeManagementServer struct {
	pb.UnimplementedEdgeManagementServer
}

type EdgePaymentServer struct {
	pb.UnimplementedEdgePaymentServer
}

func (s *EdgeManagementServer) ActivateNewDevice(ctx context.Context, in *pb.NewDevice) (*pb.ActivationStatus, error) {
	log.Printf("Received: %v", in.GetDeviceID())
	return &pb.ActivationStatus{
		MerchantID:   in.GetMerchantID(),
		DeviceID:     in.GetDeviceID(),
		DeviceStatus: "Activated",
	}, nil
}

func (s *EdgePaymentServer) ProcessPayment(ctx context.Context, in *pb.Payment) (*pb.PaymentStatus, error) {
	log.Printf("Received: %v", in.GetPaymentRefNo())
	return &pb.PaymentStatus{
		PaymentRefNo:  in.GetPaymentRefNo(),
		PaymentStatus: "Payment Processed",
	}, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}
	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)
	pb.RegisterEdgeManagementServer(s, &EdgeManagementServer{})
	pb.RegisterEdgePaymentServer(s, &EdgePaymentServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
