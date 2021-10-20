package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ViniSantos88/test-rabbitmq-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Falha na conexão com o grpc server: %v", err)
	}
	defer connection.Close()

	client := pb.NewBankAccountServiceClient(connection)

	//	AddBankAccount(client)
	//AddBankAccountStream(client)
	//AddAccounts(client)
	AddAccountStreamBoth(client)
}

func AddBankAccount(client pb.BankAccountServiceClient) {
	req := &pb.BankAccount{
		Id:          "12345",
		Name:        "Vinícius",
		BankBalance: 100.60,
	}

	res, err := client.AddBankAccount(context.Background(), req)
	if err != nil {
		log.Fatalf("Não foi possivel realizar a request via grpc: %v", err)
	}

	log.Println(res)

}

func AddBankAccountStream(client pb.BankAccountServiceClient) {

	req := &pb.BankAccount{
		Id:          "12345",
		Name:        "Vinícius",
		BankBalance: 100.60,
	}

	responseStream, err := client.AddBankAccountStream(context.Background(), req)
	if err != nil {
		log.Fatalf("Falha na request grpc: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Falha no recebimento do stream: %v", err)
		}

		fmt.Println("Status:", stream.Status)
	}
}

func AddAccounts(client pb.BankAccountServiceClient) {
	reqs := []*pb.BankAccount{
		&pb.BankAccount{
			Id:          "1",
			Name:        "Vinicius",
			BankBalance: 189.90,
		},
		&pb.BankAccount{
			Id:          "2",
			Name:        "Thomas",
			BankBalance: 200.11,
		},
		&pb.BankAccount{
			Id:          "3",
			Name:        "Camila",
			BankBalance: 202.11,
		},
	}

	stream, err := client.AddAccounts(context.Background())
	if err != nil {
		log.Fatalf("Falha na request grpc: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Falha no recebimento do stream: %v", err)
	}

	fmt.Println(res)
}

func AddAccountStreamBoth(client pb.BankAccountServiceClient) {
	stream, err := client.AddAccountStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Falha no envio ao stream: %v", err)
	}

	reqs := []*pb.BankAccount{
		&pb.BankAccount{
			Id:          "1",
			Name:        "Vinicius",
			BankBalance: 189.90,
		},
		&pb.BankAccount{
			Id:          "2",
			Name:        "Thomas",
			BankBalance: 200.11,
		},
		&pb.BankAccount{
			Id:          "3",
			Name:        "Camila",
			BankBalance: 202.11,
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Enviando bankAccount", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Falha ao receber stream: %v", err)
				break
			}
			fmt.Printf("Recebendo name %v com status: %v\n", res.GetBankAccount().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}
