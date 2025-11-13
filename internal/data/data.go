package data

import (
	"context"
	"eGZ-stu-log/internal/conf"
	"eGZ-stu-log/internal/data/ent"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewPlaceRepo)

// Data .
type Data struct {
	DB *ent.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logs := log.NewHelper(logger)
	dsn := c.Database.Source
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		logs.Errorf("failed opening connection to mysql: %v", err)
		return nil, nil, err
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		logs.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}
	client = client.Debug()
	data := &Data{
		DB: client,
	}
	cleanup := func() {
		logs.Info("closing the data resources")
		if err := client.Close(); err != nil {
			log.Errorf("failed to close database: %v", err)
		}
	}

	return data, cleanup, nil
}
