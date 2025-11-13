package biz

import (
	"eGZ-stu-log/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

type PlaceUseCase struct {
	repo data.PlaceRepo
	log  *log.Helper
}

func NewPlaceUseCase(repo data.PlaceRepo, logger log.Logger) *PlaceUseCase {
	return &PlaceUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
