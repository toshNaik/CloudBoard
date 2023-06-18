package initializers

import (
	"github.com/toshnaik/CloudBoard/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	
}