package service

import (
	"eGZ-stu-log/internal/biz"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type UploadService struct {
	importBiz *biz.ImportUseCase
}

func NewUploadService(importBiz *biz.ImportUseCase) *UploadService {
	return &UploadService{
		importBiz: importBiz,
	}
}
func (u *UploadService) UploadHandler(ctx http.Context) error {
	var err error
	var savedPath string
	resp := map[string]interface{}{
		"code":    200,
		"message": "上传成功",
	}
	req := ctx.Request()
	mode := req.URL.Query().Get("mode")

	// 保存文件
	file, handler, err := req.FormFile("file")
	if err != nil {
		return ctx.String(400, "NoFileInput")
	}

	defer file.Close()
	switch mode {
	case "importStudent":
		savedPath, _ = saveFile(file, handler, "resource/upload/tmp")
		err = u.importBiz.ImportStudent(ctx, savedPath)
		if err == nil {
			resp["code"] = 200
			resp["message"] = "导入成功！"
		} else {
			resp["code"] = 400
			resp["message"] = err.Error()
		}
	case "image":
		savedPath, _ = saveFile(file, handler, "resource/upload/images")
		resp["path"] = savedPath
	case "importRule":
		savedPath, _ = saveFile(file, handler, "resource/upload/tmp")
		err = u.importBiz.ImportRule(ctx, savedPath)
		if err == nil {
			resp["code"] = 200
			resp["message"] = "导入成功！"
		} else {
			resp["code"] = 400
			resp["message"] = err.Error()
		}
	default:
		resp["code"] = 400
		resp["message"] = "mode参数错误"
		return ctx.JSON(400, resp)
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
