package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func openDataBase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "events.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}
	return db
}

//var err error

func InitialMigration() {
	db, err := gorm.Open("sqlite3", "events.db")
	db.DropTableIfExists(&event{}, &User{}, &Token{}, &ticket{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&event{}, &ticket{}, User{}, Token{})

}
