package grpc_client

import (
	"context"
	"log"
	"time"

	pb "transaction-service/proto/category"

	"google.golang.org/grpc"
)

type CategoryClient struct {
	client pb.CategoryServiceClient
}

func NewCategoryClient() *CategoryClient {
	conn, err := grpc.Dial("category-service:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to category-service: %v", err)
	}
	c := pb.NewCategoryServiceClient(conn)
	return &CategoryClient{client: c}
}

type CategoryInfo struct {
	Name string
	Type string
	Budget float64
	OwnerID uint
}

func (cc *CategoryClient) GetCategoryInfo(catID uint) (*CategoryInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := cc.client.GetCategoryByID(ctx, &pb.GetCategoryRequest{Id: uint32(catID)})
	if err != nil {
		return nil, err
	}

	return &CategoryInfo{
		Name: res.Name,
		Type: res.Type,
		Budget: float64(res.Budget),
		OwnerID: uint(res.Ownerid),
	}, nil
}
