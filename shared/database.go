package shared

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("lms.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&User{}, &Session{})
	pass_, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	DB.Create(&User{Username: "user", Password: string(pass_)})
}