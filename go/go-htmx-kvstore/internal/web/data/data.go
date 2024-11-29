package data

import "gorm.io/gorm"

type KeyValue struct {
	gorm.Model
	Key    string
	Value  string
	User   User
	UserID uint
}

type User struct {
	gorm.Model
	Username string
	Password string
	Token    string
}
