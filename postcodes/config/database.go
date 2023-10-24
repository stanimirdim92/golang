package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func dBUrl() string {
	DBPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		DBPort,
		os.Getenv("DB_DATABASE"),
	)
}

func InitDB() (db *gorm.DB) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dBUrl(),
		DefaultStringSize:         256,   // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  false, // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    false, // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   false, // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // smart configure based on used version
	}), &gorm.Config{})

	if db == nil {
		log.Fatal(err)
	}

	return db
}
