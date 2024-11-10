package main

import (
	"log"
	"net"

	payment_service "github.com/tsaqiffatih/booking-system/payment"
	"github.com/tsaqiffatih/booking-system/payment/paymentPb"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "your-database-dsn"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := payment_service.NewPaymenRepository(db)
	serverKey := "your-midtrans-server-key"
	isProduction := false
	service := payment_service.NewPaymentService(repo, serverKey, isProduction)
	handler := payment_service.NewPaymentHandler(service)

	grpcServer := grpc.NewServer()
	paymentPb.RegisterPaymentServiceServer(grpcServer, handler)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting gRPC server on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
