package repository

import (
	"ethical-be/app/config"
	roleModel "ethical-be/modules/v1/utilities/role/model"
	model2 "ethical-be/modules/v1/utilities/user/model"
	model "ethical-be/modules/v1/utilities/user/model/user_zone"
	helperDatabases "ethical-be/pkg/helpers/databases"
	"fmt"

	"gorm.io/gorm"
)

type IUserZoneRepository interface {
	AssignUserZone(userZone *model.UserZone) (uint, error, *model.UserZone)
	GetUserZoneByUserIdZoneId(userId *string, zoneId *string) (uint, error, *model.GetBySalesZoneIdUserIdRawQuery)
	GetUserZoneByUserId(userId *string) (uint, error, *model.GetBySalesZoneIdUserIdRawQuery)
	UpdateById(userZone *model.UserZone) (uint, error, *model.UserZone)
	UpdateUserAndZoneByUserIdZoneId(userId *string, zoneId *string, zoneType *model.SalesZoneType) (uint, error, *model2.Users)
	GetListUserByZoneId(zoneId *string, zoneType *model.SalesZoneType) (uint, error, *model.GetListUserByZoneRawQuery)
	GetBySalesZoneIDMultiple(id []uint, zoneType string, roleName string) ([]model.UserZone, error)
	GetByZoneType(zoneType string) ([]model.UserZone, error)
	GetBySalesZoneUserID(userId string) (model.UserZone, error)
	GetUserZoneByUserIDZoneIDZoneType(userId *string, zoneType *model.SalesZoneType, zoneId *string) (model.UserZone, error)
	GetAllUserVacant() ([]model.UserZone, error)
	GetUserVacantBySalesZoneIDMultiple(id []uint, zoneType string) ([]model.UserZone, error)
	GetAllUserNonVacant() ([]model.UserZone, error)
	GetUserNonVacantBySalesZoneIDMultiple(id []uint, zoneType string) ([]model.UserZone, error)
	GetSubordinateEmployeesByUserIDZoneIDZoneType(userId *string, zoneType *model.SalesZoneType, zoneId *string, roleName *string) (uint, error, *[]model.UserZoneWithEpochEntity)
	DeleteUserZoneByID(userZoneID string) (uint, error)
	GetZoneChildVacantBySalesZoneData(salesZoneType string, salesZoneId *int) ([]model.ZoneType, error)
	GetZoneChildRoleImpersonate(salesZoneType string, salesZoneId *int) ([]model.ZoneTypeRole, error)
}

type userZoneRepository struct {
	conf            *config.Conf
	conn            *gorm.DB
	helperDatabases helperDatabases.IHelperDatabases
}

func InitUserZoneRepository(conn *gorm.DB, helperDatabases helperDatabases.IHelperDatabases, conf *config.Conf) IUserZoneRepository {
	return &userZoneRepository{
		conf:            conf,
		conn:            conn,
		helperDatabases: helperDatabases,
	}
}

func (repo *userZoneRepository) AssignUserZone(userZone *model.UserZone) (uint, error, *model.UserZone) {
	if err := repo.conn.Create(&userZone).Error; err != nil {
		return 500, err, nil
	}
	repo.conn.Preload("user_zones").Find(&userZone)
	return 200, nil, userZone
}

func (repo *userZoneRepository) DeleteUserZoneByID(userZoneID string) (uint, error) {
	if err := repo.conn.Delete(model.UserZone{}, userZoneID).Error; err != nil {
		return 500, err
	}
	return 200, nil
}

func (repo *userZoneRepository) GetUserZoneByUserIdZoneId(userId *string, zoneId *string) (uint, error, *model.GetBySalesZoneIdUserIdRawQuery) {
	var userZone model.GetBySalesZoneIdUserIdRawQuery

	rawQuery := "select uz.id, uz.id, u.name as username, u.nip as user_nip, uz.sales_zone_id, " +
		"case when uz.sales_zone_type::sales_zone_type = 'group_territories' " +
		"then ( select gt.name from public.group_territories as gt where true and gt.id = uz.sales_zone_id ) " +
		"else " +
		"case when uz.sales_zone_type::sales_zone_type = 'areas' " +
		"then ( select a.name from public.areas as a where true and a.id = uz.sales_zone_id ) " +
		"else ( select r.name from public.regions as r where true and r.id = uz.sales_zone_id ) " +
		"end " +
		"end name_sales_zone, uz.sales_zone_type, round( cast(extract(epoch from uz.assigned_date) as numeric) ) as assigned_date, round( cast(extract(epoch from uz.finished_date) as numeric) ) as finished_date, round( cast(extract(epoch from uz.created_at) as numeric) ) as created_at, round( cast(extract(epoch from uz.updated_at) as numeric) ) as updated_at from public.user_zones as uz " +
		"inner join public.users as u on uz.user_id = u.id " +
		"where uz.user_id = @userId and uz.sales_zone_id = @zoneId and uz.finished_date is null and u.status = 'active';"

	err := repo.conn.Raw(rawQuery, map[string]interface{}{"userId": *userId, "zoneId": *zoneId}).Scan(&userZone).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &userZone
}

func (repo *userZoneRepository) GetUserZoneByUserId(userId *string) (uint, error, *model.GetBySalesZoneIdUserIdRawQuery) {
	var userZone model.GetBySalesZoneIdUserIdRawQuery

	rawQuery := "select uz.id, uz.id, uz.user_id, u.name as username, u.nip as user_nip, uz.sales_zone_id, " +
		"case when uz.sales_zone_type::sales_zone_type = 'group_territories' " +
		"then ( select gt.name from public.group_territories as gt where true and gt.id = uz.sales_zone_id ) " +
		"else case when uz.sales_zone_type::sales_zone_type = 'areas' " +
		"then ( select a.name from public.areas as a where true and a.id = uz.sales_zone_id ) " +
		"else case when uz.sales_zone_type::sales_zone_type = 'regions' " +
		"then ( select r.name from public.regions as r where true and r.id = uz.sales_zone_id ) " +
		"else ( select d.name from public.districts as d where true and d.id = uz.sales_zone_id )" +
		"end " +
		"end " +
		"end name_sales_zone, uz.sales_zone_type, round( cast(extract(epoch from uz.assigned_date) as numeric) ) as assigned_date, round( cast(extract(epoch from uz.finished_date) as numeric) ) as finished_date, round( cast(extract(epoch from uz.created_at) as numeric) ) as created_at, round( cast(extract(epoch from uz.updated_at) as numeric) ) as updated_at from public.user_zones as uz " +
		"inner join public.users as u on uz.user_id = u.id " +
		"where uz.user_id = @userId and uz.finished_date is null and u.status = 'active';"

	err := repo.conn.Raw(rawQuery, map[string]interface{}{"userId": *userId}).Scan(&userZone).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &userZone
}

func (repo *userZoneRepository) UpdateById(userZone *model.UserZone) (uint, error, *model.UserZone) {
	err := repo.conn.Save(&userZone).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, userZone
}

func (repo *userZoneRepository) UpdateUserAndZoneByUserIdZoneId(userId *string, zoneId *string, zoneType *model.SalesZoneType) (uint, error, *model2.Users) {
	var resultUserScan model2.Users
	var updateTableZone string

	if *zoneType == "group_territories" {
		updateTableZone = "update group_territories "
	} else if *zoneType == "areas" {
		updateTableZone = "update areas "
	} else if *zoneType == "areas" {
		updateTableZone = "update areas "
	} else if *zoneType == "regions" {
		updateTableZone = "update regions "
	} else if *zoneType == "districts" {
		updateTableZone = "update districts "
	}

	rawQueryUpdateMultiple := "with cte as ( " +
		"update public.user_zones set finished_date = current_timestamp, updated_at = current_timestamp " +
		"where true and user_zones.user_id = @userId and user_zones.sales_zone_id = @zoneId  and user_zones.sales_zone_type::sales_zone_type = @zoneType returning * " +
		") " +
		updateTableZone +
		"set is_vacant = true, updated_at = current_timestamp where true and id = ( select c.sales_zone_id from cte as c ) returning *"

	err := repo.conn.Raw(rawQueryUpdateMultiple, map[string]interface{}{"userId": *userId, "zoneId": *zoneId, "zoneType": *zoneType}).Scan(&resultUserScan).Error
	if err != nil {
		return 500, err, nil
	}
	return 200, nil, &resultUserScan
}

func (repo *userZoneRepository) GetListUserByZoneId(zoneId *string, zoneType *model.SalesZoneType) (uint, error, *model.GetListUserByZoneRawQuery) {
	var resultQuery model.GetListUserByZoneRawQuery

	rawQuery := "select uz.id, uz.sales_zone_id, uz.user_id, uz.sales_zone_type, us.name as username, us.nip as user_nip, us.role_id as user_role_id, " +
		"case when uz.sales_zone_type::sales_zone_type = 'group_territories' " +
		"then ( select gt.name from public.group_territories as gt where true and gt.id = uz.sales_zone_id ) " +
		"else " +
		"case when uz.sales_zone_type::sales_zone_type = 'areas' " +
		"then ( select a.name from public.areas as a where true and a.id = uz.sales_zone_id ) " +
		"else case when uz.sales_zone_type::sales_zone_type = 'regions'" +
		"then ( select r.name from public.regions as r where true and r.id = uz.sales_zone_id ) " +
		"else ( select d.name from public.districts as d where true and d.id = uz.sales_zone_id )" +
		"end " +
		"end " +
		"end name_sales_zone, round( cast(extract(epoch from uz.assigned_date) as numeric) ) as assigned_date, round( cast(extract(epoch from uz.finished_date) as numeric) ) as finished_date, round( cast(extract(epoch from uz.created_at) as numeric) ) as created_at, round( cast(extract(epoch from uz.updated_at) as numeric) ) as updated_at " +
		"from user_zones as uz inner join users as us on uz.id = us.id " +
		"where true and uz.sales_zone_id = @zoneId and " +
		"uz.sales_zone_type = @zoneType and uz.finished_date is null and us.status = 'active'"

	err := repo.conn.Raw(rawQuery, map[string]interface{}{"zoneId": *zoneId, "zoneType": *zoneType}).Scan(&resultQuery).Error
	if err != nil {
		return 500, err, nil
	}

	return 200, nil, &resultQuery
}

func (repo *userZoneRepository) GetBySalesZoneIDMultiple(id []uint, zoneType string, roleName string) ([]model.UserZone, error) {

	var (
		userZone []model.UserZone
	)

	err := repo.conn.Preload("Users", "deleted_at is null").Preload("Users.Role", "deleted_at is null").
		Table("user_zones").
		Joins("join users on users.id = user_zones.user_id ").
		Joins("join roles on users.role_id = roles.id").
		Where("roles.name = ?", roleName).
		Where("user_zones.sales_zone_id IN (?)", id).
		Where("user_zones.sales_zone_type = ?", zoneType).
		Where("user_zones.finished_date is NULL").
		Group("user_zones.user_id, user_zones.id").
		Find(&userZone).Error
	return userZone, err
}

func (repo *userZoneRepository) GetByZoneType(zoneType string) ([]model.UserZone, error) {
	var userZone []model.UserZone

	err := repo.conn.Preload("Users", "deleted_at is null").Preload("Users.Role", "deleted_at is null").
		Table("user_zones").
		Where("user_zones.sales_zone_type = ?", zoneType).
		Where("user_zones.finished_date is NULL").
		Find(&userZone).Error
	return userZone, err
}

func (repo *userZoneRepository) GetBySalesZoneUserID(userId string) (model.UserZone, error) {
	var userZone model.UserZone
	err := repo.conn.Where("finished_date is NULL").Where("user_id = ?", userId).Find(&userZone).Error
	return userZone, err
}
func (repo *userZoneRepository) GetUserZoneByUserIDZoneIDZoneType(userId *string, zoneType *model.SalesZoneType, zoneId *string) (model.UserZone, error) {
	var userZone model.UserZone
	err := repo.conn.Preload("Users", "deleted_at is null").Preload("Users.Role", "deleted_at is null").
		Table("user_zones").
		Where("user_zones.sales_zone_type = ?", zoneType).
		Where("user_zones.user_id = ?", userId).
		Where("user_zones.sales_zone_id = ?", zoneId).
		Where("user_zones.finished_date is NULL").
		Find(&userZone).Error
	return userZone, err
}

func (repo *userZoneRepository) GetAllUserVacant() ([]model.UserZone, error) {
	var userZone []model.UserZone
	err := repo.conn.Preload("Users", "deleted_at is null").Preload("Users.Role", "deleted_at is null").
		Table("user_zones").
		Where("user_zones.finished_date is NOT NULL").
		Find(&userZone).Error
	return userZone, err
}

func (repo *userZoneRepository) GetUserVacantBySalesZoneIDMultiple(id []uint, zoneType string) ([]model.UserZone, error) {

	var (
		userZone []model.UserZone
	)

	err := repo.conn.Preload("Users").Preload("Users.Role", "deleted_at is null").
		Table("user_zones").
		Joins("join users on users.id = user_zones.user_id ").
		Joins("join roles on users.role_id = roles.id").
		Where("user_zones.sales_zone_id IN (?)", id).
		Where("user_zones.sales_zone_type = ?", zoneType).
		Where("user_zones.finished_date is not NULL").
		Group("user_zones.user_id, user_zones.id").
		Find(&userZone).Error
	return userZone, err
}

func (repo *userZoneRepository) GetSubordinateEmployeesByUserIDZoneIDZoneType(userId *string, zoneType *model.SalesZoneType, zoneId *string, roleName *string) (uint, error, *[]model.UserZoneWithEpochEntity) {
	var userZone []model.UserZoneWithEpochEntity
	var CTEQuery string

	fmt.Println("userId: ", *userId)

	// this repository has been used on module sales product and outlet, which is user_id must order as asc. what if change on order by, please to make sure again
	// on module sales product and outlet handler: dryCheckGetSalesProductValidation and dryCheckGetSalesOutletValidation (so impactful).
	// what if there must be a change in the order by. please add params in this function for optional order by for the desired field.
	if *zoneType == model.District {
		CTEQuery = "with dataUserRegionsUnderDistrict as ( " +
			"select uz2.* from public.user_zones as uz2 inner join public.users as u2 on uz2.user_id = u2.id where true and uz2.finished_date is null and uz2.sales_zone_id in ( " +
			"select r2.id from public.districts as d2 inner join public.regions as r2 on d2.id = r2.district_id where true and r2.district_id in ( " +
			"select d.id from public.districts as d " +
			"inner join public.regions as r on r.district_id = d.id " +
			"inner join public.user_zones as uz on uz.sales_zone_id = d.id " +
			"inner join public.users as u on uz.user_id = u.id " +
			"inner join public.roles as ro on u.role_id = ro.id " +
			"where true and uz.user_id = @userId and uz.sales_zone_type = @zoneType and uz.sales_zone_id = @zoneId and ro.name = @roleName and uz.finished_date is null and u.status = 'active' group by d.id " +
			") " +
			") and uz2.finished_date is null and u2.status = 'active' " +
			"), dataUserAreaUnderRegion as ( " +
			"select uz3.* from public.user_zones as uz3 inner join public.users as u3 on uz3.user_id = u3.id where true and uz3.finished_date is null and uz3.sales_zone_id in ( " +
			"select a.id from public.areas as a inner join public.regions as r3 on a.region_id = r3.id where true and r3.id in ( select dug.sales_zone_id from dataUserRegionsUnderDistrict as dug )" +
			") and uz3.finished_date is null and u3.status = 'active' " +
			"), dataUserGTUnderArea as ( " +
			"select uz4.* from public.user_zones as uz4 inner join public.users as u4 on uz4.user_id = u4.id where true and uz4.finished_date is null and uz4.sales_zone_id in ( " +
			"select gt.id from public.areas as a2 inner join public.group_territories as gt on gt.area_id = a2.id where true and a2.id in ( select dua.sales_zone_id from dataUserAreaUnderRegion as dua ) " +
			") and uz4.finished_date is null and u4.status = 'active' " +
			"), dataEmployee as ( " +
			"select * from dataUserRegionsUnderDistrict union all select * from dataUserAreaUnderRegion union all select * from dataUserGTUnderArea" +
			") select d.id as user_zone_id, d.user_id, d.sales_zone_id, d.sales_zone_type, round( cast(extract(epoch from d.assigned_date) as numeric) ) as assigned_date, round( cast(extract(epoch from d.created_at) as numeric) ) as created_at, round( cast(extract(epoch from d.updated_at) as numeric) ) as updated_at " +
			"from dataEmployee as d order by d.user_id asc "
	} else if *zoneType == model.Region {
		CTEQuery = "with dataUserAreaUnderRegion as ( select uz.* from user_zones as uz where true and uz.finished_date is null " +
			"and uz.sales_zone_id in ( select a.id from areas as a " +
			"inner join public.regions as r on a.region_id = r.id " +
			"inner join public.user_zones as uz on uz.sales_zone_id = r.id " +
			"inner join public.users as u on uz.user_id = u.id " +
			"inner join public.roles as ro on u.role_id = ro.id " +
			"where true and uz.user_id = @userId and uz.sales_zone_type = @zoneType and uz.sales_zone_id = @zoneId and ro.name = @roleName and uz.finished_date is null and u.status = 'active' ) " +
			"), " +
			"dataUsersGTUnderArea as ( select uz2.* from public.user_zones as uz2 inner join public.users as u2 on uz2.user_id = u2.id where true and u2.status = 'active' and uz2.sales_zone_type = 'group_territories' and uz2.finished_date is null and uz2.sales_zone_id in ( select gt.id from public.group_territories as gt inner join public.areas as a2 on gt.area_id = a2.id where true and a2.id in ( select dua.sales_zone_id from dataUserAreaUnderRegion as dua ) ) " +
			"), " +
			"dataUserAreaAndGT as ( select dua.* from dataUserAreaUnderRegion as dua union all select dug.* from dataUsersGTUnderArea as dug " +
			") select duaa.id as user_zone_id, duaa.user_id, duaa.sales_zone_id, duaa.sales_zone_type, round( cast(extract(epoch from duaa.assigned_date) as numeric) ) as assigned_date, round( cast(extract(epoch from duaa.created_at) as numeric) ) as created_at, round( cast(extract(epoch from duaa.updated_at) as numeric) ) as updated_at " +
			"from dataUserAreaAndGT as duaa order by duaa.user_id asc "
	} else if *zoneType == model.Area {
		CTEQuery = "with dataUserGTUnderArea as ( select uz.* from public.user_zones as uz where true and uz.finished_date is null and uz.sales_zone_type = 'group_territories' and " +
			"uz.sales_zone_id in ( select gt.id from public.areas as a " +
			"inner join public.group_territories as gt on gt.area_id = a.id " +
			"inner join public.user_zones as uz on uz.sales_zone_id = a.id " +
			"inner join public.users as u on uz.user_id = u.id " +
			"inner join public.roles as ro on u.role_id = ro.id " +
			"where true and uz.user_id = @userId and uz.sales_zone_type = @zoneType and uz.sales_zone_id = @zoneId and ro.name = @roleName and uz.finished_date is null and u.status = 'active'  ) " +
			") select dug.id as user_zone_id, dug.user_id, dug.sales_zone_id, dug.sales_zone_type, " +
			"round( cast(extract(epoch from dug.assigned_date) as numeric) ) as assigned_date, round( cast(extract(epoch from dug.created_at) as numeric) ) as created_at, " +
			"round( cast(extract(epoch from dug.updated_at) as numeric) ) as updated_at " +
			"from dataUserGTUnderArea as dug order by dug.user_id asc"
	}

	fmt.Println("DEBUG QUERY: ", CTEQuery)

	err := repo.conn.Raw(CTEQuery, map[string]interface{}{"userId": *userId, "zoneId": *zoneId, "zoneType": *zoneType, "roleName": *roleName}).Scan(&userZone).Error
	if err != nil {
		return 500, err, nil
	}

	return 200, nil, &userZone
}

func (repo *userZoneRepository) GetAllUserNonVacant() ([]model.UserZone, error) {
	var userZone []model.UserZone
	err := repo.conn.Preload("Users", "deleted_at is null").Preload("Users.Role", "deleted_at is null").
		Table("user_zones").
		Where("user_zones.finished_date is NULL").
		Find(&userZone).Error
	return userZone, err
}

func (repo *userZoneRepository) GetUserNonVacantBySalesZoneIDMultiple(id []uint, zoneType string) ([]model.UserZone, error) {

	var (
		userZone []model.UserZone
	)

	err := repo.conn.Preload("Users", "deleted_at is null").Preload("Users.Role", "deleted_at is null").
		Table("user_zones").
		Joins("join users on users.id = user_zones.user_id ").
		Joins("join roles on users.role_id = roles.id").
		Where("user_zones.sales_zone_id IN (?)", id).
		Where("user_zones.sales_zone_type = ?", zoneType).
		Where("user_zones.finished_date is NULL").
		Group("user_zones.user_id, user_zones.id").
		Find(&userZone).Error
	return userZone, err
}

// Create a scope for statement used in GetZoneChildVacantBySalesZoneData
func WhereIdSalesZoneIds(salesZoneType model.SalesZoneType, salesZoneIdSubQuery ...interface{}) func(db *gorm.DB) *gorm.DB {
	var parentZoneTypeField string
	switch salesZoneType {
	case model.District:
		parentZoneTypeField = "id"
	case model.Region:
		parentZoneTypeField = "district_id"
	case model.Area:
		parentZoneTypeField = "region_id"
	case model.GroupTerritory:
		parentZoneTypeField = "area_id"
	default:
		parentZoneTypeField = "id"
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(parentZoneTypeField+" IN ?", salesZoneIdSubQuery).Table(string(salesZoneType))
	}
}

func SelectBySalesZoneType(salesZoneType model.SalesZoneType, role roleModel.Roles, roleOnly bool) func(db *gorm.DB) *gorm.DB {
	var Query string
	return func(db *gorm.DB) *gorm.DB {
		if roleOnly == true {
			Query = fmt.Sprintf("%+v", *role.ID) +
				" as role_id, '" + *role.Name +
				"' as role_name, '" + *role.Label + "' as role_label"
		} else {
			Query = "'" + string(salesZoneType) + "' as zone_type, id as zone_id, name, is_vacant, " +
				fmt.Sprintf("%+v", *role.ID) +
				" as role_id, '" + *role.Name +
				"' as role_name, '" + *role.Label + "' as role_label"
		}

		return db.Select(Query)
	}
}

func (r *userZoneRepository) GetZoneChildVacantBySalesZoneData(salesZoneType string, salesZoneId *int) ([]model.ZoneType, error) {

	// Get Roles
	var roles []roleModel.Roles
	r.conn.Find(&roles)

	var NSM roleModel.Roles
	var SM roleModel.Roles
	var ASM roleModel.Roles
	var FF roleModel.Roles

	for _, role := range roles {
		if roleModel.RoleType(*role.Name) == roleModel.NSM {
			NSM = role
		}
		if roleModel.RoleType(*role.Name) == roleModel.SM {
			SM = role
		}
		if roleModel.RoleType(*role.Name) == roleModel.ASM {
			ASM = role
		}
		if roleModel.RoleType(*role.Name) == roleModel.FieldForce {
			FF = role
		}
	}

	var subQuery *gorm.DB
	var txDistrict *gorm.DB

	if salesZoneType == string(model.District) {
		// Create new dry sql to get district by zone id
		txDistrict = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.District,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.District,
				NSM,
				false,
			)).
			Where("is_vacant = ?", true)

		// Create new dry sql for subquery to get all
		// current child 1 level below self
		subQuery = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.District,
				*salesZoneId,
			)).Select("id")

	} else if salesZoneType != string(model.District) &&
		salesZoneType != string(model.Region) &&
		salesZoneType != string(model.Area) &&
		salesZoneType != string(model.GroupTerritory) {
		// if sales zone is not in (district region area gt)

		// Create new dry sql to get district by zone id
		txDistrict = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("districts")
			},
			SelectBySalesZoneType(
				model.District,
				NSM,
				false,
			)).
			Where("is_vacant = ?", true)

		// Create new dry sql for subquery to get all
		// current child 1 level below self
		subQuery = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("districts")
			}).Select("id")
	}

	var txRegion *gorm.DB

	if salesZoneType == string(model.Region) {
		txRegion = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Region,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.Region,
				SM,
				false,
			)).
			Where("is_vacant = ?", true)

		subQuery = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("regions")
			}).
			Where("id = ?", salesZoneId).
			Select("id")
	} else {
		txRegion = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Region,
				subQuery,
			),
			SelectBySalesZoneType(
				model.Region,
				SM,
				false,
			)).
			Where("is_vacant = ?", true)

		subQuery = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Region,
				subQuery,
			)).
			Select("id")
	}

	var txArea *gorm.DB

	if salesZoneType == string(model.Area) {
		txArea = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Area,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.Area,
				ASM,
				false,
			)).
			Where("is_vacant = ?", true)

		subQuery = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("areas")
			}).
			Where("id = ?", salesZoneId).
			Select("id")
	} else {
		txArea = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Area,
				subQuery,
			),
			SelectBySalesZoneType(
				model.Area,
				ASM,
				false,
			)).
			Where("is_vacant = ?", true)

		subQuery = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Area,
				subQuery,
			)).
			Select("id")
	}

	var txGroupTerritory *gorm.DB

	if salesZoneType == string(model.GroupTerritory) {
		txGroupTerritory = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.GroupTerritory,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.GroupTerritory,
				FF,
				false,
			)).
			Where("is_vacant = ?", true)
	} else {
		txGroupTerritory = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.GroupTerritory,
				subQuery,
			),
			SelectBySalesZoneType(
				model.GroupTerritory,
				FF,
				false,
			)).
			Where("is_vacant = ?", true)
	}

	var zoneType []model.ZoneType
	var txRaw *gorm.DB
	var txErr error
	switch salesZoneType {
	case string(model.District):
		txRaw = r.conn.Raw("? UNION ? UNION ? order by zone_type desc, zone_id asc",
			txRegion,
			txArea,
			txGroupTerritory,
		).Scan(&zoneType)
		txErr = txRaw.Error
	case string(model.Region):
		txRaw = r.conn.Raw("? UNION ? order by zone_type desc, zone_id asc",
			txArea,
			txGroupTerritory,
		).Scan(&zoneType)
		txErr = txRaw.Error
	case string(model.Area):
		txRaw = r.conn.Raw("? order by zone_type desc, zone_id asc",
			txGroupTerritory,
		).Scan(&zoneType)
		txErr = txRaw.Error
	case string(model.GroupTerritory):
		zoneType = []model.ZoneType{}
		txErr = nil
	default:
		txRaw = r.conn.Raw("? UNION ? UNION ? UNION ? order by zone_type desc, zone_id asc",
			txDistrict,
			txRegion,
			txArea,
			txGroupTerritory,
		).Scan(&zoneType)
		txErr = txRaw.Error
	}
	if txErr != nil {
		return nil, txErr
	} else {
		return zoneType, nil
	}
}

func (r *userZoneRepository) GetZoneChildRoleImpersonate(salesZoneType string, salesZoneId *int) ([]model.ZoneTypeRole, error) {

	// Get Roles
	var roles []roleModel.Roles
	r.conn.Find(&roles)

	var NSM roleModel.Roles
	var SM roleModel.Roles
	var ASM roleModel.Roles
	var FF roleModel.Roles

	for _, role := range roles {
		if roleModel.RoleType(*role.Name) == roleModel.NSM {
			NSM = role
		}
		if roleModel.RoleType(*role.Name) == roleModel.SM {
			SM = role
		}
		if roleModel.RoleType(*role.Name) == roleModel.ASM {
			ASM = role
		}
		if roleModel.RoleType(*role.Name) == roleModel.FieldForce {
			FF = role
		}
	}

	var subQuery *gorm.DB
	var txDistrict *gorm.DB

	if salesZoneType == string(model.District) {
		// Create new dry sql to get district by zone id
		txDistrict = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.District,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.District,
				NSM,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")

		// Create new dry sql for subquery to get all
		// current child 1 level below self
		subQuery = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.District,
				*salesZoneId,
			)).Select("id")

	} else if salesZoneType != string(model.District) &&
		salesZoneType != string(model.Region) &&
		salesZoneType != string(model.Area) &&
		salesZoneType != string(model.GroupTerritory) {
		// if sales zone is not in (district region area gt)

		// Create new dry sql to get district by zone id
		txDistrict = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("districts")
			},
			SelectBySalesZoneType(
				model.District,
				NSM,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")

		// Create new dry sql for subquery to get all
		// current child 1 level below self
		subQuery = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("districts")
			}).Select("id")
	}

	var txRegion *gorm.DB

	if salesZoneType == string(model.Region) {
		txRegion = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Region,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.Region,
				SM,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")

		subQuery = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("regions")
			}).
			Where("id = ?", salesZoneId).
			Select("id")
	} else {
		txRegion = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Region,
				subQuery,
			),
			SelectBySalesZoneType(
				model.Region,
				SM,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")

		subQuery = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Region,
				subQuery,
			)).
			Select("id")
	}

	var txArea *gorm.DB

	if salesZoneType == string(model.Area) {
		txArea = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Area,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.Area,
				ASM,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")

		subQuery = r.conn.Scopes(
			func(db *gorm.DB) *gorm.DB {
				return db.Table("areas")
			}).
			Where("id = ?", salesZoneId).
			Select("id")
	} else {
		txArea = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Area,
				subQuery,
			),
			SelectBySalesZoneType(
				model.Area,
				ASM,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")

		subQuery = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.Area,
				subQuery,
			)).
			Select("id")
	}

	var txGroupTerritory *gorm.DB

	if salesZoneType == string(model.GroupTerritory) {
		txGroupTerritory = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.GroupTerritory,
				*salesZoneId,
			),
			SelectBySalesZoneType(
				model.GroupTerritory,
				FF,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")

	} else {
		txGroupTerritory = r.conn.Scopes(
			WhereIdSalesZoneIds(
				model.GroupTerritory,
				subQuery,
			),
			SelectBySalesZoneType(
				model.GroupTerritory,
				FF,
				true,
			)).
			Where("is_vacant = ?", true).
			Group("role_id, role_name, role_label")
	}

	var zoneTypeRole []model.ZoneTypeRole
	var txRaw *gorm.DB
	var txErr error
	switch salesZoneType {
	case string(model.District):
		txRaw = r.conn.Raw("? UNION ? UNION ? order by role_id asc",
			txRegion,
			txArea,
			txGroupTerritory,
		).Scan(&zoneTypeRole)
		txErr = txRaw.Error
	case string(model.Region):
		txRaw = r.conn.Raw("? UNION ? order by role_id asc",
			txArea,
			txGroupTerritory,
		).Scan(&zoneTypeRole)
		txErr = txRaw.Error
	case string(model.Area):
		txRaw = r.conn.Raw("? order by role_id asc",
			txGroupTerritory,
		).Scan(&zoneTypeRole)
		txErr = txRaw.Error
	case string(model.GroupTerritory):
		zoneTypeRole = []model.ZoneTypeRole{}
		txErr = nil
	default:
		txRaw = r.conn.Raw("? UNION ? UNION ? UNION ? order by role_id asc",
			txDistrict,
			txRegion,
			txArea,
			txGroupTerritory,
		).Scan(&zoneTypeRole)
		txErr = txRaw.Error
	}

	if txErr != nil {
		return nil, txErr
	} else {
		return zoneTypeRole, nil
	}
}
