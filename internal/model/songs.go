package model

import "time"

type Song struct {
	ID          string    `json:"id"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link"`
	Lyrics      string    `json:"lyrics"`
	Name        string    `json:"name"`
	Group       string    `json:"group"`
}
