package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"time"
)

var (
	DB         *gorm.DB
	err        error
	sqlitePath string
)

func Init() {
	home, _ := os.UserHomeDir()
	sqlitePath = home + "/.cy/cy.db"
	if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
		_ = os.MkdirAll(path.Dir(sqlitePath), os.ModePerm)
	} else {

	}
	DB, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	_ = DB.AutoMigrate(&AcmeAccount{}, &WebsiteSSL{})
}

type BaseModel struct {
	ID        uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
