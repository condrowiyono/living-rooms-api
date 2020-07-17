package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Image will serve all image needed
// Including Type banner, poster, profile_picture, general, etc
type Image struct {
	gorm.Model
	Type    string `json:"type"`
	Keyword string `json:"keyword"`
	Source  string `json:"source"`
	Path    string `json:"path" gorm:"type:text"`
}
