package rsongs

import (
	"github.com/OlegSadJktu/goffective/internal/model/msongs"
	"github.com/go-pg/pg/v10"
)

type SongRepo struct {
	db *pg.DB
}

func New(db *pg.DB) *SongRepo {
	return &SongRepo{
		db: db,
	}
}

