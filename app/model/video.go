package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Video will add blopper, behind the scene, etc
type Video struct {
	gorm.Model
	Type     string `json:"type" gorm:"type:text"`
	Source   string `json:"source"`
	VideoURL string `json:"player_url" gorm:"type:text"`
	ShowID   int    `json:"show_id"`
	ShowType string `json:"show_type"`
}
