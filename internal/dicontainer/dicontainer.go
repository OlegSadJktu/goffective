package dicontainer

import (
	"github.com/OlegSadJktu/goffective/internal/controllers/csongs"
	"github.com/OlegSadJktu/goffective/internal/repository/rsongs"
	"github.com/OlegSadJktu/goffective/internal/service/ssongs"
	"github.com/go-pg/pg/v10"
)

type DIContainer struct {
	db    *pg.DB
	repos struct {
		*rsongs.SongRepo
	}
	controllers struct {
		*csongs.SongController
	}
	services struct {
		*ssongs.SongService
	}
}

func New(db *pg.DB) *DIContainer {
	return &DIContainer{
		db: db,
	}
}

func (c *DIContainer) SongsRepo() *rsongs.SongRepo {
	if c.repos.SongRepo == nil {
		c.repos.SongRepo = rsongs.New(c.db)
	}
	return c.repos.SongRepo
}

func (c *DIContainer) SongsController() *csongs.SongController {
	if c.controllers.SongController == nil {
		c.controllers.SongController = csongs.New(c.SongsService())
	}
	return c.controllers.SongController
}

func (c *DIContainer) SongsService() *ssongs.SongService {
	if c.services.SongService == nil {
		c.services.SongService = ssongs.New(c.SongsRepo())
	}
	return c.services.SongService
}
