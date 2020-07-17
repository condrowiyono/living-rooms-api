package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //asas
)

// Person will serve as struct model
// Person act in all aspect of human, including crew or cast, etc
type Person struct {
	gorm.Model
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Picture string `json:"picture"`
}
