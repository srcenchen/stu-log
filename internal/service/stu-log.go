package service

import (
	"context"
	"eGZ-stu-log/internal/biz"
	"eGZ-stu-log/internal/data"
	"eGZ-stu-log/internal/data/ent"
	"eGZ-stu-log/internal/data/ent/dorm"
	"eGZ-stu-log/internal/data/ent/grade"
	"eGZ-stu-log/internal/data/ent/rule"
	"eGZ-stu-log/internal/data/ent/student"
	"eGZ-stu-log/internal/data/ent/stulog"
	"time"

	pb "eGZ-stu-log/api/base_info/v1"
)

type StuLogService struct {
	pb.UnimplementedStuLogServer
	data      *data.Data
	exportBiz *biz.ExportStuLogUseCase
}

func NewStuLogService(data *data.Data, exportBiz *biz.ExportStuLogUseCase) *StuLogService {
	return &StuLogService{
		data:      data,
		exportBiz: exportBiz,
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
	dormInfo := &ent.Dorm{}
	if stuLogInfo.Edges.Dorm != nil {
		dormInfo = stuLogInfo.Edges.Dorm
	}
	return &pb.StuLogItem{
		Id:           stuLogInfo.ID,
		Content:      stuLogInfo.Content,
		Time:         stuLogInfo.Time.Format("2006-01-02 15:04:05"),
		Score:        stuLogInfo.Score,
		Rule:         stuLogInfo.Edges.Rule.Content + "-" + stuLogInfo.Edges.Rule.Group,
		Grade:        stuLogInfo.Edges.Grade.GradeName,
		Revoked:      stuLogInfo.Revoked,
		ImageUrls:    imageUrls,
		DormBuilding: dormInfo.Building,
		DormNum:      dormInfo.DormNum,
		Students:     students,
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
		WithDorm().
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
		WithDorm().
		WithStudents(func(sq *ent.StudentQuery) {
			sq.WithClass()
		})
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 99999
	}
	if req.GradeId != nil && *req.GradeId != -1 {
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
	if req.OnlyDorm != nil && *req.OnlyDorm {
		dbQuery.Where(stulog.HasDorm())
	}
	total, err := dbQuery.Count(ctx)
	stuLogList, err := dbQuery.Order(ent.Desc("time")).Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize)).All(ctx)
	if err != nil {
		return nil, err
	}
	var stuLogListReply []*pb.StuLogItem
	for _, stuLogInfo := range stuLogList {
		stuLogListReply = append(stuLogListReply, toStuLogItem(stuLogInfo))
	}
	return &pb.GetStuLogListReply{
		Total: int64(total),
		List:  stuLogListReply,
	}, nil
}

// ExportStuLog 导出学生日志
func (s *StuLogService) ExportStuLog(ctx context.Context, req *pb.ExportStuLogRequest) (*pb.ExportStuLogReply, error) {
	dbQuery := s.data.DB.StuLog.Query().
		WithGrade().
		WithClass().
		WithImages().
		WithRule().
		WithDorm().
		WithStudents(func(sq *ent.StudentQuery) {
			sq.WithClass()
		})
	if req.GradeId != nil && *req.GradeId != -1 {
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
	stuLogList, err := dbQuery.Order(ent.Desc("time")).All(ctx)
	if err != nil {
		return nil, err
	}
	outPutPath, err := s.exportBiz.ExportStuLog(ctx, stuLogList)
	return &pb.ExportStuLogReply{ExportPath: outPutPath}, nil
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
	if req.DormId != nil {
		dormInfo, _ := s.data.DB.Dorm.Query().Where(dorm.ID(*req.DormId)).First(ctx)
		stuLogInfo.SetDorm(dormInfo)
	}
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

func (s *StuLogService) RevokeStuLog(ctx context.Context, req *pb.RevokeStuLogRequest) (*pb.RevokeStuLogReply, error) {
	logInfoOrigin, err := s.data.DB.StuLog.Query().Where(stulog.ID(req.Id)).First(ctx)
	if err != nil {
		return nil, err
	}
	if logInfoOrigin.Revoked == req.Revoke {
		return &pb.RevokeStuLogReply{
			Message: "操作失败",
		}, nil
	}
	logInfo, err := s.data.DB.StuLog.UpdateOneID(req.Id).SetRevoked(req.Revoke).Save(ctx)
	if err != nil {
		return nil, err
	}
	// 通过id查询Student
	stuLogsInfo, _ := s.data.DB.StuLog.Query().WithStudents(func(sq *ent.StudentQuery) {
		sq.WithClass()
	}).Where(stulog.ID(req.Id)).First(ctx)
	if req.Revoke {
		for _, studentInfo := range stuLogsInfo.Edges.Students {
			_, _ = s.data.DB.Student.UpdateOneID(studentInfo.ID).SetScore(studentInfo.Score - logInfo.Score).Save(ctx)
		}
	} else {
		for _, studentInfo := range stuLogsInfo.Edges.Students {
			_, _ = s.data.DB.Student.UpdateOneID(studentInfo.ID).SetScore(studentInfo.Score + logInfo.Score).Save(ctx)
		}
	}
	return &pb.RevokeStuLogReply{
		Message: "操作成功",
	}, nil
}

func (s *StuLogService) GetRuleList(ctx context.Context, req *pb.GetRuleRequest) (*pb.GetRuleReply, error) {
	ruleList, err := s.data.DB.Rule.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	var ruleInfo []*pb.RuleItem
	for _, rule := range ruleList {
		ruleInfo = append(ruleInfo, &pb.RuleItem{
			Id:    rule.ID,
			Group: rule.Group,
			Rule:  rule.Content,
			Score: rule.Score,
		})
	}
	return &pb.GetRuleReply{List: ruleInfo}, nil

}
