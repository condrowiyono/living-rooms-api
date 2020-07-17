package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Movie hold every component detail about a movie
type Movie struct {
	gorm.Model
	TmdbID      int          `json:"tmdb_id"`
	ImdbID      string       `json:"imdb_id"`
	Overview    string       `json:"overview" gorm:"type:text"`
	ReleaseDate string       `json:"release_date"`
	Runtime     int          `json:"runtime"`
	Title       string       `json:"title"`
	Director    string       `json:"director"`
	Writer      string       `json:"writer"`
	Actors      []Person     `json:"actors" gorm:"many2many:movies_actors;association_autocreate:false;"`
	Productions []Production `json:"productions" gorm:"many2many:movies_productions;association_autocreate:false;"`
	Countries   []Country    `json:"countries" gorm:"many2many:movies_countries;association_autocreate:false;"`
	Banners     []Image      `json:"banners" gorm:"many2many:movies_banners;"`
	Posters     []Image      `json:"posters" gorm:"many2many:movies_posters;"`
	Genres      []Genre      `json:"genres" gorm:"many2many:movies_genres;association_autocreate:false;"`
	Player      Player       `json:"player" gorm:"polymorphic:Show;"`
	Videos      []Video      `json:"videos" gorm:"polymorphic:Show;"`
}
