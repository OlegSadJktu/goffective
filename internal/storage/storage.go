package storage

import (
	"context"

	"github.com/go-pg/pg/v10"
)

type Storage struct {
	db *pg.DB
}

func New(url string) (*Storage, error) {
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
  err = db.Ping(context.Background())
  if err != nil {
    return nil, err
  }
	return &Storage{
		db: db,
	}, nil

}
