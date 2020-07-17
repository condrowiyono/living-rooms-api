package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Production Studio that produce
type Production struct {
	gorm.Model
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}
