package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type databaseInfo struct {
	uname string
	psw   string
	name  string
	host  string
}

var dbConn *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv error: %v", err)
	}

	dbInfo := databaseInfo{
		uname: os.Getenv("db_user"),
		psw:   os.Getenv("db_pass"),
		name:  os.Getenv("db_name"),
		host:  os.Getenv("db_host"),
	}

	engine := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbInfo.host, dbInfo.uname, dbInfo.name, dbInfo.psw)

	conn, err := gorm.Open("postgres", engine)
	if err != nil {
		log.Fatalf("db conn error: %v", err)
	}

	dbConn = conn
	dbConn.Debug().AutoMigrate(&User{})
}

func DB() *gorm.DB {
	return dbConn
}
