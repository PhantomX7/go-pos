package mysql_test

import (
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
)

const databaseName = "cleanarch"

func SetupTests() *gorm.DB {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	db, err := gorm.Open(mocket.DriverName, databaseName)

	if err != nil {
		panic(err)
	}

	return db
}

