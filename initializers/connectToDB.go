package initializers

import (
	// redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

// var Redis *redis.Client
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
	// Redis = redis.NewClient(&redis.Options{
	// 	Addr: "10.135.2.75:6379",
	// 	Password: "ea90215a-9474-42b4-80ff-02c1a3312bbd",
	// 	DB: 0,
	// })
}
