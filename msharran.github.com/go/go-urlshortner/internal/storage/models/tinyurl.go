package models

import "gorm.io/gorm"

type TinyURL struct {
	gorm.Model
	Long  string
	Short string `gorm:"unique"`
}
