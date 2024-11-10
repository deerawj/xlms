package main

import "gorm.io/gorm"

// User is a struct that represents a user in the database
type User struct {
    gorm.Model
    ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    Username string `json:"username" gorm:"unique"`
    Password string `json:"password"`
}

type Session struct {
    gorm.Model
    ID      uint   `json:"id" gorm:"primaryKey,autoIncrement"`
    UserID  uint   `json:"user_id"` 
    Token   string `json:"token"`

    User    User   `gorm:"foreignKey:UserID"`
}