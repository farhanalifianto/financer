package grpc_server

import (
	"category-service/model"
	pb "category-service/proto/category"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type CategoryServer struct {
	pb.UnimplementedCategoryServiceServer
	DB *gorm.DB
}

func (s *CategoryServer) GetCategoryByID(ctx context.Context, req *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
	var cat model.Category
	if err := s.DB.First(&cat, req.Id).Error; err != nil {
		return nil,  status.Errorf(codes.NotFound, "category not found")
	}

	return &pb.CategoryResponse{
		Id:   uint32(cat.ID),
		Name: cat.Name,
		Type: cat.Type,
	}, nil
}
