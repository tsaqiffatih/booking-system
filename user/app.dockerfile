# Gunakan base image yang sesuai dengan versi Go
FROM golang:1.20-alpine AS builder

# Set environment variables
ENV GO111MODULE=on
WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk caching dependensi
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh kode sumber ke dalam kontainer
COPY . .

# Build aplikasi
RUN go build -o user_service ./user/cmd/main.go

# Image runtime stage yang minimalis
FROM alpine:3.18

# Set environment untuk aplikasi
WORKDIR /app
COPY --from=builder /app/user/user_service /app/user_service

# Expose port sesuai dengan konfigurasi user_service
EXPOSE 50051

# Menjalankan aplikasi
CMD ["./user_service"]
