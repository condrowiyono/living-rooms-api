package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Country will serve all country needed, including language and origin country
type Country struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}
