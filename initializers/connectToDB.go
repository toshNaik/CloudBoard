package initializers

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("cboard.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
}