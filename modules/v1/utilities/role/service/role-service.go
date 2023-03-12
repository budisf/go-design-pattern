package service

import (
	"errors"
	model "ethical-be/modules/v1/utilities/role/model"
	"ethical-be/modules/v1/utilities/role/model/relations"
	"ethical-be/modules/v1/utilities/role/repository"
	userRepository "ethical-be/modules/v1/utilities/user/repository"
	apiresponse "ethical-be/pkg/api-response"
	helperDatabases "ethical-be/pkg/helpers/databases"
	"fmt"
	"strconv"
)

type IRoleService interface {
	Save(role *model.RoleRequestDTO) (uint, error)
	GetById(id *string) (uint, error, *relations.RolesRelation)
	GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity, userIDLogin *string) (uint, error, *apiresponse.Pagination)
	UpdateById(id *string, role *model.RoleRequestUpdateDTO) (uint, error, *model.Roles)
	DeleteById(id *string) (uint, error)
	UpdateChangeParentRole(idRole *string, idHeadRole *string) (uint, error, *relations.RolesRelation)
	GetChildPositionByUserId(userId string) ([]model.RoleResponseDTO, error)
	getChildPositionByRoleID(roleId uint) ([]model.RoleResponseDTO, error)
}

type roleService struct {
	roleRepository repository.IRoleRepository
	userRepository userRepository.IUserRepository
}

func InitRoleService(roleRepository repository.IRoleRepository, userRepository userRepository.IUserRepository) IRoleService {
	return &roleService{
		roleRepository: roleRepository,
		userRepository: userRepository,
	}
}

func (service *roleService) Save(role *model.RoleRequestDTO) (uint, error) {
	var saveData relations.RolesRelation

	if role.ParentId != nil {
		roleIdInt := int(*role.ParentId)
		roleIdString := strconv.Itoa(roleIdInt)
		httpStatusResponseGetById, errorMessageGetById, dataRole := service.roleRepository.GetByIdRelation(&roleIdString)

		if errorMessageGetById != nil || httpStatusResponseGetById == 500 {
			return 500, errors.New(errorMessageGetById.Error())
		}

		if dataRole.ID == nil {
			return 404, nil
		}

	}

	saveData.Name = role.Name
	saveData.Label = role.Label
	saveData.ParentId = role.ParentId
	httpStatusResponse, errorMessage := service.roleRepository.Save(&saveData.Roles)

	return httpStatusResponse, errorMessage
}

func (service *roleService) GetById(id *string) (uint, error, *relations.RolesRelation) {
	httpStatusResponseGetById, errorMessageGetById, dataRole := service.roleRepository.GetByIdRelation(id)

	if errorMessageGetById != nil || httpStatusResponseGetById == 500 {
		return 500, errors.New(errorMessageGetById.Error()), nil
	}

	if dataRole.ID == nil {
		return 404, nil, nil
	}

	return httpStatusResponseGetById, nil, dataRole
}

func (service *roleService) GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity, userIDLogin *string) (uint, error, *apiresponse.Pagination) {

	httpStatusGetRoleByID, messageError, dataRole := service.userRepository.GetById(userIDLogin)
	if httpStatusGetRoleByID == 500 && messageError != nil {
		return 500, messageError, nil
	} else if dataRole.ID == nil {
		return 404, errors.New(fmt.Sprintf("user id %v does'nt exist", userIDLogin)), nil
	}
	httpStatus, errorMessage, responseObject := service.roleRepository.GetAllByPaginate(structPagination, dataRole.RoleId, dataRole.Role.Name)
	if httpStatus == 500 && errorMessage != nil {
		return 500, errorMessage, nil
	}
	return httpStatus, errorMessage, responseObject
}

func (service *roleService) UpdateById(id *string, role *model.RoleRequestUpdateDTO) (uint, error, *model.Roles) {
	httpResponseStatus, errorMessage, dataRole := service.roleRepository.GetByIdRelation(id)
	if httpResponseStatus == 500 || errorMessage != nil {
		return 500, errorMessage, nil
	}
	if dataRole.ID == nil {
		return 404, nil, nil
	}

	dataRole.Label = role.Label
	dataRole.Name = role.Name
	httpResponseStatusUpdate, errorMessageUpdate, dataRoleUpdate := service.roleRepository.UpdateById(&dataRole.Roles)

	if httpResponseStatusUpdate == 500 || errorMessageUpdate != nil {
		return 500, errorMessageUpdate, nil
	}

	return 200, nil, dataRoleUpdate
}

func (service *roleService) DeleteById(id *string) (uint, error) {
	httpResponseStatus, errorMessage, dataRole := service.roleRepository.GetByIdRelation(id)
	if httpResponseStatus == 500 || errorMessage != nil {
		return 500, errorMessage
	}
	if dataRole.ID == nil {
		return 404, nil
	}

	httpResponseStatusUpdate, errorMessageUpdate, _ := service.roleRepository.DeleteById(&dataRole.Roles)

	if httpResponseStatusUpdate == 500 || errorMessageUpdate != nil {
		return 500, errorMessageUpdate
	}

	return 200, nil
}

func (service *roleService) UpdateChangeParentRole(idRole *string, idHeadRole *string) (uint, error, *relations.RolesRelation) {
	var uintIdHeadRole *uint

	if idHeadRole != nil {
		httpResponseStatusHead, errorMessageHead, dataRoleHead := service.roleRepository.GetByIdRelation(idHeadRole)
		if httpResponseStatusHead == 500 || errorMessageHead != nil {
			return 500, errorMessageHead, nil
		}
		if dataRoleHead.ID == nil {
			return 404, errors.New("ID Head Role doesn't exists"), nil
		}
		uint64IdHeadRole, _ := strconv.ParseUint(*idHeadRole, 10, 32)
		uintIdHeadRoleTemp := uint(uint64IdHeadRole)
		uintIdHeadRole = &uintIdHeadRoleTemp
	} else {
		uintIdHeadRole = nil
	}

	httpResponseStatus, errorMessage, dataRole := service.roleRepository.GetByIdRelation(idRole)
	if httpResponseStatus == 500 || errorMessage != nil {
		return 500, errorMessage, nil
	}
	if dataRole.ID == nil {
		return 404, errors.New("ID Role doesn't exists"), nil
	}

	dataRole.ParentId = uintIdHeadRole

	httpResponseStatusUpdate, errorMessageUpdate, _ := service.roleRepository.UpdateById(&dataRole.Roles)

	if httpResponseStatusUpdate == 500 || errorMessageUpdate != nil {
		return 500, errorMessageUpdate, nil
	}

	httpResponseStatus, errorMessage, dataRole = service.roleRepository.GetByIdRelation(idRole)
	if httpResponseStatus == 500 || errorMessage != nil {
		return 500, errorMessage, nil
	}

	return 200, nil, dataRole
}

func (service *roleService) GetChildPositionByUserId(userId string) ([]model.RoleResponseDTO, error) {

	_, errUser, user := service.userRepository.GetById(&userId)
	if errUser != nil {
		return nil, errUser
	}
	roleId := user.RoleId
	if roleId == nil {
		return nil, errors.New("404")
	}

	results, err := service.getChildPositionByRoleID(*roleId)
	return results, err

}

func (service *roleService) getChildPositionByRoleID(roleId uint) ([]model.RoleResponseDTO, error) {
	result, err := service.roleRepository.GetChildRole(roleId)

	var resultQuery []model.RoleResponseDTO

	if len(result) != 0 {
		for _, value := range result {
			resultQuery = append(resultQuery, model.RoleResponseDTO{
				RoleId:    value.ID,
				Name:      value.Name,
				Label:     value.Label,
				ParentId:  value.ParentId,
				CreatedAt: &value.CreatedAt,
				UpdatedAt: &value.UpdatedAt,
			})
			result2, _ := service.getChildPositionByRoleID(*value.ID)
			resultQuery = append(resultQuery, result2...)

		}
	}
	return resultQuery, err
}
