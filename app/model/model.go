package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// User will serve as struct model
type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null" json:"email"`
	Username string `gorm:"unique;not null" json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// SampleLink will serve as struct model
type SampleLink struct {
	gorm.Model
	ImdbID string `json:"imdb_id"`
	Name   string `json:"name"`
	Link   string `json:"link"`
}

// Playlist will serve as struct model
type Playlist struct {
	gorm.Model
	Name        string
	Description string
	Shows       []Show `gorm:"many2many:playlists_shows;"`
}

// Show will serve as struct model
type Show struct {
	gorm.Model
	Actors     string
	Awards     string
	BoxOffice  string
	Country    string
	DVD        string
	Director   string
	Genre      string
	Language   string
	Plot       string `gorm:"type:text"`
	Poster     string
	Production string
	Rated      string
	Released   string
	Response   string
	Runtime    string
	Title      string
	Type       string
	Website    string
	Writer     string
	Year       string
	ImdbID     string `json:"imdbID"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	Banner     string
	Trailer    string
	PlayerURL  string `gorm:"type:text"`
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(
		&SampleLink{},
		&Playlist{},
		&Show{},
		&User{},
	)
	return db
}
