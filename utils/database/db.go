package database

import "github.com/jinzhu/gorm"

var db *gorm.DB

func InitDB(d *gorm.DB) {
	db = d
}

// gorm will create a clone of database for transactions
func BeginTransactions() *gorm.DB{
	return db.Begin()
}


