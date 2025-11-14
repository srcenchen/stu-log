package service

import (
	"context"

	pb "eGZ-stu-log/api/base_info/v1"
)

type ClassService struct {
	pb.UnimplementedClassServer
}

func NewClassService() *ClassService {
	return &ClassService{}
}

func (s *ClassService) CreateClass(ctx context.Context, req *pb.CreateClassRequest) (*pb.CreateClassReply, error) {
    return &pb.CreateClassReply{}, nil
}
