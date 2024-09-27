package ssongs

import "github.com/OlegSadJktu/goffective/internal/repository/rsongs"

type SongService struct {
	repo *rsongs.SongRepo
}

func New(repo *rsongs.SongRepo) *SongService {
	return &SongService{
		repo: repo,
	}

}
