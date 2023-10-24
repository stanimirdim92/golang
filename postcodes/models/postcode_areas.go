package models

import (
	"postcodes/config"
)

type PostcodeAreas struct {
	Id   uint   `json:"id" gorm:"primary_key, AUTO_INCREMENT"`
	Name string `json:"name"`
}

func GetPostcodeAreas() (postcodeAreass []PostcodeAreas, err error) {
	var postcodeAreas []PostcodeAreas

	result := config.DB.Find(&postcodeAreas)

	return postcodeAreas, result.Error
}
