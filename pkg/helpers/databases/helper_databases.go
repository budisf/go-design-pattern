package helperDatabases

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type IHelperDatabases interface {
	PaginationPostgresSQL(structResultQuery interface{}, rawQuery *string, structPagination *QueryParamPaginationEntity, fieldToSearchString []string, sqlOrder [1]string) (interface{}, error)
}

type helperDatabases struct {
	conn *gorm.DB
}

func InitHelperDatabase(conn *gorm.DB) IHelperDatabases {
	return &helperDatabases{
		conn: conn,
	}
}

/*
					|--------------------------------------------------------------------------
					| pagination with custom raw query
					|--------------------------------------------------------------------------
					|
					| This function is for return data from complex query and combine pagination, search string
				    | and custom struct receive
					|
					| @database purpose: PostgresSQL
					| @variable explanation:
					| - structResultQuery interface{} : for receive custom struct from need Scanner SQL from GORM and final result it is.
			        | - rawQuery *string: complex query whose passed off on inner params function
		            | - structPagination *QueryParamPaginationEntity: is struct include many param include there: page, limit, offset, search
					| - fieldToSearchString []string: this function only receive list of fields as reference from another func who passed it.
	                | - sqlOrder [1]strings: only single list order field

|
*/
func (h *helperDatabases) PaginationPostgresSQL(structResultQuery interface{}, rawQuery *string, structPagination *QueryParamPaginationEntity, fieldToSearchString []string, sqlOrder [1]string) (interface{}, error) {
	// initiate variable
	var newQueryLikeField string
	var checkStruct bool
	var messageError error
	var searchQueryString string
	var valueForSearchQueryLike string
	var orderBy string
	/*
		check type data for scan result query on gorm
	*/
	switch structResultQuery.(type) {
	case nil:
		// type of i is type of x (interface{})
		checkStruct = false
		messageError = errors.New("not throw result query is nil")
	case int:
		// type of i is int
		checkStruct = false
		messageError = errors.New("not throw result query type of is int")
	case float64:
		// type of i is float64
		checkStruct = false
		messageError = errors.New("not throw result query type of is int")
	case func(int) float64:
		// type of i is func(int) float64
		checkStruct = false
		messageError = errors.New("not throw result query type of is func(int) float64")
	case func(int) int64:
		// type of i is func(int) int64
		checkStruct = false
		messageError = errors.New("not throw result query type of is func(int) int64")
	case bool, string:
		// type is bool or string
		checkStruct = false
		messageError = errors.New("not throw result query type of is bool or string")
	default:
		// type of i is type of x (interface{}) or struct
		checkStruct = true
		messageError = nil
	}
	/*
		to create new string for query LIKE
	*/
	var tempSearchString []string
	if len(fieldToSearchString) > 0 {
		for i := 0; i < len(fieldToSearchString); i++ {
			newQueryLikeField = "COALESCE(" + fieldToSearchString[i] + ",'')"
			tempSearchString = append(tempSearchString, newQueryLikeField)
		}
	}
	// full count
	if len(tempSearchString) > 0 {
		searchQueryString = " where lower(concat(" + strings.Join(tempSearchString, ",' ',") + ")) like "
		valueForSearchQueryLike = " '%" + strings.ToLower(*structPagination.Search) + "%' "
	} else {
		searchQueryString = " where 1 = "
		valueForSearchQueryLike = "1"
	}

	limitString := strconv.Itoa(int(*structPagination.Limit))
	offsetString := strconv.Itoa(int(*structPagination.Offset))
	limitOffsetQuery := " LIMIT " + limitString + " OFFSET " + offsetString

	fmt.Println("DEBUG WHERE - 1: ", searchQueryString)
	fmt.Println("DEBUG WHERE - 2: ", valueForSearchQueryLike)

	if structPagination.OrderBy != nil && *structPagination.OrderBy == "desc" {
		orderBy = " DESC "
	} else {
		orderBy = " ASC "
	}

	sqlScript := "SELECT " +
		"raw.* FROM ( " +
		" SELECT ( " +
		"SELECT COUNT(1) FROM ( " +
		*rawQuery + searchQueryString + valueForSearchQueryLike +
		") as raw" +
		" ) AS full_count, raw.* FROM ( " +
		*rawQuery + searchQueryString + valueForSearchQueryLike + "ORDER BY " + sqlOrder[0] + orderBy +
		") as raw " +
		") AS raw " + limitOffsetQuery

	fmt.Println("DEBUG QUERY: ", sqlScript)

	if checkStruct {
		err := h.conn.Raw(sqlScript).Scan(&structResultQuery).Error
		if err != nil {
			messageError = err
		}
	} else {
		messageError = errors.New("something goes wrong when check struct")
	}

	return structResultQuery, messageError
}
