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

// GetGrades 获取全部年级信息
func (s *GradeService) GetGrades(ctx context.Context, req *pb.GetGradesRequest) (*pb.GetGradesReply, error) {
	gradeInfos, err := s.data.DB.Grade.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	var gradeList []*pb.GradeItem
	for _, gradeInfo := range gradeInfos {
		gradeList = append(gradeList, &pb.GradeItem{
			GradeName: gradeInfo.GradeName,
			Id:        gradeInfo.ID,
		})
	}
	return &pb.GetGradesReply{
		List: gradeList,
	}, nil
}
