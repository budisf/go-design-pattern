package helpers

import (
	"errors"
	"ethical-be/app/config"
	res "ethical-be/pkg/api-response"
	helperDatabases "ethical-be/pkg/helpers/databases"
	singleton "ethical-be/pkg/value-object"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

var (
	conf, err = config.Init()
)

/*
   |--------------------------------------------------------------------------
   | Convert Type Data Interface to Integer
   |--------------------------------------------------------------------------
   |
   | This function is for convert data type interface to integer,
   | this function only can convert interface (float64 and string) to int,
   | You can add more data type if you need it.
*/

func ConvertInterfaceToInteger(data interface{}) int {

	var result int
	switch data.(type) {
	case float64:
		toInt, ok := data.(float64)
		result = int(toInt)
		fmt.Println(toInt, ok)
	default:
		toInt, ok := data.(string)
		result, _ = strconv.Atoi(toInt)
		fmt.Println(toInt, ok)
	}
	return result
}

/*
	|--------------------------------------------------------------------------
	| to print logging error on the terminal
	|--------------------------------------------------------------------------
	|
	| This function is return detail location error with line, time, anda message

|
*/
func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, _, line, _ := runtime.Caller(1)

		log.SetFlags(log.Lmsgprefix)
		log.Printf("------------------------------------------------------------------------------------\n[Error]\t\t: %s\n[Line]\t\t: %d\n[Time]\t\t: %s\n[Message]\t: %v", runtime.FuncForPC(pc).Name(), line, time.Now().Format(time.RFC850), err)
		// PostToWebHook(err)
		b = true
	}
	return
}

/*
|--------------------------------------------------------------------------
| to remove all whitespace
|--------------------------------------------------------------------------
|
| This function is for return new string from remove whitespace on string/ setence
*/
func RemoveWhiteSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

/*
|--------------------------------------------------------------------------
| to remove all special character on query string
|--------------------------------------------------------------------------
|
| This function is for return new string from query string to avoid XSS and sql injection
|
*/
func Escape(letter *string) *string {
	var newString string
	re := regexp.MustCompile(`[%\w\d\s]+`)

	for _, match := range re.FindAllString(*letter, -1) {
		if len(match) > 1 {
			newString += match
		} else {
			newString = match
		}
	}
	re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	finalString := re_leadclose_whtsp.ReplaceAllString(newString, "")
	finalString = re_inside_whtsp.ReplaceAllString(finalString, " ")
	return &finalString
}

func QueryParamPaginateTransform(ctx *gin.Context) (*helperDatabases.QueryParamPaginationEntity, error) {
	var result helperDatabases.QueryParamPaginationDTO
	var resultEntity helperDatabases.QueryParamPaginationEntity
	var offset int = *singleton.PaginationValueObject.Offset
	var page int = *singleton.PaginationValueObject.Page
	var limit int = *singleton.PaginationValueObject.Limit
	var search string = *singleton.PaginationValueObject.Search
	var orderBy string = *singleton.PaginationValueObject.OrderBy

	if err := ctx.ShouldBindQuery(&result); err != nil {
		return nil, err
	}
	if result.Page != nil && *result.Page < 1 {
		return nil, errors.New("Page query params not allowed under 1")
	} else {
		// {endpoint}?page=1&limit=1&search=fulan

		if result.OrderBy != nil {
			orderBy = *result.OrderBy
		}

		if result.Limit != nil {
			limit = *result.Limit
		}

		if result.Page != nil {
			page = *result.Page
		}

		if result.Search != nil {
			search = *Escape(result.Search)
		}

		if result.Offset != nil {
			offset = ((page - 1) * limit) + *result.Offset
		} else {
			offset = ((page - 1) * limit)
		}

		resultEntity = helperDatabases.QueryParamPaginationEntity{
			Page:    &page,
			Offset:  &offset,
			Limit:   &limit,
			Search:  &search,
			OrderBy: &orderBy,
		}
	}

	return &resultEntity, nil
}

func StringToUint64(data *string) (*uint64, error) {
	newData, err := strconv.ParseUint(*data, 10, 32)
	return &newData, err
}

func PaginationMetadata(count *int64, limit int, page *int, endpoint string) helperDatabases.ResponseBackPaginationDTO {

	var totalPage float64
	previousPage := *page - 1
	nextPage := *page + 1
	totalData := int(*count)
	totalPage1 := float64(totalData) / float64(limit)
	totalPage1 = math.Ceil(totalPage1)

	if totalPage1 == 0 {
		totalPage = 1
	} else {
		totalPage = totalPage1
	}

	nextPageString := fmt.Sprintf("%s/%s/%spage=%d&limit=%d", conf.App.Url, conf.App.Name_api, endpoint, nextPage, limit)
	previousPageUrlString := fmt.Sprintf("%s/%s/%spage=%d&limit=%d", conf.App.Url, conf.App.Name_api, endpoint, previousPage, limit)
	firstPageUrlString := fmt.Sprintf("%s/%s/%spage=%d&limit=%d", conf.App.Url, conf.App.Name_api, endpoint, 1, limit)
	lastPageUrlString := fmt.Sprintf("%s/%s/%spage=%v&limit=%d", conf.App.Url, conf.App.Name_api, endpoint, totalPage, limit)

	results := helperDatabases.ResponseBackPaginationDTO{
		TotalData:        &totalData,
		TotalDataPerPage: &limit,
		CurrentPage:      page,
		PreviousPage:     &previousPage,
		TotalPage:        &totalPage,
		NextPageUrl:      &nextPageString,
		PreviousPageUrl:  &previousPageUrlString,
		FirstPageUrl:     &firstPageUrlString,
		LastPageUrl:      &lastPageUrlString,
	}

	return results

}

/*
|--------------------------------------------------------------------------
| to convert utc to time.Time with zone Asia/Jakarta
|--------------------------------------------------------------------------
|
| This function is for return new string from query string to avoid XSS and sql injection
|
*/
var countryTz = map[string]string{
	"Hungary": "Europe/Budapest",
	"Egypt":   "Africa/Cairo",
	"Jakarta": "Asia/Jakarta",
}

func UTCTotTime(utc *uint, countryName string) (*time.Time, error) {
	loc, err := time.LoadLocation(countryTz[countryName])
	if err != nil {
		return nil, err
	}
	tm := time.Unix(int64(*utc), 0).In(loc)
	// unitTimeInRFC3339 := tm.Format(time.RFC3339)
	return &tm, nil
}

func GetUserIdFromMiddleware(c *gin.Context) int {
	// c.Get("user_id") will return uint based on model.Users
	// and from repository function we use int therefore
	// we convert to int
	uidValue, exists := c.Get("user_id")
	if exists != true {
		c.JSON(http.StatusInternalServerError, res.ServerError("UserId authentication requested but middleware not registered on routes"))
	}
	uidString := fmt.Sprintf("%v", uidValue)
	uidInt, _ := strconv.Atoi(uidString)

	return uidInt
}

func GetUserRoleIdFromMiddleware(c *gin.Context) int {
	// c.Get("user_id") will return uint based on model.Users
	// and from repository function we use int therefore
	// we convert to int
	uidValue, exists := c.Get("user_role_id")
	if exists != true {
		c.JSON(http.StatusInternalServerError, res.ServerError("UserRoleId authentication requested but middleware not registered on routes"))
	}
	uidString := fmt.Sprintf("%v", uidValue)
	uidInt, _ := strconv.Atoi(uidString)

	return uidInt
}

/*
|--------------------------------------------------------------------------
| GetThisMondayAndSunday
|--------------------------------------------------------------------------
|
| Get date of monday and sunday on current week
*/
func GetThisMondayAndSunday() (time.Time, time.Time) {
	/*
		Get date of monday on current week
	*/
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, offset)
	sunday := monday.AddDate(0, 0, 6)
	return monday, sunday
}

/*
|--------------------------------------------------------------------------
| GetMonthAndYearFromUnix
|--------------------------------------------------------------------------
|
| This function is for get month and year and convert to integer from unix (epoch)
*/
func GetMonthAndYearFromUnix(date uint) (uint, uint) {
	convDate := time.Unix(int64(date), 0)

	month := uint(convDate.Month())
	year := uint(convDate.Year())

	return month, year
}

/*
|--------------------------------------------------------------------------
| convertInfToZero
|--------------------------------------------------------------------------
|
| to convert +Inf(angka tak terhingga) to zero
*/
func ConvertInfToZero(x float64) float64 {
	if math.IsInf(x, 1) {
		return math.Copysign(0, x)
	}
	return x
}

/*
	|--------------------------------------------------------------------------
	| getoffsetpagination
	|--------------------------------------------------------------------------
	|
	| to get offset for pagination purposes
*/

func GetOffsetPagination(page int, limit int) int {
	pageA := page - 1
	if page == 0 {
		pageA = 0
	}
	offset := pageA * limit

	return offset
}

/*
	|--------------------------------------------------------------------------
	| queryparampagination
	|--------------------------------------------------------------------------
	|
	| to get page and limit from query param
*/

func QueryParamPagination(c *gin.Context) (*int, *int, error) {
	pageTemp := c.Query("page")
	limitTemp := c.Query("limit")

	if pageTemp == "" && limitTemp == "" {
		return nil, nil, errors.New("Page and limit required")
	}

	page, errPage := strconv.Atoi(pageTemp)
	if errPage != nil {
		return nil, nil, errPage
	}

	limit, errLimit := strconv.Atoi(limitTemp)
	if errLimit != nil {
		return nil, nil, errLimit
	}

	return &page, &limit, nil
}

/*
	|--------------------------------------------------------------------------
	| response
	|--------------------------------------------------------------------------
	|
	| status code responses
*/

func ResponseError(statusCode int, err error, c gin.Context) {
	if err != nil {
		if statusCode == 404 {
			c.JSON(http.StatusNotFound, res.NotFound(err.Error()))
			return
		}
		if statusCode == 403 {
			c.JSON(http.StatusForbidden, res.StatusForbidden(err.Error()))
			return
		}
		if statusCode == 400 {
			c.JSON(http.StatusBadRequest, res.BadRequest(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, res.ServerError(err.Error()))
		return
	}
}

/*
	|--------------------------------------------------------------------------
	| RangeMonth
	|--------------------------------------------------------------------------
	|
	| gets the difference of month
*/

func RangeMonth(start time.Time, end time.Time) int {
	year := end.Year() - start.Year()                    //2023 - 2022 = 1
	month := (int(end.Month()) - int(start.Month()) + 1) //2 - 12 + 1 = -9
	diff := (year * 12) + month                          //12 + (-9) = 3

	return diff
}
