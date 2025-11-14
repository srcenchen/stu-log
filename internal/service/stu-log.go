package service

import (
	"context"
	"eGZ-stu-log/internal/data"
	"eGZ-stu-log/internal/data/ent/student"
	"time"

	pb "eGZ-stu-log/api/base_info/v1"
)

type StuLogService struct {
	pb.UnimplementedStuLogServer
	data *data.Data
}

func NewStuLogService(data *data.Data) *StuLogService {
	return &StuLogService{
		data: data,
	}
}

func (s *StuLogService) ReportStuLog(ctx context.Context, req *pb.ReportStuLogRequest) (*pb.ReportStuLogReply, error) {
	// 首先创建核心违纪信息
	stuLogInfo := s.data.DB.StuLog.Create().
		SetContent(req.Content).
		AddGradeIDs(req.GradeId).
		SetTime(time.Now()).
		SetScore(req.Score).
		AddImageIDs(req.ImageIds...).
		AddRuleIDs(req.RuleId).
		AddStudentIDs(req.StudentIds...)
	// 根据ids查询classId 塞进去
	for _, stuId := range req.StudentIds {
		studentInfo, _ := s.data.DB.Student.Query().WithClass().Where(student.ID(stuId)).First(ctx)
		stuLogInfo.AddClass(studentInfo.Edges.Class)
		_, _ = s.data.DB.Student.UpdateOneID(stuId).SetScore(studentInfo.Score + req.Score).Save(ctx)
	}
	_, err := stuLogInfo.Save(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ReportStuLogReply{
		Message: "Success",
	}, nil
}
