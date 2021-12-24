package models

import (
	"os/user"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func InitDB() (*gorm.DB, error) {
	var err error

	//Get Current User
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	filepath := usr.HomeDir + "/gomodoro.db" // file path in user home directory

	db, err := gorm.Open("sqlite3", filepath) //Create DB connection

	if err != nil {
		return nil, err
	}

	db.Exec(`PRAGMA foreign_keys=ON`) // Need this to use foreign keys on sqlite.

	db.AutoMigrate(&Task{}, &SubTask{}) //To mmigrate models to the db

	return db, nil
}
