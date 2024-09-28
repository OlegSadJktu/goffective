package rsongs

import (
	"fmt"
	"time"

	msongs "github.com/OlegSadJktu/goffective/internal/model"
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

func (r *SongRepo) Create(model *msongs.Song) (string, error) {
	_, err := r.db.Model(model).Returning("*").Insert()

	if err != nil {
		return "", err
	}
	return model.ID, nil
}

func (r *SongRepo) GetAll(min, max time.Time, name, group string, offset int) ([]*msongs.Song, error) {
	var songs []*msongs.Song
	query := r.db.Model(&songs)
	if !min.IsZero() {
		minformat := min.Format(time.RFC3339)
		query.Where("? < release_date", minformat)
	}
	if !max.IsZero() {
		maxformat := max.Format(time.RFC3339)
		query.Where("? > release_date", maxformat)
	}
	if len(name) > 0 {
		query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if len(group) > 0 {
		query.Where("group ILIKE ?", fmt.Sprintf("%%%s%%", group))
	}
	query.Offset(offset)
	err := query.Select()
	if err != nil {
		return nil, nil
	}
	return songs, nil
}

func (r *SongRepo) GetById(id string) (*msongs.Song, error) {
	var song msongs.Song
	err := r.db.Model(&song).Returning("*").Where("id = ?", id).First()
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *SongRepo) Update(song *msongs.Song) (*msongs.Song, error) {
	_, err := r.db.Model(song).Where("id = ?", song.ID).Returning("*").Update()
	return song, err
}

func (r *SongRepo) Delete(song *msongs.Song) error {
	_, err := r.db.Model(song).Where("id = ?", song.ID).Returning("*").Delete()
	return err
}
