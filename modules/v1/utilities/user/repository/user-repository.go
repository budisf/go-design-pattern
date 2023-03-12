package repository

import (
	"ethical-be/app/config"
	"ethical-be/modules/v1/utilities/user/model"
	apiresponse "ethical-be/pkg/api-response"
	helperDatabases "ethical-be/pkg/helpers/databases"
	"fmt"
	"math"

	"gorm.io/gorm"
)

type IUserRepository interface {
	Save(user *model.Users) (uint, error)
	GetById(id *string) (uint, error, *model.Users)
	GetByNip(nip *string) (uint, error, *model.Users)
	GetByRoleId(roleId uint) (uint, error, *[]model.Users)
	UpdateById(user *model.Users) (uint, error, *model.Users)
	GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity) (uint, error, *apiresponse.Pagination)
	GetUserByAuthServerId(authUser int) (uint, error, *model.Users)
}

type userRepository struct {
	conf            *config.Conf
	conn            *gorm.DB
	helperDatabases helperDatabases.IHelperDatabases
}

func InitUserRepository(conn *gorm.DB, helperDatabases helperDatabases.IHelperDatabases, conf *config.Conf) IUserRepository {
	return &userRepository{
		conf:            conf,
		conn:            conn,
		helperDatabases: helperDatabases,
	}
}

func (repo *userRepository) Save(user *model.Users) (uint, error) {
	if err := repo.conn.Create(&user).Error; err != nil {
		return 500, err
	}
	repo.conn.Preload("users").Find(&user)
	return 200, nil
}

func (repo *userRepository) GetById(id *string) (uint, error, *model.Users) {
	var user model.Users
	err := repo.conn.Preload("Role").Where("status = 'active' and id = ? ", id).Find(&user).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &user
}

func (repo *userRepository) UpdateById(user *model.Users) (uint, error, *model.Users) {
	err := repo.conn.Save(&user).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, user
}

func (repo *userRepository) GetByNip(nip *string) (uint, error, *model.Users) {
	var user model.Users
	err := repo.conn.Where("status = 'active' and nip = ?", nip).Find(&user).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &user
}

func (repo *userRepository) GetByRoleId(roleId uint) (uint, error, *[]model.Users) {
	var user []model.Users
	err := repo.conn.Preload("Role").Where("status = 'active' and role_id = ?", roleId).Find(&user).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &user
}

type ResultQueryWithFullCount struct {
	FullCount *int `json:"full_count"` // this properties must implement to throw on to struct global response meta data pagination (for get full count the data)
	// this expected result query from your raw query
	model.UserResponseEpochDTO
}

func (repo *userRepository) GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity) (uint, error, *apiresponse.Pagination) {
	var resultQueryWithFullCount []ResultQueryWithFullCount
	var resultQuery []model.UserResponseEpochDTO
	var totalCountData int
	var currentPage int
	var LastPage float64
	var nextPageUrl *string
	var previousPageUrl *string
	var firstPageUrl *string
	var lastPageUrl *string

	rawQuery := "select raw2.* from ( " +
		// @start template
		// your query please take this below. sub query (raw2) is wrapper which is to passed on func pagination helper
		"select u.id, u.name, u.auth_server_id, u.nip, u.role_id,  round( cast(extract(epoch from u.created_at) as numeric) ) as created_at, round( cast(extract(epoch from u.updated_at) as numeric) ) as updated_at " +
		"from public.users as u  " +
		// @end template
		" ) as raw2 "
	/*
		for fieldToSearchString
		{using name sub query itself}.{name_field}
	*/
	fieldToSearchString := []string{
		"raw2.name", "raw2.nip",
	}
	/*
		for fieldOrderAscBy
		{using name sub query itself}.{name_field}
	*/
	fieldOrderAscBy := [1]string{
		"raw2.id",
	}
	resultRawQuery, errorResult := repo.helperDatabases.PaginationPostgresSQL(resultQueryWithFullCount, &rawQuery, structPagination, fieldToSearchString, fieldOrderAscBy)

	dataUser := resultRawQuery.([]ResultQueryWithFullCount)
	dataOfPage := int(int64(len(dataUser)))
	var previousPage int
	if *structPagination.Page == 1 {
		previousPage = 1
	} else {
		previousPage = int(int64(*structPagination.Page)) - 1
	}

	if len(dataUser) == 0 {
		resultQuery = []model.UserResponseEpochDTO{}
		totalCountData = 0
		currentPage = int(int64(*structPagination.Page))
		LastPage = 1
		nextPageUrl = nil
		previousPageUrl = nil
		firstPageUrlString := fmt.Sprintf("%s/%s/v1/user?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, previousPage, *structPagination.Limit)
		firstPageUrl = &firstPageUrlString
		lastPageUrlString := fmt.Sprintf("%s/%s/v1/user?page=%v&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, LastPage, *structPagination.Limit)
		lastPageUrl = &lastPageUrlString
	} else {
		for index, value := range dataUser {
			if index == 0 {
				totalCountData = *value.FullCount
			}
			resultQuery = append(resultQuery, model.UserResponseEpochDTO{
				ID:           value.ID,
				RoleId:       value.RoleId,
				Name:         value.Name,
				Nip:          value.Nip,
				AuthServerId: value.AuthServerId,
				CreatedAt:    value.CreatedAt,
				UpdatedAt:    value.UpdatedAt,
			})
		}
		currentPage = int(int64(*structPagination.Page))
		LastPage = math.Ceil(float64(totalCountData) / float64(*structPagination.Limit))
		nextPageString := fmt.Sprintf("%s/%s/v1/user?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, *structPagination.Page+1, *structPagination.Limit)
		nextPageUrl = &nextPageString
		previousPageUrlString := fmt.Sprintf("%s/%s/v1/user?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, previousPage, *structPagination.Limit)
		previousPageUrl = &previousPageUrlString
		firstPageUrlString := fmt.Sprintf("%s/%s/v1/user?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, previousPage, *structPagination.Limit)
		firstPageUrl = &firstPageUrlString
		lastPageUrlString := fmt.Sprintf("%s/%s/v1/user?page=%v&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, LastPage, *structPagination.Limit)
		lastPageUrl = &lastPageUrlString
	}
	responseBackJson := apiresponse.Pagination{
		MetaData: &helperDatabases.ResponseBackPaginationDTO{
			TotalData:        &totalCountData,
			TotalDataPerPage: &dataOfPage,
			CurrentPage:      &currentPage,
			PreviousPage:     &previousPage,
			TotalPage:        &LastPage,
			NextPageUrl:      nextPageUrl,
			PreviousPageUrl:  previousPageUrl,
			FirstPageUrl:     firstPageUrl,
			LastPageUrl:      lastPageUrl,
		},
		Records: resultQuery,
	}

	if errorResult != nil {
		return 500, errorResult, nil
	}

	return 200, nil, &responseBackJson
}

func (repo *userRepository) GetUserByAuthServerId(authUser int) (uint, error, *model.Users) {
	var user model.Users
	err := repo.conn.Where("status = 'active' and auth_server_id = ?", authUser).Find(&user).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &user
}
