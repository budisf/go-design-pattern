package driver

import (
	"ethical-be/modules/v1/utilities/role/handler"
	"ethical-be/modules/v1/utilities/role/repository"
	"ethical-be/modules/v1/utilities/role/service"
)

var (
	RoleRepository = repository.InitRoleRepository(DB, HelperDatabase, &Conf)
	RoleService    = service.InitRoleService(RoleRepository, UserRepository)
	RoleHandler    = handler.InitRoleHandler(RoleService)
)
