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
	Role     string `json:"role"`
	Password string `json:"password"`
}

// Playlist will serve as struct model
type Playlist struct {
	gorm.Model
	Name        string
	Description string
	Players     []Player `gorm:"many2many:playlists_players;"`
}

// Person will serve as struct model
// Person act in all aspect of human, including crew or cast, etc
type Person struct {
	gorm.Model
	TmdbID   int     `json:"tmdb_id"`
	Pictures []Image `json:"pictures" gorm:"many2many:persons_pictures;"`
	Name     string  `json:"name"`
}

// Image will serve all image needed
// Including Type banner, poster, profile_picture, general, etc
type Image struct {
	gorm.Model
	Type    string `json:"type"`
	Keyword string `json:"keyword"`
	Source  string `json:"source"`
	Path    string `json:"path" gorm:"type:text"`
}

// Genre will serve all genre needed
type Genre struct {
	gorm.Model
	Name string `json:"name"`
}

// Country will serve all country needed, including language and origin country
type Country struct {
	gorm.Model
	Name string `json:"name"`
}

type Production struct {
	gorm.Model
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

// Movie hold every component detail about a movie
// Movie is one of other media consumption suppored here
type Movie struct {
	gorm.Model
	Banners          []Image      `json:"banners" gorm:"many2many:movies_banners;"`
	Genres           []Genre      `json:"genres" gorm:"many2many:movies_genres;association_autocreate:false;"`
	TmdbID           int          `json:"tmdb_id"`
	ImdbID           string       `json:"imdb_id"`
	OriginalLanguage string       `json:"original_language"`
	Overview         string       `json:"overview" gorm:"type:text"`
	Posters          []Image      `json:"posters" gorm:"many2many:movies_posters;"`
	ReleaseDate      string       `json:"release_date"`
	Runtime          int          `json:"runtime"`
	Languages        []Country    `json:"language" gorm:"many2many:movies_languages;association_autocreate:false;"`
	Title            string       `json:"title"`
	VoteAverage      float64      `json:"vote_average"`
	ImdbRating       float64      `json:"imdb_rating"`
	Awards           string       `json:"awards"`
	Actors           []Person     `json:"actors" gorm:"many2many:movies_actors;association_autocreate:false;"`
	Productions      []Production `json:"productions" gorm:"many2many:movies_productions;association_autocreate:false;"`
	Countries        []Country    `json:"countries" gorm:"many2many:movies_countries;association_autocreate:false;"`
	Crews            []Person     `json:"crews" gorm:"many2many:movies_crews;association_autocreate:false;"`
	Director         string       `json:"director"`
	Rated            string       `json:"rated"`
	Player           Player       `json:"player" gorm:"polymorphic:Show;"`
	Videos           []Video      `json:"videos" gorm:"polymorphic:Show;"`
}

// Player will act player for every media supported in this app.
// That include movie, tv series, anime, bal-balan, etc
type Player struct {
	gorm.Model
	Type      string `json:"type" gorm:"type:text"`
	Source    string `json:"source"`
	PlayerURL string `json:"player_url" gorm:"type:text"`
	ShowID    int    `json:"show_id"`
	ShowType  string `json:"show_type"`
}

type Video struct {
	gorm.Model
	Type     string `json:"type" gorm:"type:text"`
	Source   string `json:"source"`
	VideoURL string `json:"player_url" gorm:"type:text"`
	ShowID   int    `json:"show_id"`
	ShowType string `json:"show_type"`
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(
		&Playlist{},
		&User{},
		&Person{},
		&Image{},
		&Genre{},
		&Country{},
		&Movie{},
		&Player{},
		&Video{},
		&Production{},
	)
	return db
}
