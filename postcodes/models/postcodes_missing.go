package models

import (
	"database/sql"
	"gorm.io/gorm/clause"
	"postcodes/config"
	"time"
)

type PostcodesMissing struct {
	Postcode    string `json:"postcode"`
	TimeCreated int64  `json:"time_created"`
}

/**
 * @receiver b
 * @return string
 */
func TableName() string {
	return "postcodes_missing"
}

/**
 * @param id int
 * @return postcodes
 * @return err
 */
func CheckPostcodeIsMissing(pc string) (isMissing bool, err error) {
	var postcodeIsMissing int64

	result := config.DB.Table(TableName()).Select("postcode").Where(&PostcodesMissing{Postcode: pc}).Count(&postcodeIsMissing)

	if postcodeIsMissing > 0 {
		return true, result.Error
	}

	return false, result.Error
}

/**
 * @param id int
 * @return postcodes
 * @return err
 */
func StorePostcodeIsMissing(pc string) (err error) {
	result := config.DB.Table(TableName()).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&PostcodesMissing{Postcode: pc, TimeCreated: time.Now().Unix()})

	return result.Error
}

/**
 * @param id int
 * @return postcodes
 * @return err
 */
func DeletePostcodeIsMissing(pc string) (err error) {
	result := config.DB.Table(TableName()).
		Where("postcode = @pc", sql.Named("pc", pc)).
		Delete(&PostcodesMissing{Postcode: pc})

	return result.Error
}
