package server

import (
	v1 "eGZ-stu-log/api/base_info/v1"
	"eGZ-stu-log/internal/conf"
	"eGZ-stu-log/internal/service"
	http2 "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	student *service.StudentService,
	class *service.ClassService,
	grade *service.GradeService,
	upload *service.UploadService,
	stuLog *service.StuLogService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterStudentHTTPServer(srv, student)
	v1.RegisterClassHTTPServer(srv, class)
	v1.RegisterGradeHTTPServer(srv, grade)
	v1.RegisterStuLogHTTPServer(srv, stuLog)
	// 静态资源绑定
	router := mux.NewRouter()
	router.PathPrefix("/resource").Handler(http2.StripPrefix("/resource", http2.FileServer(http2.Dir("./resource"))))
	srv.HandlePrefix("/", router)
	// route
	route := srv.Route("/v1")
	route.POST("/upload", upload.UploadHandler)
	return srv
}
