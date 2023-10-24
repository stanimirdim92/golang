package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"postcodes/models"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
)

func Index(ctx echo.Context) error {

	req := ctx.Request()
	fmt.Println(req.Host)
	page, err := strconv.ParseInt(ctx.QueryParam("page"), 10, 64)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.ParseInt(ctx.QueryParam("page_size"), 10, 64)
	if err != nil {
		pageSize = 50
	}

	var postcodes []models.Postcode
	postcodes, err = models.GetPostcodes(map[string]interface{}{
		"page_start": (page - 1) * pageSize,
		"page_end":   pageSize,
	})

	message := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Postcodes Found",
	}

	if err != nil {
		message["code"] = http.StatusResetContent
		message["message"] = "No Postcodes Found"
		return formatJSON(ctx, nil, message, make(map[string]interface{}), 0, pageSize, page)
	}

	return formatJSON(ctx, postcodes, make(map[string]interface{}), message, 0, pageSize, page)
}

/**
 * @param c
 * @param items
 * @param errors
 * @param success
 * @param results
 * @param pages
 * @param page
 */
func formatJSON(ctx echo.Context, items interface{}, errors map[string]interface{}, success map[string]interface{}, results int64, pages int64, page int64) error {
	status := http.StatusOK
	if len(errors) != 0 {
		status = http.StatusNotFound
		if errors["code"] == nil {
			status = http.StatusNotFound
		}
	} else if len(success) != 0 {
		if success["code"] == nil {
			status = http.StatusOK
		}
	}

	return ctx.JSON(status, map[string]interface{}{
		"items":   items,
		"results": results,
		"pages":   pages,
		"page":    page,
		"errors":  errors,
		"success": success,
	})
}

func FindPostcode(ctx echo.Context) error {
	postcode := strings.ReplaceAll(ctx.Param("postcode"), "/", "")

	re := regexp.MustCompile(`/[^A-Za-z0-9]/`)
	postcode = re.ReplaceAllLiteralString(postcode, "")
	postcode = strings.ToUpper(postcode)

	message := map[string]interface{}{
		"code":    http.StatusResetContent,
		"message": "Postcodes not found",
	}

	isAlreadyMissing, _ := models.CheckPostcodeIsMissing(postcode)
	if isAlreadyMissing {
		return formatJSON(ctx, nil, message, make(map[string]interface{}), 0, 0, 0)
	}

	if len(postcode) < 5 || len(postcode) > 8 {
		_ = models.StorePostcodeIsMissing(postcode)
		return formatJSON(ctx, nil, message, make(map[string]interface{}), 0, 0, 0)
	}

	// first and last char must be letter and not a number
	if !unicode.IsLetter([]rune(postcode[:1])[0]) || !unicode.IsLetter([]rune(postcode[len(postcode)-1:])[0]) {
		_ = models.StorePostcodeIsMissing(postcode)
		return formatJSON(ctx, nil, message, make(map[string]interface{}), 0, 0, 0)
	}

	// must contain at least 1 letter and at least 2 numbers
	re = regexp.MustCompile(`/[A-Za-z].*[0-9]|[0-9].*[A-Za-z]/`)
	postcodeIsValid := re.MatchString(postcode)
	if !postcodeIsValid {
		re = regexp.MustCompile(`/[A-Za-z]/`)
		if len(re.ReplaceAllLiteralString(postcode, "")) == 0 {
			_ = models.StorePostcodeIsMissing(postcode)
			return formatJSON(ctx, nil, message, make(map[string]interface{}), 0, 0, 0)
		}
	}

	var postcodes []models.Postcode
	postcodes, err := models.GetPostcode(postcode)

	if err != nil {
		return formatJSON(ctx, nil, message, make(map[string]interface{}), 0, 0, 0)
	}

	if len(postcodes) != 0 {
		success := map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Postcodes found",
		}
		return formatJSON(ctx, postcodes, make(map[string]interface{}), success, 0, 0, 0)
	}

	if os.Getenv("IDEAL_POSTCODES") != "" && os.Getenv("IDEAL_POSTCODES_URL") != "" {
		_ = models.DeletePostcode(postcode)
		/**
		 * No postcodes found in local database.
		 * Get them from https://ideal-postcodes.co.uk
		 */
		idealPostcodes := getPostcodesFromIdealPostcodes(postcode)
		if len(idealPostcodes.Result) != 0 {
			_ = models.DeletePostcodeIsMissing(postcode)

			for key, postcodeElem := range idealPostcodes.Result {
				postcodeElem.IsFromIdeal = true
				postcodeElem.PostTown = strings.Title(strings.ToLower(postcodeElem.PostTown))
				postcodeElem.PostcodeStripped = strings.ReplaceAll(postcodeElem.Postcode, " ", "")

				idealPostcodes.Result[key] = postcodeElem
			}

			err := models.StorePostcodes(idealPostcodes.Result)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}

			postcodes, _ = models.GetPostcode(postcode)
		} else {
			_ = models.StorePostcodeIsMissing(postcode)
		}
	}

	if len(postcodes) == 0 {
		return formatJSON(ctx, nil, message, make(map[string]interface{}), 0, 0, 0)
	}

	message["code"] = http.StatusOK
	message["message"] = "Postcode found"

	return formatJSON(ctx, postcodes, make(map[string]interface{}), message, 0, 0, 0)
}

func getPostcodesFromIdealPostcodes(postcode string) APIResponse {
	var URL *url.URL
	URL, _ = url.Parse(os.Getenv("IDEAL_POSTCODES_URL") + postcode)

	parameters := url.Values{}
	parameters.Add("api_key", os.Getenv("IDEAL_POSTCODES"))
	URL.RawQuery = parameters.Encode()

	response, err := http.Get(URL.String())

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	contents, _ := ioutil.ReadAll(response.Body)
	var result APIResponse
	_ = json.Unmarshal(contents, &result)

	return result
}

type APIResponse struct {
	Message string            `json:"message"`
	Code    int64             `json:"code"`
	Result  []models.Postcode `json:"result"`
}
