package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Genre will serve all genre needed
type Genre struct {
	gorm.Model
	Name string `json:"name"`
}
