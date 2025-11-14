package service

import (
	"context"
	pb "eGZ-stu-log/api/base_info/v1"
	"eGZ-stu-log/internal/data"
	"eGZ-stu-log/internal/data/ent"
	"eGZ-stu-log/internal/data/ent/class"
	"eGZ-stu-log/internal/data/ent/grade"
	"eGZ-stu-log/internal/data/ent/student"
)

type StudentService struct {
	pb.UnimplementedStudentServer
	data *data.Data
}

func NewStudentService(data *data.Data) *StudentService {
	return &StudentService{
		data: data,
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
	// 查询总页数
	total, err := s.data.DB.Student.Query().Count(ctx)
	totalPages := (int32(total) + req.PageSize - 1) / req.PageSize
	dbQuery := s.data.DB.Student.Query().WithDorm().WithGrade().WithClass().Offset(int((req.Page - 1) * req.PageSize))
	if req.GradeId != nil {
		dbQuery.Where(student.HasGradeWith(grade.ID(*req.GradeId)))
	}
	if req.ClassId != nil {
		dbQuery.Where(student.HasClassWith(class.ID(*req.ClassId)))
	}
	studentList, err := dbQuery.Limit(int(req.PageSize)).All(ctx)
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
		Students:   studentListReply,
		TotalPages: totalPages,
	}, nil
}
