package main

import (
	"fmt"
	"github.com/subosito/gotenv"
	"log"
	"os"

	"github.com/PhantomX7/go-pos/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	if err := seedRoles(db); err != nil {
		panic(err)
	}

}

func seedRoles(db *gorm.DB) error {
	roles := []models.Role{
		models.Role{Name: "admin"},
		models.Role{Name: "guest"},
	}

	if !db.First(&models.Role{}).RecordNotFound() {
		return nil
	}

	for _, role := range roles {
		err := db.Create(&role).Error
		if err != nil {
			log.Println("error-seeding-role:", err)
			return err
		}
	}
	return nil
}
