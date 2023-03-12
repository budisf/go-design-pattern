package service

import (
	"errors"
	repository2 "ethical-be/modules/v1/utilities/role/repository"
	"ethical-be/modules/v1/utilities/user/model"
	"ethical-be/modules/v1/utilities/user/repository"
	apiresponse "ethical-be/pkg/api-response"
	helperDatabases "ethical-be/pkg/helpers/databases"
	"fmt"
)

type IUserService interface {
	Save(user *model.UserRequestDTO) (uint, error)
	GetById(id *string) (uint, error, *model.Users)
	UpdateById(id *string, user *model.UserRequestUpdateDTO) (uint, error, *model.Users)
	GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity) (uint, error, *apiresponse.Pagination)
	DeleteById(id *string) (uint, error)
}

type userService struct {
	userRepository repository.IUserRepository
	roleRepository repository2.IRoleRepository
}

func InitUserRepository(userRepository repository.IUserRepository, roleRepository repository2.IRoleRepository) IUserService {
	return &userService{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

func (service *userService) Save(user *model.UserRequestDTO) (uint, error) {
	var saveData model.Users

	httpStatusCheckRole, errCheckRole, dataRole := service.userRepository.GetByNip(user.Nip)
	if httpStatusCheckRole == 500 || errCheckRole != nil {
		return 500, errCheckRole
	} else if dataRole.Nip != nil {
		return 500, errors.New("nip does exists on database")
	}

	saveData.Name = user.Name
	saveData.Nip = user.Nip
	saveData.Email = user.Email
	//saveData.RoleId = user.RoleId
	saveData.AuthServerId = user.AuthServerId

	httpStatusResponse, err := service.userRepository.Save(&saveData)
	if httpStatusResponse == 500 || err != nil {
		return 500, err
	}
	fmt.Printf("test3")
	return 200, nil
}

func (service *userService) GetById(id *string) (uint, error, *model.Users) {
	httpStatus, errorMessage, dataUser := service.userRepository.GetById(id)
	if httpStatus == 500 || errorMessage != nil {
		return 500, errorMessage, nil
	} else if dataUser.ID == nil {
		return 404, errors.New("user doesn't exists"), nil
	}

	return httpStatus, errorMessage, dataUser
}

func (service *userService) UpdateById(id *string, user *model.UserRequestUpdateDTO) (uint, error, *model.Users) {
	httpStatus, errorMessage, dataUser := service.userRepository.GetById(id)

	if httpStatus == 500 || errorMessage != nil {
		return 500, errorMessage, nil
	} else if dataUser.ID == nil {
		return 404, errors.New("user doesn't exists"), nil
	}

	dataUser.Name = user.Name
	dataUser.Email = user.Email

	httpStatusUpdate, errorMessageUpdate, dataUserUpdate := service.userRepository.UpdateById(dataUser)
	if httpStatusUpdate == 500 {
		return 500, errorMessageUpdate, nil
	}

	return 200, nil, dataUserUpdate
}

func (service *userService) GetAllByPaginate(structPagination *helperDatabases.QueryParamPaginationEntity) (uint, error, *apiresponse.Pagination) {
	httpStatus, errorMessage, responseObject := service.userRepository.GetAllByPaginate(structPagination)
	return httpStatus, errorMessage, responseObject
}

func (service *userService) DeleteById(id *string) (uint, error) {
	httpStatus, errorMessage, dataUser := service.userRepository.GetById(id)

	if httpStatus == 500 || errorMessage != nil {
		return 500, errorMessage
	} else if dataUser.ID == nil {
		return 404, errors.New("user doesn't exists")
	}

	httpStatusUpdate, errorMessageUpdate, _ := service.userRepository.UpdateById(dataUser)
	if httpStatusUpdate == 500 {
		return 500, errorMessageUpdate
	}

	return 200, nil
}
