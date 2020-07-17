package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// TvCreator adalah
type TvCreator struct {
	gorm.Model
	Name string `json:"name"`
}

type TvSeason struct {
	gorm.Model
	ReleaseDate  string      `json:"release_date"`
	EpisodeCount int         `json:"episode_count"`
	Name         string      `json:"name"`
	Overview     string      `json:"overview" gorm:"type:text"`
	SeasonNumber int         `json:"season_number"`
	Poster       string      `json:"poster"`
	Episodes     []TvEpisode `json:"episodes"`
}

type TvEpisode struct {
	gorm.Model
	AirDate       string `json:"air_date"`
	EpisodeNumber int    `json:"episode_number"`
	SeasonNumber  int    `json:"season_number"`
	Name          string `json:"name"`
	Overview      string `json:"overview" gorm:"type:text"`
	Still         string `json:"still_path"`
	Player        Player `json:"player" gorm:"polymorphic:Show;"`
}

// Tv hold every component detail about a tv show and drakor
type Tv struct {
	gorm.Model
	TmdbID       int          `json:"tmdb_id"`
	Overview     string       `json:"overview" gorm:"type:text"`
	ReleaseDate  string       `json:"release_date"`
	Runtime      int          `json:"runtime"`
	Name         string       `json:"name"`
	EpisodeCount int          `json:"episode_count"`
	SeasonCount  int          `json:"season_count"`
	Season       []TvSeason   `json:"seasons"`
	Posters      []Image      `json:"posters" gorm:"many2many:tv_posters;"`
	Banners      []Image      `json:"banners" gorm:"many2many:tv_banners;"`
	Countries    []Country    `json:"countries" gorm:"many2many:tv_countries;association_autocreate:false;"`
	Genres       []Genre      `json:"genres" gorm:"many2many:tv_genres;association_autocreate:false;"`
	Actors       []Person     `json:"actors" gorm:"many2many:tv_actors;association_autocreate:false;"`
	Network      []Network    `json:"networks" gorm:"many2many:tv_networks;association_autocreate:false;"`
	Creators     []TvCreator  `json:"creators" gorm:"many2many:tv_actors;association_autocreate:false;"`
	Productions  []Production `json:"productions" gorm:"many2many:tv_productions;association_autocreate:false;"`
}
