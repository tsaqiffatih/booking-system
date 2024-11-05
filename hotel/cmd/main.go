package main

import (
	"log"
	"net"

	hotel_service "github.com/tsaqiffatih/booking-system/hotel"
	"github.com/tsaqiffatih/booking-system/hotel/hotelPb"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=12345678 dbname=hotel-service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&hotel_service.Hotel{}, &hotel_service.Room{})

	repo := hotel_service.NewHotelRepository(db)
	service := hotel_service.NewHotelService(repo)
	handler := hotel_service.NewHotelHandler(service)

	grpcServer := grpc.NewServer()
	hotelPb.RegisterHotelServiceServer(grpcServer, handler)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Hotel service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
