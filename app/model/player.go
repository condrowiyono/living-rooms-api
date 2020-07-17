package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

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
