package service

import (
	"context"
	pb "eGZ-stu-log/api/base_info/v1"
	"eGZ-stu-log/internal/data"
)

type GradeService struct {
	pb.UnimplementedGradeServer
	data *data.Data
}

func NewGradeService(data *data.Data) *GradeService {
	return &GradeService{
		data: data,
	}
}

func (s *GradeService) CreateGrade(ctx context.Context, req *pb.CreateGradeRequest) (*pb.CreateGradeReply, error) {
	_, err := s.data.DB.Grade.Create().SetGradeName(req.GradeName).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CreateGradeReply{
		Message: "年级创建成功",
	}, nil
}
