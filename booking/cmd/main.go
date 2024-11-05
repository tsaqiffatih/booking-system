package main

import (
	"log"
	"net"

	booking_service "github.com/tsaqiffatih/booking-system/booking"
	"github.com/tsaqiffatih/booking-system/booking/bookingPb"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=12345678 dbname=booking-service port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database for booking service")
	}

	db.AutoMigrate(&booking_service.Booking{})

	bookingRepo := booking_service.NewBookingRepository(db)
	bookingService := booking_service.NewBookingService(bookingRepo)
	bookingHandler := booking_service.NewBookingHandler(bookingService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookingPb.RegisterBookingServiceServer(grpcServer, bookingHandler)

	log.Println("gRPC Booking server is running on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
