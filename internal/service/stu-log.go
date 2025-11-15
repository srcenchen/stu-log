package service

import (
	"context"
	"eGZ-stu-log/internal/data"
	"eGZ-stu-log/internal/data/ent"
	"eGZ-stu-log/internal/data/ent/grade"
	"eGZ-stu-log/internal/data/ent/rule"
	"eGZ-stu-log/internal/data/ent/student"
	"eGZ-stu-log/internal/data/ent/stulog"
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

func toStuLogItem(stuLogInfo *ent.StuLog) *pb.StuLogItem {
	var imageUrls []string
	for _, image := range stuLogInfo.Edges.Images {
		imageUrls = append(imageUrls, image.ImageUrl)
	}
	var students []*pb.StuLogStudentItem
	for _, studentItem := range stuLogInfo.Edges.Students {
		students = append(students, &pb.StuLogStudentItem{
			Id:     studentItem.ID,
			Name:   studentItem.Name,
			Class:  studentItem.Edges.Class.ClassName,
			StuNum: studentItem.StuNum,
		})
	}
	return &pb.StuLogItem{
		Id:        stuLogInfo.ID,
		Content:   stuLogInfo.Content,
		Time:      stuLogInfo.Time.String(),
		Score:     stuLogInfo.Score,
		Rule:      stuLogInfo.Edges.Rule.Content,
		Grade:     stuLogInfo.Edges.Grade.GradeName,
		Revoked:   stuLogInfo.Revoked,
		ImageUrls: imageUrls,
		Students:  students,
	}
}

// GetStuLog 获取学生日志
func (s *StuLogService) GetStuLog(ctx context.Context, req *pb.GetStuLogRequest) (*pb.StuLogItem, error) {
	// 根据ID查询
	stuLogInfo, err := s.data.DB.StuLog.Query().
		WithGrade().
		WithClass().
		WithImages().
		WithRule().
		WithStudents(func(sq *ent.StudentQuery) {
			sq.WithClass()
		}).
		Where(stulog.ID(req.Id)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toStuLogItem(stuLogInfo), nil
}

// GetStuLogList 获取学生日志
func (s *StuLogService) GetStuLogList(ctx context.Context, req *pb.GetStuLogListRequest) (*pb.GetStuLogListReply, error) {
	dbQuery := s.data.DB.StuLog.Query().
		WithGrade().
		WithClass().
		WithImages().
		WithRule().
		WithStudents(func(sq *ent.StudentQuery) {
			sq.WithClass()
		})
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.GradeId != nil {
		dbQuery.Where(stulog.HasGradeWith(grade.ID(*req.GradeId)))
	}
	if req.RuleId != nil {
		dbQuery.Where(stulog.HasRuleWith(rule.ID(*req.RuleId)))
	}
	if req.StudentId != nil {
		dbQuery.Where(stulog.HasStudentsWith(student.ID(*req.StudentId)))
	}
	if req.StartTime != nil {
		dbQuery.Where(stulog.TimeGTE(time.Unix(*req.StartTime, 0)))
	}
	if req.EndTime != nil {
		dbQuery.Where(stulog.TimeLTE(time.Unix(*req.EndTime, 0)))
	}
	total, err := dbQuery.Count(ctx)
	totalPages := (int64(total) + req.PageSize - 1) / req.PageSize
	stuLogList, err := dbQuery.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize)).All(ctx)
	if err != nil {
		return nil, err
	}
	var stuLogListReply []*pb.StuLogItem
	for _, stuLogInfo := range stuLogList {
		stuLogListReply = append(stuLogListReply, toStuLogItem(stuLogInfo))
	}
	return &pb.GetStuLogListReply{
		TotalPages: totalPages,
		List:       stuLogListReply,
	}, nil
}

// ReportStuLog 上报学生日志
func (s *StuLogService) ReportStuLog(ctx context.Context, req *pb.ReportStuLogRequest) (*pb.ReportStuLogReply, error) {
	// 首先创建核心违纪信息
	stuLogInfo := s.data.DB.StuLog.Create().
		SetContent(req.Content).
		SetGradeID(req.GradeId).
		SetTime(time.Now()).
		SetScore(req.Score).
		AddImageIDs(req.ImageIds...).
		SetRuleID(req.RuleId).
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
