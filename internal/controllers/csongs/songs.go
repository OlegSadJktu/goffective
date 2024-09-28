package csongs

import (
	"net/http"
	"strconv"

	"github.com/OlegSadJktu/goffective/internal/common/responses"
	"github.com/OlegSadJktu/goffective/internal/common/types"
	msongs "github.com/OlegSadJktu/goffective/internal/model"
	"github.com/OlegSadJktu/goffective/internal/service/ssongs"
	"github.com/gin-gonic/gin"
)

type SongController struct {
	service *ssongs.SongService
}

func New(service *ssongs.SongService) *SongController {
	return &SongController{
		service: service,
	}
}

type (
	CreateSongRequest struct {
		ReleaseDate types.CustomTime `json:"release_date" example:"23.03.2024" swaggertype:"primitive,string"`
		Lyrics      string           `json:"text" example:"Oh baby baby baby ooh\n\nOh baby baby"`
		Link        string           `json:"link" example:"https://nescam.io"`
		Name        string           `json:"name" example:"SomeSongUsername"`
		Group       string           `json:"group" example:"Bring me the horizon"`
	}
	CreateSongResponse struct {
		ID          string           `json:"id"`
		ReleaseDate types.CustomTime `json:"release_date" example:"23.03.2024" swaggertype:"primitive,string"`
		Lyrics      string           `json:"text" example:"Oh baby baby baby ooh\n\nOh baby baby"`
		Link        string           `json:"link" example:"https://nescam.io"`
		Name        string           `json:"name" example:"SomeSongUsername"`
		Group       string           `json:"group" example:"Bring me the horizon"`
	}
	UpdateSongRequest struct {
		ReleaseDate types.CustomTime `json:"release_date" example:"23.03.2024" swaggertype:"primitive,string"`
		Lyrics      string           `json:"text" example:"Oh baby baby baby ooh\n\nOh baby baby"`
		Link        string           `json:"link" example:"https://nescam.io"`
		Name        string           `json:"name" example:"SomeSongUsername"`
		Group       string           `json:"group" example:"Bring me the horizon"`
	}
	GetSongsRequest struct {
		ReleaseDateMin types.CustomTime `json:"release_date_min" form:"release_date_min" example:"23.01.2023" swaggertype:"primitive,string"`
		ReleaseDateMax types.CustomTime `json:"release_date_max" form:"release_date_max" example:"23.01.2024" swaggertype:"primitive,string"`
		Offset         int              `json:"offset" form:"offset" example:"3"`
		Name           string           `json:"name" form:"name" example:"on"`
		Group          string           `json:"group" form:"group" example:"the"`
	}
	GetSongsResponse []struct {
		ID          string           `json:"id"`
		ReleaseDate types.CustomTime `json:"release_date" example:"23.03.2024" swaggertype:"primitive,string"`
		Lyrics      string           `json:"text" example:"Oh baby baby baby ooh\n\nOh baby baby"`
		Link        string           `json:"link" example:"https://nescam.io"`
		Name        string           `json:"name" example:"SomeSongUsername"`
		Group       string           `json:"group" example:"Bring me the horizon"`
	}
	GetSongByIdRequest struct {
		ID            string `json:"id" example:"2"`
		CoupletOffset int    `json:"couplet_offset" example:"3"`
	}
	GetSongByIdResponse struct {
		ID          string           `json:"id"`
		ReleaseDate types.CustomTime `json:"release_date" example:"23.03.2024" swaggertype:"primitive,string"`
		Lyrics      string           `json:"text" example:"Oh baby baby baby ooh\n\nOh baby baby"`
		Link        string           `json:"link" example:"https://nescam.io"`
		Name        string           `json:"name" example:"SomeSongUsername"`
		Group       string           `json:"group" example:"Bring me the horizon"`
	}
)

// @Tags		songs
// @Accept		json
// @Produce	json
// @Router		/songs [post]
// @Param		input	body		CreateSongRequest	true	"Body"
// @Success	200		{object}	msongs.Song
// @Failure 400 {object} responses.Response
func (c *SongController) Create(ctx *gin.Context) {
	var req CreateSongRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.Error(err))
	}
	mod := &msongs.Song{
		ReleaseDate: req.ReleaseDate.Time,
		Link:        req.Link,
		Lyrics:      req.Lyrics,
		Name:        req.Name,
		Group:       req.Group,
	}
	res, err := c.service.Create(mod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// @Tags songs
// @Produce json
// @Router /songs [get]
// @Param input query GetSongsRequest true "Body"
// @Success 200 {array} msongs.Song
func (c *SongController) Get(ctx *gin.Context) {

	min, max, name, group, offset :=
		ctx.Query("release_date_min"),
		ctx.Query("release_date_max"),
		ctx.Query("name"),
		ctx.Query("query"),
		ctx.Query("offset")
	var minT, maxT types.CustomTime
	minT.UnmarshalJSON([]byte(min))
	maxT.UnmarshalJSON([]byte(max))
	offsetInt, _ := strconv.Atoi(offset)
	songs, err := c.service.GetAll(
		minT.Time,
		maxT.Time,
		name,
		group,
		offsetInt,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, songs)
}

// @Tags songs
// @Accepts json
// @Produce json
// @Param id path int true "Song id"
// @Param couplet_offset query int false "Couplet offset"
// @Success 200 {object} msongs.Song
// @Router /songs/{id} [get]
func (c *SongController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")
	offset := ctx.Query("couplet_offset")
	offsetInt, _ := strconv.Atoi(offset)

	song, err := c.service.GetById(id, offsetInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, song)
}

// @Tags songs
// @Accept json
// @Produce json
// @Success 200 {object} msongs.Song
// @Param id path int true "Song id"
// @Param input body UpdateSongRequest true "input"
// @Router /songs/{id} [put]
func (c *SongController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateSongRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.Error(err))
	}
	mod := &msongs.Song{
		ID:          id,
		ReleaseDate: req.ReleaseDate.Time,
		Link:        req.Link,
		Lyrics:      req.Lyrics,
		Name:        req.Name,
		Group:       req.Group,
	}
	response, err := c.service.Update(id, mod)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Tags songs
// @Produce json
// @Success 200 {object} msongs.Song
// @Param id path int true "Song id"
// @Router /songs/{id} [delete]
func (c *SongController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	song, err := c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.Error(err))
	}
	ctx.JSON(http.StatusOK, song)
}
