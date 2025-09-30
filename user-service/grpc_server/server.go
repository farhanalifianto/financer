package grpc_server

import (
	"context"
	"user-service/model"
	pb "user-service/proto/user"

	"log"

	"gorm.io/gorm"
)

type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	DB *gorm.DB
}

func (s *UserGRPCServer) GetUserByID(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    var user model.User
    if err := s.DB.First(&user, req.Id).Error; err != nil {
        log.Printf("User not found: %v", err)  // tambahin log
        return nil, err
    }

    return &pb.GetUserResponse{
        Id:    uint32(user.ID),
        Email: user.Email,
        Name:  user.Name,
        Role:  user.Role,
    }, nil
}