
### Structure Project
```
booking_system/
│
├── user/
│   ├── cmd/
│   │   └── main.go   
│   ├── user_handler.go
│   ├── user_model.go
│   ├── user_repository.go
│   ├── user_service.go
│   ├── user_client.go
│   ├── userPb/
│   ├── user.proto
│   └── app.dockerfile
│
├── booking/
│   ├── cmd/
│   │   └── main.go   
│   ├── booking_handler.go
│   ├── booking_model.go
│   ├── booking_repository.go
│   ├── booking_service.go
│   ├── booking_client.go
│   ├── bookingPb/
│   ├── booking.proto
│   └── app.dockerfile
│
├── payment/
│   ├── cmd/
│   │   └── main.go   
│   ├── payment_handler.go
│   ├── payment_model.go
│   ├── payment_scheduler.go
│   ├── payment_repository.go
│   ├── payment_service.go
│   ├── payment_client.go 
│   ├── paymentPb/
│   ├── payment.proto
│   └── app.dockerfile
│
├── notification/
│   ├── cmd/
│   │   └── main.go   
│   ├── notification_handler.go
│   ├── notification_model.go
│   ├── notification_repository.go
│   ├── notification_service.go
│   ├── notification_client.go
│   ├── notificationPb/
│   ├── notification.proto
│   └── app.dockerfile
│
├── hotel/
│   ├── cmd/
│   │   └── main.go   
│   ├── room_handler.go
│   ├── room_model.go
│   ├── room_repository.go
│   ├── room_service.go
│   ├── room_client.go
│   ├── roomPb/
│   ├── room.proto
│   └── app.dockerfile
│
├── hotel/ belum terbuat
│
├── Makefile
├── go.mod
├── docker-compose.yml
└── README.md
```