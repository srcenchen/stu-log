package service

import (
	"context"
	"eGZ-stu-log/internal/data"

	pb "eGZ-stu-log/api/base_info/v1"
)

type ClassService struct {
	pb.UnimplementedClassServer
	data *data.Data
}

func NewClassService(data *data.Data) *ClassService {
	return &ClassService{
		data: data,
	}
}

func (s *ClassService) CreateClass(ctx context.Context, req *pb.CreateClassRequest) (*pb.CreateClassReply, error) {
	_, err := s.data.DB.Class.Create().SetClassName(req.ClassName).SetGradeID(req.GradeId).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CreateClassReply{
		Message: req.ClassName + "创建成功",
	}, nil
}
