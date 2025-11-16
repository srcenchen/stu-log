package service

import (
	"context"
	pb "eGZ-stu-log/api/base_info/v1"
	"eGZ-stu-log/internal/biz"
	"eGZ-stu-log/internal/data"
	"eGZ-stu-log/internal/data/ent"
	"eGZ-stu-log/internal/data/ent/class"
	"eGZ-stu-log/internal/data/ent/grade"
	"eGZ-stu-log/internal/data/ent/student"
)

type StudentService struct {
	pb.UnimplementedStudentServer
	data      *data.Data
	exportBiz *biz.ExportStuLogUseCase
}

func NewStudentService(data *data.Data, exportBiz *biz.ExportStuLogUseCase) *StudentService {
	return &StudentService{
		data:      data,
		exportBiz: exportBiz,
	}
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {
	// 设置学生
	_, err := s.data.DB.Student.Create().
		SetName(req.Name).
		SetStuNum(req.StuNum).
		SetGradeID(req.GradeId).
		SetSex(req.Sex).
		SetClassID(req.ClassId).
		Save(ctx)

	if err != nil {
		return nil, err
	}
	return &pb.CreateStudentReply{
		Message: "学生添加成功",
	}, nil
}
func (s *StudentService) QueryStudent(ctx context.Context, req *pb.QueryStudentRequest) (*pb.QueryStudentReply, error) {
	// 根据id查询学生
	studentQuery, err := s.data.DB.Student.Query().WithDorm().WithGrade().WithClass().Where(student.ID(req.Id)).First(ctx)
	if err != nil {
		return nil, err
	}
	var dorm = &ent.Dorm{}
	if studentQuery.Edges.Dorm != nil {
		dorm = studentQuery.Edges.Dorm
	}
	return &pb.QueryStudentReply{
		Student: &pb.StudentItem{
			Id:           studentQuery.ID,
			Name:         studentQuery.Name,
			StuNum:       studentQuery.StuNum,
			Score:        studentQuery.Score,
			GradeId:      studentQuery.Edges.Grade.ID,
			GradeName:    studentQuery.Edges.Grade.GradeName,
			Sex:          studentQuery.Sex,
			ClassId:      studentQuery.Edges.Class.ID,
			ClassName:    studentQuery.Edges.Class.ClassName,
			DormBuilding: dorm.Building,
			DormName:     dorm.DormNum,
			DormId:       dorm.ID,
			DormPos:      studentQuery.DormPos,
		},
	}, nil
}
func (s *StudentService) QueryStudentList(ctx context.Context, req *pb.QueryStudentListRequest) (*pb.QueryStudentListReply, error) {
	dbQuery := s.data.DB.Student.Query().WithDorm().WithGrade().WithClass()
	if req.GradeId != nil && *req.GradeId != -1 {
		dbQuery.Where(student.HasGradeWith(grade.ID(*req.GradeId)))
	}
	if req.ClassId != nil {
		dbQuery.Where(student.HasClassWith(class.ID(*req.ClassId)))
	}
	if req.Search != nil && *req.Search != "" {
		dbQuery.Where(student.Or(student.NameContains(*req.Search), student.StuNumContains(*req.Search)))
	}
	total, err := dbQuery.Count(ctx)
	studentList, err := dbQuery.Order(ent.Asc("score")).Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize)).All(ctx)
	if err != nil {
		return nil, err
	}
	var studentListReply []*pb.StudentItem

	for _, studentQuery := range studentList {
		var dorm = &ent.Dorm{}
		if studentQuery.Edges.Dorm != nil {
			dorm = studentQuery.Edges.Dorm
		}
		studentListReply = append(studentListReply, &pb.StudentItem{
			Id:           studentQuery.ID,
			Name:         studentQuery.Name,
			StuNum:       studentQuery.StuNum,
			Score:        studentQuery.Score,
			GradeId:      studentQuery.Edges.Grade.ID,
			GradeName:    studentQuery.Edges.Grade.GradeName,
			Sex:          studentQuery.Sex,
			ClassId:      studentQuery.Edges.Class.ID,
			ClassName:    studentQuery.Edges.Class.ClassName,
			DormBuilding: dorm.Building,
			DormName:     dorm.DormNum,
			DormId:       dorm.ID,
			DormPos:      studentQuery.DormPos,
		})
	}
	return &pb.QueryStudentListReply{
		Students: studentListReply,
		Total:    int32(total),
	}, nil
}

// ExportStudent 导出
func (s *StudentService) ExportStudent(ctx context.Context, req *pb.ExportStuRequest) (*pb.ExportStuReply, error) {
	dbQuery := s.data.DB.Student.Query().
		WithGrade().
		WithClass().WithDorm()
	if req.GradeId != nil && *req.GradeId != -1 {
		dbQuery.Where(student.HasGradeWith(grade.ID(*req.GradeId)))
	}
	stuLogList, err := dbQuery.Order(ent.Asc("score")).All(ctx)
	if err != nil {
		return nil, err
	}
	outPutPath, err := s.exportBiz.ExportStudent(ctx, stuLogList)
	return &pb.ExportStuReply{ExportPath: outPutPath}, nil
}
