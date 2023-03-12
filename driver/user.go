package driver

import (
	"ethical-be/modules/v1/utilities/user/handler"
	handler2 "ethical-be/modules/v1/utilities/user/handler"
	"ethical-be/modules/v1/utilities/user/repository"
	"ethical-be/modules/v1/utilities/user/service"
	service2 "ethical-be/modules/v1/utilities/user/service"
	service3 "ethical-be/modules/v1/utilities/user/service"
	helperDatabases "ethical-be/pkg/helpers/databases"
)

var (
	HelperDatabase     = helperDatabases.InitHelperDatabase(DB)
	UserZoneRepository = repository.InitUserZoneRepository(DB, HelperDatabase, &Conf)
	UserRepository     = repository.InitUserRepository(DB, HelperDatabase, &Conf)
	UserRoleRepository = repository.InitUserRoleRepository(DB, HelperDatabase, &Conf)
	UserZoneService    = service3.InitUserZoneService(UserZoneRepository, UserRoleRepository, UserRepository, RegionsRepository, AreaRepository, GtRepository, DistrictRepository)
	UserService        = service2.InitUserRepository(UserRepository, RoleRepository)
	UserRoleService    = service.InitUserRoleService(UserRoleRepository, RoleRepository, UserRepository)
	UserHandler        = handler.InitUserHandler(UserService)
	UserRoleHandler    = handler2.InitUserRoleHandler(UserRoleService)
	UserZoneHandler    = handler2.InitUserZoneHandler(UserZoneService, UserRoleService)
)
