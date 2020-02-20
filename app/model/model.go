package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// SampleLink will serve as struct model
type SampleLink struct {
	gorm.Model
	ImdbID string `json:"imdb_id"`
	Name   string `json:"name"`
	Link   string `json:"link"`
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&SampleLink{})
	return db
}
