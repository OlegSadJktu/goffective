package csongs

import (
	"github.com/OlegSadJktu/goffective/internal/service/ssongs"
)

type SongController struct {
	service *ssongs.SongService
}

func New(service *ssongs.SongService) *SongController {
	return &SongController{
		service: service,
	}
}
