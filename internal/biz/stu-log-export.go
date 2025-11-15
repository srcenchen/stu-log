package biz

import (
	"context"
	"eGZ-stu-log/internal/data"
	"eGZ-stu-log/internal/data/ent"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xuri/excelize/v2"
)

type ExportStuLogUseCase struct {
	data *data.Data
	log  *log.Helper
}

func NewExportStuLogUseCase(data *data.Data, logger log.Logger) *ExportStuLogUseCase {
	return &ExportStuLogUseCase{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (s *ExportStuLogUseCase) ExportStuLog(ctx context.Context, stuLogs []*ent.StuLog) (string, error) {
	// 创建新的 Excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// 创建工作表
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println("Error creating sheet:", err)
		return "", err
	}

	// 设置表头（第一行）
	headers := []string{"事件", "学生", "条例", "备注", "年级", "班级", "分数", "上报时间", "宿舍楼", "宿舍号", "已撤销"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+i, 1)
		_ = f.SetCellValue("Sheet1", cell, header)
	}

	// 写入每条日志数据
	for i, log := range stuLogs {
		// 解析学生表
		var studentArray []string
		for _, student := range log.Edges.Students {
			studentArray = append(studentArray, student.Name)
		}
		// 解析班级表
		var classArray []string
		for _, class := range log.Edges.Class {
			classArray = append(classArray, class.ClassName)
		}
		row := i + 2                                                                           // 第二行开始写数据
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)                             // 事件编号
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), strings.Join(studentArray, " ")) // 学生
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), log.Edges.Rule.Content)          // 条例
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), log.Content)                     // 备注
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), log.Edges.Grade.GradeName)       // 年级
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), strings.Join(classArray, " "))   // 班级
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), log.Score)                       // 分数
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), log.Time.String())               // 上报时间
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("K%d", row), log.Revoked)                     // 撤销
		if log.Edges.Dorm != nil {
			f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), log.Edges.Dorm.Building) // 宿舍（楼栋）
			f.SetCellValue("Sheet1", fmt.Sprintf("J%d", row), log.Edges.Dorm.DormNum)  // 宿舍（房间号）
		}
	}
	// 设置默认工作表
	f.SetActiveSheet(index)

	// 添加样式
	// 设置全局字体大小为11，居中对齐，允许文本换行
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		fmt.Println("Error creating style:", err)
		return "", err
	}

	// 应用样式到所有单元格（包括表头和数据行）
	for i := 0; i < len(headers); i++ {
		for j := 0; j < len(stuLogs)+1; j++ { // 修正：去掉多余的+1
			cell, _ := excelize.CoordinatesToCellName(i+1, j+1)
			_ = f.SetCellStyle("Sheet1", cell, cell, style)
		}
	}

	// 设置行高为57（包括表头和数据行）
	for i := 1; i < len(stuLogs)+1; i++ { // 修正：调整循环范围
		_ = f.SetRowHeight("Sheet1", i+1, 57)
	}

	// 设置列宽
	for i := 0; i < len(headers); i++ {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		_ = f.SetColWidth("Sheet1", colName, colName, float64(len(headers[i])+5)) // 基础宽度+5
	}

	// 调整特定列的宽度以适应内容
	_ = f.SetColWidth("Sheet1", "A", "A", 8)  // 事件
	_ = f.SetColWidth("Sheet1", "B", "B", 20) // 学生
	_ = f.SetColWidth("Sheet1", "C", "C", 30) // 条例
	_ = f.SetColWidth("Sheet1", "D", "D", 20) // 备注
	_ = f.SetColWidth("Sheet1", "E", "E", 10) // 年级
	_ = f.SetColWidth("Sheet1", "F", "F", 10) // 班级
	_ = f.SetColWidth("Sheet1", "G", "G", 8)  // 分数
	_ = f.SetColWidth("Sheet1", "H", "H", 20) // 上报时间
	_ = f.SetColWidth("Sheet1", "I", "I", 10) // 宿舍楼栋
	_ = f.SetColWidth("Sheet1", "J", "J", 10) // 宿舍房间号

	// 设置图片列的宽度
	for i := 0; i < 10; i++ { // 假设最多10张图片
		colName, _ := excelize.ColumnNumberToName(11 + i)
		_ = f.SetColWidth("Sheet1", colName, colName, 20)
	}

	// 保存文件
	err = f.SaveAs("./resource/upload/tmp/export.xlsx")
	if err != nil {
		fmt.Println("Error saving file:", err)
		return "", err
	}

	fmt.Println("日志已成功导出至 ./resource/upload/temp/export.xlsx")
	return "/resource/upload/temp/export.xlsx", err
}
