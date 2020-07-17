package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Artist s
type Artist struct {
	gorm.Model
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Picture string `json:"picture"`
}

// Concert s
type Concert struct {
	gorm.Model
	Title       string  `json:"title"`
	ArtistID    int     `json:"-"`
	ConcertDate string  `json:"concert_date"`
	ReleaseDate string  `json:"release_date"`
	Place       string  `json:"place"`
	Overview    string  `json:"overview" gorm:"type:text"`
	Setlist     string  `json:"setlist" gorm:"type:text"`
	Artist      Artist  `json:"artist"`
	Player      Player  `json:"player" gorm:"polymorphic:Show;"`
	Videos      []Video `json:"videos" gorm:"polymorphic:Show;"`
	Banners     []Image `json:"banners" gorm:"many2many:concerts_banners;"`
}
