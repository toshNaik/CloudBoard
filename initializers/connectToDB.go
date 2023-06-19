package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/toshnaik/CloudBoard/utils"
)

var DB *gorm.DB

func ConnectToDB() {
	sqlDB, err := utils.ConnectTCPSocket()
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
