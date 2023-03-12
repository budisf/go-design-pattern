package service

import (
	"errors"
	repository2 "ethical-be/modules/v1/utilities/role/repository"
	userRoleModel "ethical-be/modules/v1/utilities/user/model/user_role"
	"ethical-be/modules/v1/utilities/user/repository"
	repository3 "ethical-be/modules/v1/utilities/user/repository"
	"fmt"
	"strconv"
)

type IUserRoleService interface {
	UpdateById(id *string, user *userRoleModel.UserRoleRequestDTO) (uint, error, *userRoleModel.UserRolesRawQueryResult)
	GetById(id *string) (uint, error, *userRoleModel.UserRolesRawQueryResult)
}

type userRoleService struct {
	userRoleRepository repository.IUseRoleRepository
	userRepository     repository3.IUserRepository
	roleRepository     repository2.IRoleRepository
}

func InitUserRoleService(userRoleRepository repository.IUseRoleRepository, roleRepository repository2.IRoleRepository, userRepository repository3.IUserRepository) IUserRoleService {
	return &userRoleService{
		userRoleRepository: userRoleRepository,
		roleRepository:     roleRepository,
		userRepository:     userRepository,
	}
}

func (service *userRoleService) UpdateById(id *string, user *userRoleModel.UserRoleRequestDTO) (uint, error, *userRoleModel.UserRolesRawQueryResult) {
	roleIdUint64 := uint64(*user.RoleId)
	roleIdString := strconv.FormatUint(roleIdUint64, 10)

	httpStatusCheckRole, errCheckRole, dataRole := service.roleRepository.GetByIdRelation(&roleIdString)
	if httpStatusCheckRole == 500 || errCheckRole != nil {
		return 500, errCheckRole, nil
	} else if dataRole.ID == nil {
		return 404, errors.New("role_id doesn't exists"), nil
	}

	httpStatus, errorMessage, dataUser := service.userRepository.GetById(id)
	if httpStatus == 500 || errorMessage != nil {
		return 500, errorMessage, nil
	} else if dataUser.ID == nil {
		return 404, errors.New("user doesn't exists"), nil
	}

	dataUser.RoleId = user.RoleId

	httpStatusUpdate, errorMessageUpdate, _ := service.userRepository.UpdateById(dataUser)
	if httpStatusUpdate == 500 || errorMessageUpdate != nil {
		return 500, errorMessageUpdate, nil
	}

	httpStatusGetRaw, errorStatusGetRaw, dataRaw := service.userRoleRepository.GetByIdUsingRawQuery(id)
	if httpStatusGetRaw == 500 || errorStatusGetRaw != nil {
		return 500, errorStatusGetRaw, nil
	}

	return 200, nil, dataRaw
}

func (service *userRoleService) GetById(id *string) (uint, error, *userRoleModel.UserRolesRawQueryResult) {
	httpStatus, errorMessage, dataUser := service.userRepository.GetById(id)

	if httpStatus == 500 {
		return 500, errorMessage, nil
	} else if dataUser.ID == nil {
		return 404, errors.New("user doesn't exists"), nil
	}

	httpStatusGetRaw, errorStatusGetRaw, dataRaw := service.userRoleRepository.GetByIdUsingRawQuery(id)
	if httpStatusGetRaw == 500 {
		return 500, errorStatusGetRaw, nil
	} else if dataRaw.ID == nil {
		return 404, errors.New(fmt.Sprintf("User ID %v doesn't have role", *id)), nil
	}

	return 200, nil, dataRaw
}
