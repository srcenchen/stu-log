package service

import (
	"eGZ-stu-log/internal/biz"
	"eGZ-stu-log/internal/data"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type UploadService struct {
	importBiz *biz.ImportUseCase
	data      *data.Data
}

func NewUploadService(importBiz *biz.ImportUseCase, data *data.Data) *UploadService {
	return &UploadService{
		importBiz: importBiz,
		data:      data,
	}
}
func (u *UploadService) UploadHandler(ctx http.Context) error {
	req := ctx.Request()
	mode := req.URL.Query().Get("mode")

	// 解析 multipart 表单
	err := req.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		return ctx.String(400, "ParseMultipartFormError")
	}

	files := req.MultipartForm.File["file"]
	if len(files) == 0 {
		return ctx.String(400, "NoFileInput")
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "上传成功",
	}

	var savedPaths []string
	var imageIDs []int64

	for _, handler := range files {
		file, err := handler.Open()
		if err != nil {
			return ctx.String(500, "FileOpenError")
		}
		defer file.Close()

		var savedPath string

		switch mode {
		case "importStudent":
			// 多文件导入学生，依次处理
			savedPath, _ = saveFile(file, handler, "resource/upload/tmp")
			err = u.importBiz.ImportStudent(ctx, savedPath)
			if err != nil {
				resp["code"] = 400
				resp["message"] = err.Error()
			}

		case "image":
			// 多图片保存并写入数据库
			savedPath, _ = saveFile(file, handler, "resource/upload/images")
			imageData, _ := u.data.DB.Image.Create().SetImageUrl(savedPath).Save(ctx)
			imageIDs = append(imageIDs, imageData.ID)

		case "importRule":
			savedPath, _ = saveFile(file, handler, "resource/upload/tmp")
			err = u.importBiz.ImportRule(ctx, savedPath)
			if err != nil {
				resp["code"] = 400
				resp["message"] = err.Error()
			}

		default:
			resp["code"] = 400
			resp["message"] = "mode参数错误"
			return ctx.JSON(400, resp)
		}

		savedPaths = append(savedPaths, savedPath)
	}

	resp["paths"] = savedPaths
	if mode == "image" {
		resp["imageIds"] = imageIDs // 多个图片
	}

	return ctx.JSON(200, resp)
}

// saveFile 保存文件
func saveFile(file multipart.File, header *multipart.FileHeader, dir string) (savedPath string, err error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	savedPath = filepath.Join(dir, fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename))
	dstFile, err := os.Create(savedPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, file)
	if err != nil {
		err = os.Remove(savedPath)
		return "", err
	}
	return
}
