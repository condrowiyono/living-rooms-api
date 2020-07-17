package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(
		&User{},
		&Person{},
		&Image{},
		&Genre{},
		&Country{},
		&Movie{},
		&Player{},
		&Video{},
		&Production{},
		&Artist{},
		&Concert{},
		&Network{},
		&TvCreator{},
		&TvEpisode{},
		&TvSeason{},
		&Tv{},
	)
	return db
}
