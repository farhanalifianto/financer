package grpc_client

import (
	"context"
	"log"
	"time"

	pb "transaction-service/proto/user" // hasil generate proto, copy dari user-service

	"google.golang.org/grpc"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient() *UserClient {
	conn, err := grpc.Dial("user-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to user-service: %v", err)
	}
	c := pb.NewUserServiceClient(conn)
	return &UserClient{client: c}
}
type UserInfo struct {
    Email string
    Name  string
}

func (uc *UserClient) GetUserInfo(userID uint) (*UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	res, err := uc.client.GetUserByID(ctx, &pb.GetUserRequest{Id: uint32(userID)})
	if err != nil {
		return nil, err
	}
	 return &UserInfo{
        Email: res.Email,
        Name:  res.Name,
    }, nil
}

