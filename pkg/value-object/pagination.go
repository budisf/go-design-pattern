package singleton

import (
	helperDatabases "ethical-be/pkg/helpers/databases"
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

var PaginationValueObject *helperDatabases.QueryParamPaginationDTO

var defaultPage int = 1       // default value page
var defaultLimit int = 100    //default value limit
var defaultOffset int = 0     //default value offset
var defaultSearch string = "" //default value search
var defaultOrderBy string = "asc"

func GetPaginationValueObject() *helperDatabases.QueryParamPaginationDTO {
	if PaginationValueObject == nil {
		lock.Lock()
		defer lock.Unlock()
		if PaginationValueObject == nil {
			fmt.Println("Initiate Pagination Default Value now.")
			PaginationValueObject = &helperDatabases.QueryParamPaginationDTO{
				Page:    func(i int) *int { return &i }(defaultPage),
				Limit:   func(i int) *int { return &i }(defaultLimit),
				Offset:  func(i int) *int { return &i }(defaultOffset),
				Search:  func(i string) *string { return &i }(defaultSearch),
				OrderBy: func(i string) *string { return &i }(defaultOrderBy),
			}
		} else {
			fmt.Println("Pagination instance already created.")
		}
	} else {
		fmt.Println("Pagination instance already created.")
	}

	return PaginationValueObject
}
