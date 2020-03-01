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
	Type string `json:"type"`
	Path string `json:"path" gorm:"type:text"`
}

// Genre will serve all genre needed
type Genre struct {
	gorm.Model
	Name string `json:"name"`
}

// Country will serve all country needed, including language and origin country
type Country struct {
	gorm.Model
	Iso639_1 string `json:"iso_639_1"`
	Name     string `json:"name"`
}

// Movie hold every component detail about a movie
// Movie is one of other media consumption suppored here
type Movie struct {
	gorm.Model
	Banners          []Image   `json:"banners" gorm:"many2many:movies_banners;"`
	Budget           int       `json:"budget"`
	Genres           []Genre   `json:"genres" gorm:"many2many:movies_genres;"`
	Homepage         string    `json:"homepage"`
	TmdbID           int       `json:"tmdb_id"`
	ImdbID           string    `json:"imdb_id"`
	OriginalLanguage string    `json:"original_language"`
	OriginalTitle    string    `json:"original_title"`
	Overview         string    `json:"overview" gorm:"type:text"`
	Popularity       float64   `json:"popularity"`
	Posters          []Image   `json:"posters" gorm:"many2many:movies_posters;"`
	ReleaseDate      string    `json:"release_date"`
	Revenue          int       `json:"revenue"`
	Runtime          int       `json:"runtime"`
	Languages        []Country `json:"language" gorm:"many2many:movies_languages;"`
	Status           string    `json:"status"`
	Tagline          string    `json:"tagline"`
	Title            string    `json:"title"`
	VoteAverage      float64   `json:"vote_average"`
	VoteCount        int       `json:"vote_count"`
	ImdbRating       int       `json:"imdb_rating"`
	Awards           string    `json:"awards"`
	Actors           []Person  `json:"actors" gorm:"many2many:movies_actors;"`
	CountryID        int       `json:"-"`
	Country          Country   `json:"original_country"`
	Crews            []Person  `json:"crews" gorm:"many2many:movies_crews;"`
	Rated            string    `json:"rated"`
	Player           Player    `json:"player" gorm:"polymorphic:Show;"`
}

// Player will act player for every media supported in this app.
// That include movie, tv series, anime, bal-balan, etc
type Player struct {
	gorm.Model
	Type      string `gorm:"type:text"`
	Source    string `json:"source"`
	PlayerURL string `json:"player_url" gorm:"type:text"`
	ShowID    int    `json:"show_id"`
	ShowType  string `json:"show_type"`
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
	)
	return db
}
