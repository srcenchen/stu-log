package biz

import (
	"eGZ-stu-log/internal/data"
	"eGZ-stu-log/internal/data/ent"
	"eGZ-stu-log/internal/data/ent/class"
	"eGZ-stu-log/internal/data/ent/dorm"
	"eGZ-stu-log/internal/data/ent/grade"
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/xuri/excelize/v2"
)

type ImportUseCase struct {
	data *data.Data
	log  *log.Helper
}

func NewImportUseCase(data *data.Data, logger log.Logger) *ImportUseCase {
	return &ImportUseCase{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uc *ImportUseCase) ImportStudent(ctx http.Context, excelPath string) (err error) {
	f, err := excelize.OpenFile(excelPath)
	defer f.Close()
	rows, err := f.GetRows("Sheet1")
	studentsInfo := rows[1:]
	errCnt := 0
	// 姓名 学号	年级	班级	性别	住校	楼名	宿舍	床铺号
	for _, row := range studentsInfo {
		// 根据年级名获取年级id
		gradeInfo, err := uc.data.DB.Grade.Query().Where(grade.GradeName(row[2])).First(ctx)
		if ent.IsNotFound(err) {
			gradeInfo, err = uc.data.DB.Grade.Create().SetGradeName(row[2]).Save(ctx)
		}
		// 根据年级id和班级名获取班级id
		classInfo, err := uc.data.DB.Class.Query().
			Where(class.ClassName(row[3])).
			Where(class.HasGradeWith(grade.ID(gradeInfo.ID))).First(ctx)
		if ent.IsNotFound(err) {
			classInfo, err = uc.data.DB.Class.Create().
				SetClassName(row[3]).
				SetGrade(gradeInfo).
				Save(ctx)
		} else if err != nil {
			return err
		}
		// 添加学生
		stu := uc.data.DB.Student.Create().
			SetName(row[0]).
			SetStuNum(row[1]).
			SetGrade(gradeInfo).
			SetClass(classInfo).
			SetSex(row[4])
		//添加宿舍情况
		if row[5] == "是" {
			dormInfo, err := uc.data.DB.Dorm.Query().Where(dorm.DormNum(row[7])).Where(dorm.Building(row[6])).First(ctx)
			if ent.IsNotFound(err) {
				dormInfo, err = uc.data.DB.Dorm.Create().
					SetDormNum(row[7]).
					SetBuilding(row[6]).
					SetSex(row[4]).
					SetGrade(gradeInfo).
					Save(ctx)
			}
			stu.SetDorm(dormInfo).SetDormPos(row[8])
		}
		_, err = stu.Save(ctx)
		if err != nil {
			errCnt++
		}
	}
	_ = os.Remove(excelPath) // 导入结束 删除
	if errCnt != 0 {
		return errors.New("部分学生没有导入成功，可能因为导入过了？失败数：" + strconv.Itoa(errCnt))
	}
	return
}

func (uc *ImportUseCase) ImportRule(ctx http.Context, excelPath string) (err error) {
	f, err := excelize.OpenFile(excelPath)
	defer f.Close()
	rows, err := f.GetRows("Sheet1")
	studentsInfo := rows[1:]
	// 分组 规则 分数
	errCnt := 0
	for _, row := range studentsInfo {
		score, _ := strconv.Atoi(row[2])
		_, ers := uc.data.DB.Rule.Create().
			SetGroup(row[0]).
			SetContent(row[1]).
			SetScore(int32(score)).Save(ctx)
		if ers != nil {
			errCnt++
		}
	}
	if errCnt != 0 {
		err = errors.New("部分条例没有导入成功，可能因为重复，失败数:" + strconv.Itoa(errCnt))
	}
	_ = os.Remove(excelPath) // 导入结束 删除
	return
}
