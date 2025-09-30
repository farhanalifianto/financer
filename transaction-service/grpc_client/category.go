package grpc_client

import (
	"context"
	"log"
	"time"

	pb "transaction-service/proto/category"

	"google.golang.org/grpc"
)

// CategoryClient untuk menampung koneksi dan client gRPC
type CategoryClient struct {
	client pb.CategoryServiceClient
}

// NewCategoryClient inisialisasi koneksi ke category-service
func NewCategoryClient() *CategoryClient {
	conn, err := grpc.Dial("category-service:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to category-service: %v", err)
	}
	c := pb.NewCategoryServiceClient(conn)
	return &CategoryClient{client: c}
}

// CategoryInfo hasil dari CategoryService
type CategoryInfo struct {
	Name string
	Type string
}

// GetCategoryInfo ambil info kategori berdasarkan ID
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
	}, nil
}
