package repository

import (
	"ethical-be/app/config"
	modelRole "ethical-be/modules/v1/utilities/role/model"
	"ethical-be/modules/v1/utilities/user/model/user_role"
	helperDatabases "ethical-be/pkg/helpers/databases"

	"gorm.io/gorm"
)

type IUseRoleRepository interface {
	GetByIdUsingRawQuery(id *string) (uint, error, *model.UserRolesRawQueryResult)
	GetByName(name string) (*modelRole.Roles, error)
}

type userRoleRepository struct {
	conf            *config.Conf
	conn            *gorm.DB
	helperDatabases helperDatabases.IHelperDatabases
}

func InitUserRoleRepository(conn *gorm.DB, helperDatabases helperDatabases.IHelperDatabases, conf *config.Conf) IUseRoleRepository {
	return &userRoleRepository{
		conf:            conf,
		conn:            conn,
		helperDatabases: helperDatabases,
	}
}

func (repo *userRoleRepository) GetByIdUsingRawQuery(id *string) (uint, error, *model.UserRolesRawQueryResult) {
	var userRole model.UserRolesRawQueryResult
	rawQuery := "select  u.id, u.name, u.auth_server_id, u.nip, u.role_id, round( cast(extract(epoch from u.created_at) as numeric) ) as created_at, round( cast(extract(epoch from u.updated_at) as numeric) ) as updated_at, " +
		" r.name as role_name, r.label as role_label, r.parent_id, r2.name as name_role_head " +
		"from users as u " +
		"inner join roles as r on u.role_id = r.id " +
		"left join public.roles as r2 on r.parent_id = r2.id " +
		"where u.id = @id and u.status = 'active'"

	err := repo.conn.Raw(rawQuery, map[string]interface{}{"id": *id}).Scan(&userRole).Error
	if err != nil {
		return 500, err, nil
	}

	return 200, nil, &userRole
}

func (repo *userRoleRepository) GetByName(name string) (*modelRole.Roles, error) {
	var userRole modelRole.Roles
	err := repo.conn.Where("name = ?", name).Find(&userRole).Error
	return &userRole, err
}
