package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ViniSantos88/test-rabbitmq-grpc/pb"
)

type BankAccountService struct {
	pb.UnimplementedBankAccountServiceServer
}

func NewBankAccountService() *BankAccountService {
	return &BankAccountService{}
}

func (*BankAccountService) AddBankAccount(ctx context.Context, req *pb.BankAccount) (*pb.BankAccount, error) {
	return &pb.BankAccount{
		Id:          "1234",
		Name:        req.GetName(),
		BankBalance: req.GetBankBalance(),
	}, nil
}

func (*BankAccountService) AddBankAccountStream(req *pb.BankAccount, stream pb.BankAccountService_AddBankAccountStreamServer) error {
	stream.Send(&pb.AccountResulStream{
		Status:      "Init",
		BankAccount: &pb.BankAccount{},
	})

	time.Sleep(time.Second * 4)

	stream.Send(&pb.AccountResulStream{
		Status:      "Incluíndo",
		BankAccount: &pb.BankAccount{},
	})

	time.Sleep(time.Second * 4)

	stream.Send(&pb.AccountResulStream{
		Status: "Account incluída",
		BankAccount: &pb.BankAccount{
			Id:          "123456",
			Name:        req.GetName(),
			BankBalance: req.GetBankBalance(),
		},
	})

	time.Sleep(time.Second * 4)

	stream.Send(&pb.AccountResulStream{
		Status: "Completado",
		BankAccount: &pb.BankAccount{
			Id:          "123456",
			Name:        req.GetName(),
			BankBalance: req.GetBankBalance(),
		},
	})

	time.Sleep(time.Second * 2)

	return nil
}

func (*BankAccountService) AddAccounts(stream pb.BankAccountService_AddAccountsServer) error {
	accounts := []*pb.BankAccount{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Accounts{
				BankAccount: accounts,
			})
		}

		if err != nil {
			log.Fatalf("Erro ao receber o stream: %v", err)
		}

		accounts = append(accounts, &pb.BankAccount{
			Id:          req.GetId(),
			Name:        req.GetName(),
			BankBalance: req.GetBankBalance(),
		})

		fmt.Println("Adicionado", req.GetName())
	}

	return nil
}

func (*BankAccountService) AddAccountStreamBoth(stream pb.BankAccountService_AddAccountStreamBothServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Erro ao receber o stream do client: %v", err)
		}

		err = stream.Send(&pb.AccountResulStream{
			Status:      "Added",
			BankAccount: req,
		})

		if err != nil {
			log.Fatalf("Erro ao enviar o stream para o client: %v", err)
		}

	}

}
