package ssongs

import (
	"errors"
	"strings"
	"time"

	msongs "github.com/OlegSadJktu/goffective/internal/model"
	"github.com/OlegSadJktu/goffective/internal/repository/rsongs"
)

const coupletTake = 3

var (
	ErrInvalidCoupletOffset = errors.New("invalid couplet offset")
)

type SongService struct {
	repo *rsongs.SongRepo
}

func New(repo *rsongs.SongRepo) *SongService {
	return &SongService{
		repo: repo,
	}
}

func (s *SongService) Create(model *msongs.Song) (*msongs.Song, error) {
	_, err := s.repo.Create(model)
	if err != nil {
		return nil, err
	}
	return model, nil

}

func (s *SongService) GetAll(min, max time.Time, name, group string, offset int) ([]*msongs.Song, error) {
	songs, err := s.repo.GetAll(min, max, name, group, offset)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *SongService) GetById(id string, offset int) (*msongs.Song, error) {
	song, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	var to int
	couplets := strings.Split(song.Lyrics, "\n\n")
	if offset+coupletTake > len(couplets) {
		to = min(len(couplets), offset+coupletTake)
	}
	if offset > to {
		return nil, ErrInvalidCoupletOffset
	}
	couplets = couplets[offset:to]
	song.Lyrics = strings.Join(couplets, "\n\n")

	return song, nil
}

func (s *SongService) Delete(id string) (*msongs.Song, error) {
	var song msongs.Song
	song.ID = id
	err := s.repo.Delete(&song)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (s *SongService) Update(id string, model *msongs.Song) (*msongs.Song, error) {
	song, err := s.repo.Update(model)
	if err != nil {
		return nil, err
	}
	return song, nil

}
