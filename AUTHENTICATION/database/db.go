package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"package/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *sql.DB

//GetDatabase is return db connection
func GetDatabase() *gorm.DB {
	databaseurl := os.Getenv("DATABASE_URL")
	databasename := os.Getenv("DATABASE_NAME")
	connection, err := gorm.Open(databasename, databaseurl)
	if err != nil {
		log.Fatalln("wrong database url")
	}
	sqldb := connection.DB()
	err = sqldb.Ping()
	if err != nil {
		log.Fatal("database is disconnected")
	}
	fmt.Println("connected to database")
	return connection
}

//Closedatabase is close database
func Closedatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}

//Initialmigration is migrate model to table
func Initialmigration() {
	connection := GetDatabase()
	connection.AutoMigrate(&model.User{})
	defer Closedatabase(connection)
	fmt.Println("migration done")
}
