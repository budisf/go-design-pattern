package repository

import (
	"ethical-be/app/config"
	model "ethical-be/modules/v1/utilities/role/model"
	"ethical-be/modules/v1/utilities/role/model/relations"
	apiresponse "ethical-be/pkg/api-response"
	helperDatabases "ethical-be/pkg/helpers/databases"
	"fmt"
	"math"
	"strconv"

	"gorm.io/gorm"
)

type IRoleRepository interface {
	Save(role *model.Roles) (uint, error)
	GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity, roleID *uint, roleName *string) (uint, error, *apiresponse.Pagination)
	UpdateById(role *model.Roles) (uint, error, *model.Roles)
	DeleteById(role *model.Roles) (uint, error, *model.Roles)
	GetByIdRelation(id *string) (uint, error, *relations.RolesRelation)
	GetChildRole(roleID uint) ([]model.Roles, error)
}

type roleRepository struct {
	conf            *config.Conf
	conn            *gorm.DB
	helperDatabases helperDatabases.IHelperDatabases
}

func InitRoleRepository(conn *gorm.DB, helperDatabases helperDatabases.IHelperDatabases, conf *config.Conf) IRoleRepository {
	return &roleRepository{
		conf:            conf,
		conn:            conn,
		helperDatabases: helperDatabases,
	}
}

func (repo *roleRepository) Save(role *model.Roles) (uint, error) {
	if err := repo.conn.Create(&role).Error; err != nil {
		return 500, err
	}
	repo.conn.Preload("roles").Find(&role)
	return 200, nil
}

func (repo *roleRepository) GetById(id *string, roleID *string) (uint, error, *model.RolesRawQuerySelfJoinResult) {
	var role model.RolesRawQuerySelfJoinResult
	rawQuery := "select r.*, r2.name as parent_role_name from public.roles as r  " +
		"left join public.roles as r2 on r.parent_id = r2.id where true and " +
		" r.id = @id " +
		"order by 1 asc"

	err := repo.conn.Raw(rawQuery, map[string]interface{}{"id": *id}).Scan(&role).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &role
}

func (repo *roleRepository) GetByIdRelation(id *string) (uint, error, *relations.RolesRelation) {
	var RolesRelation relations.RolesRelation

	err := repo.conn.Preload("ParentRole").Where("id = ?", id).Find(&RolesRelation).
		Error

	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &RolesRelation
}

type ResultQueryWithFullCount struct {
	FullCount *int `json:"full_count"` // this properties must implement to throw on to struct global response meta data pagination (for get full count the data)
	// this expected result query from your raw query
	model.RolesWithEpochEntity
}

func (repo *roleRepository) GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity, roleID *uint, roleName *string) (uint, error, *apiresponse.Pagination) {
	var resultQueryWithFullCount []ResultQueryWithFullCount
	var resultQuery []model.RolesWithEpochEntity
	var totalCountData int
	var currentPage int
	var LastPage float64
	var nextPageUrl *string
	var previousPageUrl *string
	var firstPageUrl *string
	var lastPageUrl *string
	var queryWhereGradedRole string
	roleIDString := strconv.FormatUint(uint64(*roleID), 10)

	if roleName != nil {
		fmt.Println("QUERY DEBUG ROLE PAGINATION - Role Name ", *roleName)
		if *roleName == "super-admin" || *roleName == "trade-team" || *roleName == "director" || *roleName == "fic" || *roleName == "marketing-director" || *roleName == "msd" {
			queryWhereGradedRole = "rbp2.name <> 'super-admin' and rbp2.name <> 'trade-team' and rbp2.name <> 'director' and rbp2.name <> 'fic' and rbp2.name <> 'marketing-director' and rbp2.name <> 'msd' "
		} else {
			queryWhereGradedRole = "1=1 "
		}

	} else {
		fmt.Println("QUERY DEBUG ROLE PAGINATION - Role Name ", roleName)
		queryWhereGradedRole = "1=1 "
	}

	rawQuery := "select raw2.* from ( " +
		// @start template
		// your query please take this below. sub query (raw2) is wrapper which is to passed on func pagination helper
		"with recursive role_by_parent as ( " +
		"select r.* from roles as r where true and r.id = " + roleIDString + " " +
		"union all " +
		"select child.* from roles as child inner join role_by_parent as rbp on child.parent_id = rbp.id " +
		") " +
		"select rbp2.id as role_id, rbp2.name, rbp2.label, rbp2.parent_id, " +
		"round( cast(extract(epoch from rbp2.created_at) as numeric) ) as created_at, " +
		"round( cast(extract(epoch from rbp2.updated_at) as numeric) ) as updated_at " +
		"from role_by_parent as rbp2 where true and " + queryWhereGradedRole +
		// @end template
		" ) as raw2 "
	/*
		for fieldToSearchString
		{using name sub query itself}.{name_field}
	*/
	fieldToSearchString := []string{
		"raw2.name",
	}
	/*
		for fieldOrderAscBy
		{using name sub query itself}.{name_field}
	*/
	fieldOrderAscBy := [1]string{
		"raw2.role_id",
	}
	fmt.Println("resultQueryWithFullCount : ", resultQueryWithFullCount)
	resultRawQuery, errorResult := repo.helperDatabases.PaginationPostgresSQL(resultQueryWithFullCount, &rawQuery, structPagination, fieldToSearchString, fieldOrderAscBy)

	dataRole := resultRawQuery.([]ResultQueryWithFullCount)
	dataOfPage := int(int64(len(dataRole)))
	var previousPage int
	if *structPagination.Page == 1 {
		previousPage = 1
	} else {
		previousPage = int(int64(*structPagination.Page)) - 1
	}
	if len(dataRole) == 0 {
		resultQuery = []model.RolesWithEpochEntity{}
		totalCountData = 0
		currentPage = int(int64(*structPagination.Page))
		LastPage = 1
		nextPageUrl = nil
		previousPageUrl = nil
		firstPageUrlString := fmt.Sprintf("%s/%s/v1/role?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, previousPage, *structPagination.Limit)
		firstPageUrl = &firstPageUrlString
		lastPageUrlString := fmt.Sprintf("%s/%s/v1/role?page=%v&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, LastPage, *structPagination.Limit)
		lastPageUrl = &lastPageUrlString
	} else {
		for index, value := range dataRole {
			if index == 0 {
				totalCountData = *value.FullCount
			}
			if *value.RoleId != *roleID {
				resultQuery = append(resultQuery, model.RolesWithEpochEntity{
					RoleId:    value.RoleId,
					Name:      value.Name,
					ParentId:  value.ParentId,
					Label:     value.Label,
					CreatedAt: value.CreatedAt,
					UpdatedAt: value.UpdatedAt,
				})
			} else {
				resultQuery = []model.RolesWithEpochEntity{}
			}
		}
		currentPage = int(int64(*structPagination.Page))
		LastPage = math.Ceil(float64(totalCountData) / float64(*structPagination.Limit))
		nextPageString := fmt.Sprintf("%s/%s/v1/role?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, *structPagination.Page+1, *structPagination.Limit)
		nextPageUrl = &nextPageString
		previousPageUrlString := fmt.Sprintf("%s/%s/v1/role?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, previousPage, *structPagination.Limit)
		previousPageUrl = &previousPageUrlString
		firstPageUrlString := fmt.Sprintf("%s/%s/v1/role?page=%d&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, previousPage, *structPagination.Limit)
		firstPageUrl = &firstPageUrlString
		lastPageUrlString := fmt.Sprintf("%s/%s/v1/role?page=%v&limit=%d", repo.conf.App.Url, repo.conf.App.Name_api, LastPage, *structPagination.Limit)
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

func (repo *roleRepository) UpdateById(role *model.Roles) (uint, error, *model.Roles) {
	err := repo.conn.Save(&role).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, role
}

func (repo *roleRepository) DeleteById(role *model.Roles) (uint, error, *model.Roles) {
	err := repo.conn.Delete(&role).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, role
}

func (repo *roleRepository) GetChildRole(roleId uint) ([]model.Roles, error) {
	var role []model.Roles
	err := repo.conn.Where("deleted_at is null").Where("parent_id = ?", roleId).Find(&role).Error
	return role, err
}
