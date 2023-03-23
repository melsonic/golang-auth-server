package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {

	loaderr := godotenv.Load()

	if loaderr != nil {
		log.Fatal("Error loading .env file")
	}

	// load username, password, dbname
	dbuname := os.Getenv("DB_USERNAME")
	dbpword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_DBNAME")

	// mysql default port is 3306
	// dsn := "neilbenz:mockinbird@tcp(localhost:3306)/octerndb"
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbuname, dbpword, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// error connecting to db
	if err != nil {
		panic("Error connecting to db")
	}
	Db = db
	fmt.Println("database connected")

}
