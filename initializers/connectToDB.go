package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	sqlDB, err := connectTCPSocket()
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}

	DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
}
