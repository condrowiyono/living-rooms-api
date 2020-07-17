package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

type Network struct {
	gorm.Model
	Name    string `json:"name"`
	Country string `json:"country"`
}
