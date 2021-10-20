package main

import (
	"log"
	"net"

	"github.com/ViniSantos88/test-rabbitmq-grpc/pb"
	"github.com/ViniSantos88/test-rabbitmq-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Falha na conexão: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBankAccountServiceServer(grpcServer, services.NewBankAccountService())
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Falha no serve: %v", err)
	}

}
