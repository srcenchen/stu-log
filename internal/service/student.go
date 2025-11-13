package service

import (
	"context"
	pb "eGZ-stu-log/api/base_info/v1"
	"eGZ-stu-log/internal/biz"
	"eGZ-stu-log/internal/data"
)

type StudentService struct {
	pb.UnimplementedStudentServer
	placeUseCase *biz.PlaceUseCase
	data         *data.Data
}

func NewStudentService(placeUseCase *biz.PlaceUseCase, data *data.Data) *StudentService {
	return &StudentService{placeUseCase: placeUseCase, data: data}
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {
	//s.data.DB.Grade.Create().SetGradeName("123").Save(ctx)
	s.data.DB.Class.Create().SetName("123").SetGradeID(1).Save(ctx)
	// 设置学生
	_, err := s.data.DB.Student.Create().SetName(req.Name).SetGradeID(int(req.GradeId)).SetSex(req.Sex).SetClassID(int(req.ClassId)).SetScore(100).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CreateStudentReply{
		Message: "success",
	}, nil
}
