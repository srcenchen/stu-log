package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type PlaceRepo interface {
	Save(context.Context) (string, error)
}

type placeRepo struct {
	data *Data
	log  *log.Helper
}

func (p placeRepo) Save(ctx context.Context) (string, error) {
	return "nil", nil
}

func NewPlaceRepo(data *Data, logger log.Logger) PlaceRepo {
	return &placeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
