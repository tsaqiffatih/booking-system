package main

import (
	"log"
	"net"
	"time"

	"github.com/tinrab/retry"
	user_service "github.com/tsaqiffatih/booking-system/user"
	"github.com/tsaqiffatih/booking-system/user/userPb"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=12345678 dbname=user-service port=5432 sslmode=disable"

	var db *gorm.DB
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
			return err
		}

		if err := db.AutoMigrate(&user_service.User{}); err != nil {
			log.Printf("failed to migrate database: %v", err)
			return err
		}

		return
	})

	if db == nil {
		log.Fatalf("database connection is nil")
	}

	userRepo := user_service.NewUserRepository(db)
	userServ := user_service.NewUserService(userRepo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	userPb.RegisterUserServiceServer(grpcServer, user_service.NewUserHandler(userServ))

	log.Println("starting user service gRPC server on port 50051 ...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
