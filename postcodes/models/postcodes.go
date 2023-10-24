package models

import (
	"database/sql"
	"postcodes/config"
	"strings"
)

type Postcode struct {
	ID                      uint    `json:"id" gorm:"primaryKey, autoIncrement, column:id"`
	Postcode                string  `json:"postcode" gorm:"column:postcode"`
	PostcodeInward          string  `json:"postcode_inward" gorm:"column:postcode_inward"`
	PostcodeOutward         string  `json:"postcode_outward" gorm:"column:postcode_outward"`
	PostTown                string  `json:"post_town" gorm:"column:post_town"`
	DependantLocality       string  `json:"dependant_locality" gorm:"column:dependant_locality"`
	DoubleDependantLocality string  `json:"double_dependant_locality" gorm:"column:double_dependant_locality"`
	Thoroughfare            string  `json:"thoroughfare" gorm:"column:thoroughfare"`
	DependantThoroughfare   string  `json:"dependant_thoroughfare" gorm:"column:dependant_thoroughfare"`
	BuildingNumber          string  `json:"building_number" gorm:"column:building_number"`
	BuildingName            string  `json:"building_name" gorm:"column:building_name"`
	SubBuildingName         string  `json:"sub_building_name" gorm:"column:sub_building_name"`
	PoBox                   string  `json:"po_box" gorm:"column:po_box"`
	DepartmentName          string  `json:"department_name" gorm:"column:department_name"`
	OrganisationName        string  `json:"organisation_name" gorm:"column:organisation_name"`
	Udprn                   string  `json:"udprn" gorm:"column:udprn"`
	Umprn                   string  `json:"umprn" gorm:"column:umprn"`
	PostcodeType            string  `json:"postcode_type" gorm:"column:postcode_type"`
	SuOrganisationIndicator string  `json:"su_organisation_indicator" gorm:"column:su_organisation_indicator"`
	DeliveryPointSuffix     string  `json:"delivery_point_suffix" gorm:"column:delivery_point_suffix"`
	Line1                   string  `json:"line_1" gorm:"column:line_1"`
	Line2                   string  `json:"line_2" gorm:"column:line_2"`
	Line3                   string  `json:"line_3" gorm:"column:line_3"`
	Longitude               float64 `json:"longitude" gorm:"column:longitude"`
	Latitude                float64 `json:"latitude" gorm:"column:latitude"`
	Eastings                float64 `json:"eastings" gorm:"column:eastings"`
	Northings               float64 `json:"northings" gorm:"column:northings"`
	Country                 string  `json:"country" gorm:"column:country"`
	TraditionalCounty       string  `json:"traditional_county" gorm:"column:traditional_county"`
	AdministrativeCounty    string  `json:"administrative_county" gorm:"column:administrative_county"`
	PostalCounty            string  `json:"postal_county" gorm:"column:postal_county"`
	County                  string  `json:"county" gorm:"column:county"`
	District                string  `json:"district" gorm:"column:district"`
	Ward                    string  `json:"ward" gorm:"column:ward"`
	Premise                 string  `json:"premise" gorm:"column:premise"`
	PostcodeArea            string  `json:"postcode_area" gorm:"column:postcode_area"`

	PostcodeStripped string `json:"postcode_stripped" gorm:"column:postcode_stripped"`
	IsFromIdeal      bool   `json:"is_from_ideal" gorm:"column:is_from_ideal"`
	IsUp             uint8  `json:"is_up" gorm:"column:is_up"`
}

/**
 * @param params
 * @return postcodes
 * @return results
 * @return err
 */
func GetPostcodes(params map[string]interface{}) (postcodes []Postcode, err error) {
	var postcode []Postcode
	db := config.DB
	if pageEnd, ok := params["page_end"]; ok && params != nil {
		pageEnd := pageEnd.(int64)
		pageStart := params["page_start"].(int64)
		if pageEnd > 0 {
			db = db.Offset(int(pageStart)).Limit(int(pageEnd))
		}
	}
	result := db.Debug().Find(&postcode)

	return postcode, result.Error
}

/**
 * @param id int
 * @return postcodes
 * @return err
 */
func GetPostcodeByID(id int) (postcodes []Postcode, err error) {
	var postcode []Postcode

	result := config.DB.First(&postcode, "id = ?", id)

	return postcode, result.Error
}

/**
 * @param params
 * @return postcodes
 * @return results
 * @return err
 */
func GetPostcode(pc string) (postcodes []Postcode, err error) {
	var postcode []Postcode
	db := config.DB

	db = db.Where("postcode = @pc or postcode_stripped = @pc2", sql.Named("pc", pc), sql.Named("pc2", strings.ReplaceAll(pc, " ", "")))

	result := db.Debug().Find(&postcode)

	return postcode, result.Error
}

func DeletePostcode(pc string) error {
	result := config.DB.
		Where("postcode = @pc or postcode_stripped = @pc2", sql.Named("pc", pc), sql.Named("pc2", strings.ReplaceAll(pc, " ", ""))).
		Delete(&Postcode{Postcode: pc, PostcodeStripped: pc})

	return result.Error
}

func StorePostcodes(postcode []Postcode) error {
	result := config.DB.CreateInBatches(postcode, 10)

	return result.Error
}
