package service

import (
	"errors"
	model "ethical-be/modules/v1/utilities/user/model/user_zone"
	"ethical-be/modules/v1/utilities/user/repository"
	areaModel "ethical-be/modules/v1/utilities/zone/area/models"
	areaRelation "ethical-be/modules/v1/utilities/zone/area/models/relations"
	modelDistrict "ethical-be/modules/v1/utilities/zone/district/models"
	districtRelation "ethical-be/modules/v1/utilities/zone/district/models/relation"
	repository4 "ethical-be/modules/v1/utilities/zone/district/repository"
	gt "ethical-be/modules/v1/utilities/zone/gt/models"
	groupTerritoriesRelation "ethical-be/modules/v1/utilities/zone/gt/models/relations"
	gtRepository "ethical-be/modules/v1/utilities/zone/gt/repository"
	modelRegion "ethical-be/modules/v1/utilities/zone/region/models"
	regionRelation "ethical-be/modules/v1/utilities/zone/region/models/relations"
	"net/http"
	"time"

	repository3 "ethical-be/modules/v1/utilities/zone/area/repository"
	repository2 "ethical-be/modules/v1/utilities/zone/region/repository"

	"fmt"
	"strconv"
)

type IUserZoneService interface {
	AssignUserZone(userId *string, userZoneDTO *model.UserZoneRequestDTO) (uint, error)
	GetUserZoneByUserId(userId *string) (uint, error, *model.GetBySalesZoneIdUserIdRawQuery)
	UpdateFinishedAssigment(userId *string, userZoneDTO *model.UserZoneRequestDTO) (uint, error)
	GetListUserByZoneId(zoneId *string, zoneType *string) (uint, error, *model.GetListUserByZoneRawQuery)
	GetChildByUserId(userId string, roleName string) (*[]model.UserResponse, error, int)
	GetUserZoneByUserIDZoneIDZoneType(userId *string, zoneType *string, zoneId *string) (uint, error, *model.UserZone)
	GetChildVacantByUserId(userId string) (*[]model.UserResponse, error, int)
	GetSubordinateEmployeesByUserIDZoneIDZoneType(userId *string, zoneType *string, zoneId *string, roleName *string) (uint, error, *[]model.UserZoneWithEpochEntity)
	GetChildNonVacantByUserId(userId string) (*[]model.UserResponse, error, int)
	GetZoneChildVacantByUserId(userId string) ([]model.ZoneType, error, int)
	dryValidationZone(salesZoneID int, salesZoneType string) (uint, error)
	GetZoneByUserID(salesZoneType *string, userID *string) (uint, error, []model.GetZoneByUserIDResponse)
	ImpersonateAccessControlSales(userId string) ([]model.ZoneTypeRole, error, int)
}

var (
	ZoneDistric  = "districts"
	ZoneArea     = "areas"
	ZoneRegion   = "regions"
	ZoneGT       = "group_territories"
	FieldForce   = "field-force"
	Asm          = "asm"
	Sm           = "sm"
	FieldForceId = 7
	AsmId        = 6
	SmId         = 5
	NsmId        = 4
	RoleIds      = []int{FieldForceId, AsmId, SmId, NsmId} //(ff,asm,sm,nsm)
)

type userZoneService struct {
	userZoneRepository repository.IUserZoneRepository
	userRoleRepository repository.IUseRoleRepository
	userRepository     repository.IUserRepository
	districtRepository repository4.RepositoryDistrict
	regionRepository   repository2.RegionsRepository
	areaRepository     repository3.AreasReporitory
	gtRepository       gtRepository.GtRepository
}

func InitUserZoneService(userZoneRepository repository.IUserZoneRepository, userRoleRepository repository.IUseRoleRepository,
	userRepository repository.IUserRepository, regionRepository repository2.RegionsRepository, areaRepository repository3.AreasReporitory,
	gtRepository gtRepository.GtRepository, districtRepository repository4.RepositoryDistrict) IUserZoneService {
	return &userZoneService{
		userZoneRepository: userZoneRepository,
		userRoleRepository: userRoleRepository,
		userRepository:     userRepository,
		regionRepository:   regionRepository,
		areaRepository:     areaRepository,
		gtRepository:       gtRepository,
		districtRepository: districtRepository,
	}
}

func (service *userZoneService) dryValidationZone(salesZoneID int, salesZoneType string) (uint, error) {

	if salesZoneType == ZoneGT {
		dataGT, errorMessageGT := service.gtRepository.FindById(salesZoneID)
		if errorMessageGT != nil {
			return 500, errorMessageGT
		} else if dataGT.ID == 0 {
			return 404, errors.New(fmt.Sprintf("the data GT with ID %v isn't exists", salesZoneID))
		}
	} else if salesZoneType == ZoneArea {
		fmt.Println("area 2")
		dataArea, errorMessageArea := service.areaRepository.FindById(salesZoneID)
		if errorMessageArea != nil {
			return 500, errorMessageArea
		} else if dataArea.ID == 0 {
			return 404, errors.New(fmt.Sprintf("the data Area with ID %v isn't exists", salesZoneID))
		}
	} else if salesZoneType == ZoneRegion {
		dataRegion, errorMessageRegion := service.regionRepository.FindById(salesZoneID)
		if errorMessageRegion != nil {
			return 500, errorMessageRegion
		} else if dataRegion.ID == 0 {
			return 404, errors.New(fmt.Sprintf("the data Area with ID %v isn't exists", salesZoneID))
		}
	} else if salesZoneType == ZoneDistric {
		fmt.Println("Waiting struct districts")
		return 200, nil
	}

	return 200, nil
}

func (service *userZoneService) AssignUserZone(userId *string, userZoneDTO *model.UserZoneRequestDTO) (uint, error) {
	var enumSalesZone model.SalesZoneType
	var userZoneEntity model.UserZone
	var dataZoneInterface interface{}
	httpStatus, errorMessage, dataUser := service.userRepository.GetById(userId)

	// check user was exist
	if httpStatus == 500 || errorMessage != nil {
		return 500, errorMessage
	} else if dataUser.ID == nil {
		return 404, errors.New("user doesn't exists")
	}

	// validation zone from dry function
	httpStatusValidation, errorValidation := service.dryValidationZone(int(*userZoneDTO.SalesZoneId), *userZoneDTO.SalesZoneType)
	if httpStatusValidation == 500 && errorValidation != nil {
		return 500, errorValidation
	} else if httpStatusValidation == 404 && errorValidation != nil {
		return 404, errorValidation
	}

	if *userZoneDTO.SalesZoneType == ZoneGT {
		enumSalesZone = model.GroupTerritory
		dataGTVacant, errorDataVacant := service.gtRepository.FindById_ByIsVacant(int(*userZoneDTO.SalesZoneId), true)
		if errorDataVacant != nil {
			return 500, errorDataVacant
		}
		if dataGTVacant.ID == 0 {
			return 500, errors.New(fmt.Sprintf("Status Zone isn't vacant. User can't assign on this zone: sales_zone_id %v and sales_zone_type %v (GT)", int(*userZoneDTO.SalesZoneId), *userZoneDTO.SalesZoneType))
		}
		dataZoneInterface = dataGTVacant
	} else if *userZoneDTO.SalesZoneType == ZoneArea {
		enumSalesZone = model.Area
		dataAreaVacant, errorDataAreaVacant := service.areaRepository.FindById_ByIsVacant(int(*userZoneDTO.SalesZoneId), true)
		if errorDataAreaVacant != nil {
			return 500, errorDataAreaVacant
		}
		if dataAreaVacant.ID == 0 {
			return 500, errors.New(fmt.Sprintf("Status Zone isn't vacant. User can't assign on this zone: sales_zone_id %v and sales_zone_type %v (Area)", int(*userZoneDTO.SalesZoneId), *userZoneDTO.SalesZoneType))
		}
		dataZoneInterface = dataAreaVacant
	} else if *userZoneDTO.SalesZoneType == ZoneRegion {
		enumSalesZone = model.Region
		dataRegionVacant, errorDataVacantRegion := service.regionRepository.FindById_ByIsVacant(int(*userZoneDTO.SalesZoneId), true)
		if errorDataVacantRegion != nil {
			return 500, errorDataVacantRegion
		}
		if dataRegionVacant.ID == 0 {
			return 500, errors.New(fmt.Sprintf("Status Zone isn't vacant. User can't assign on this zone: sales_zone_id %v and sales_zone_type %v (Region)", int(*userZoneDTO.SalesZoneId), *userZoneDTO.SalesZoneType))
		}
		dataZoneInterface = dataRegionVacant
	} else {
		enumSalesZone = model.District
		dataDistrictVacant, errorDataVacantDistrict := service.districtRepository.FindById_ByIsVacant(int(*userZoneDTO.SalesZoneId), true)
		if errorDataVacantDistrict != nil {
			return 500, errorDataVacantDistrict
		}
		if dataDistrictVacant.ID == 0 {
			return 500, errors.New(fmt.Sprintf("Status Zone isn't vacant. User can't assign on this zone: sales_zone_id %v and sales_zone_type %v (District)", int(*userZoneDTO.SalesZoneId), *userZoneDTO.SalesZoneType))
		}
		dataZoneInterface = dataDistrictVacant
	}

	// check user id when assign sales zone id and sales zone type have zone id assigned or not on table user_zone
	// user FF, ASM, SM, NSM only have one sales zone id and sales zone type (current need business).
	salesZoneIDString := strconv.FormatUint(uint64(*userZoneDTO.SalesZoneId), 10)
	httpStatusGetUserZoneByUserID, errorMessageUserZoneByUserID, dataUserZoneByUserID := service.userZoneRepository.GetListUserByZoneId(&salesZoneIDString, &enumSalesZone)
	if httpStatusGetUserZoneByUserID == 500 && errorMessageUserZoneByUserID != nil {
		return 500, errorMessageUserZoneByUserID
	} else if dataUserZoneByUserID.ID != nil {
		return 500, errors.New(fmt.Sprintf("User ID %v couldn't been registered on sales_zone_id: %v, sales_zone_type: %v with Name Sales Zone: %v. User/ Zone only have one responsible person", *userId, *dataUserZoneByUserID.SalesZoneId, dataUserZoneByUserID.SalesZoneType, *dataUserZoneByUserID.NameSalesZone))
	}

	userIdUint64, _ := strconv.ParseUint(*userId, 10, 32)
	userIdUint := uint(userIdUint64)

	assignedDate := time.Now()

	userZoneEntity.UserId = &userIdUint
	userZoneEntity.SalesZoneId = userZoneDTO.SalesZoneId
	userZoneEntity.SalesZoneType = enumSalesZone
	userZoneEntity.AssignedDate = &assignedDate
	userZoneEntity.CreatedAt = &assignedDate

	httpStatusAssign, errorAssign, dataInsertUserZone := service.userZoneRepository.AssignUserZone(&userZoneEntity)
	if httpStatusAssign == 500 {
		return 500, errorAssign
	}
	userIDZoneString := strconv.FormatUint(uint64(*dataInsertUserZone.ID), 10)
	stringSalesZone := fmt.Sprintf("%v", enumSalesZone)
	// to update on table zone
	if stringSalesZone == ZoneDistric {
		var districtData modelDistrict.Districts
		d := dataZoneInterface.(districtRelation.DistritcRelation)
		districtData.ID = d.ID
		districtData.Name = d.Name
		districtData.CreatedAt = d.CreatedAt
		districtData.UpdatedAt = time.Now()
		districtData.DeletedAt = d.DeletedAt
		districtData.IsVacant = false

		_, errorUpdateDistrict := service.districtRepository.Update(districtData)
		if errorUpdateDistrict != nil {
			// given rollback to delete the data from table user zone when event error (500) update on table district
			httpDeleteUserZone, messageErrorDelete := service.userZoneRepository.DeleteUserZoneByID(userIDZoneString)
			if httpDeleteUserZone == 500 && messageErrorDelete != nil {
				return 500, messageErrorDelete
			}
			return 500, errorUpdateDistrict
		}
	} else if stringSalesZone == ZoneRegion {
		var regionData modelRegion.Regions
		r := dataZoneInterface.(regionRelation.RegionRelation)
		fmt.Println("r.DistrictID: ", r.DistrictID)
		regionData.ID = r.ID
		regionData.Name = r.Name
		regionData.CreatedAt = r.CreatedAt
		regionData.UpdatedAt = time.Now()
		regionData.DeletedAt = r.DeletedAt
		regionData.IsVacant = false
		regionData.DistrictID = r.DistrictID

		_, errorUpdateRegion := service.regionRepository.Update(regionData)
		if errorUpdateRegion != nil {
			// given rollback to delete the data from table user zone when event error (500) update on table regions
			httpDeleteUserZone, messageErrorDelete := service.userZoneRepository.DeleteUserZoneByID(userIDZoneString)
			if httpDeleteUserZone == 500 && messageErrorDelete != nil {
				return 500, messageErrorDelete
			}
			return 500, nil
		}
	} else if stringSalesZone == ZoneArea {
		var areaData areaModel.Areas
		a := dataZoneInterface.(areaRelation.AreaRelation)
		areaData.ID = a.ID
		areaData.Name = a.Name
		areaData.CreatedAt = a.CreatedAt
		areaData.UpdatedAt = time.Now()
		areaData.DeletedAt = a.DeletedAt
		areaData.IsVacant = false
		areaData.RegionID = a.RegionID

		_, errorUpdateArea := service.areaRepository.Update(areaData)
		if errorUpdateArea != nil {
			// given rollback to delete the data from table user zone when event error (500) update on table areas
			httpDeleteUserZone, messageErrorDelete := service.userZoneRepository.DeleteUserZoneByID(userIDZoneString)
			if httpDeleteUserZone == 500 && messageErrorDelete != nil {
				return 500, messageErrorDelete
			}
			return 500, errorUpdateArea
		}
	} else if stringSalesZone == ZoneGT {
		var gtData gt.GroupTerritories
		g := dataZoneInterface.(groupTerritoriesRelation.GroupTerritoriesRelation)
		gtData.ID = g.ID
		gtData.Name = g.Name
		gtData.CreatedAt = g.CreatedAt
		gtData.UpdatedAt = time.Now()
		gtData.DeletedAt = g.DeletedAt
		gtData.IsVacant = false
		gtData.AreaID = g.AreaID

		_, errorUpdateGT := service.gtRepository.Update(gtData)
		if errorUpdateGT != nil {
			// given rollback to delete the data from table user zone when event error (500) update on table gt
			httpDeleteUserZone, messageErrorDelete := service.userZoneRepository.DeleteUserZoneByID(userIDZoneString)
			if httpDeleteUserZone == 500 && messageErrorDelete != nil {
				return 500, messageErrorDelete
			}
			return 500, errorUpdateGT
		}
	}

	return 200, nil
}

func (service *userZoneService) GetUserZoneByUserId(userId *string) (uint, error, *model.GetBySalesZoneIdUserIdRawQuery) {
	httpStatus, errorUserZone, data := service.userZoneRepository.GetUserZoneByUserId(userId)

	if httpStatus == 500 || errorUserZone != nil {
		return 500, errorUserZone, nil
	} else if data.ID == nil {
		return 404, errors.New("data user not found"), nil
	}

	return httpStatus, errorUserZone, data
}

func (service *userZoneService) UpdateFinishedAssigment(userId *string, userZoneDTO *model.UserZoneRequestDTO) (uint, error) {
	var enumSalesZone model.SalesZoneType

	if *userZoneDTO.SalesZoneType == ZoneGT {
		enumSalesZone = model.GroupTerritory
	} else if *userZoneDTO.SalesZoneType == ZoneArea {
		enumSalesZone = model.Area
	} else if *userZoneDTO.SalesZoneType == ZoneRegion {
		enumSalesZone = model.Region
	} else {
		enumSalesZone = model.District
	}

	// validation zone from dry function
	httpStatusValidation, errorValidation := service.dryValidationZone(int(*userZoneDTO.SalesZoneId), *userZoneDTO.SalesZoneType)
	if httpStatusValidation == 500 && errorValidation != nil {
		return 500, errorValidation
	} else if httpStatusValidation == 404 && errorValidation != nil {
		return 404, errorValidation
	}

	salesZoneIDString := strconv.FormatUint(uint64(*userZoneDTO.SalesZoneId), 10)

	httpStatus, errorUserZone, data := service.userZoneRepository.GetUserZoneByUserIdZoneId(userId, &salesZoneIDString)
	fmt.Println("data: ", data)
	if errorUserZone != nil {
		return httpStatus, errorUserZone
	} else if data.ID == nil {
		return 404, errors.New("data user not found")
	}

	httpStatusUpdateAssigned, errorUpdateAssigned, _ := service.userZoneRepository.UpdateUserAndZoneByUserIdZoneId(userId, &salesZoneIDString, &enumSalesZone)
	if httpStatusUpdateAssigned == 500 {
		return 500, errorUpdateAssigned
	}
	return 200, nil
}

func (service *userZoneService) GetListUserByZoneId(zoneId *string, zoneType *string) (uint, error, *model.GetListUserByZoneRawQuery) {
	var enumSalesZone model.SalesZoneType
	salesZoneIDUint64, errConvertStringToUint := strconv.ParseUint(*zoneId, 10, 32)
	if errConvertStringToUint != nil {
		return 500, errConvertStringToUint, nil
	}

	// validation zone from dry function
	httpStatusValidation, errorValidation := service.dryValidationZone(int(salesZoneIDUint64), *zoneType)
	if httpStatusValidation == 500 && errorValidation != nil {
		return 500, errorValidation, nil
	} else if httpStatusValidation == 404 && errorValidation != nil {
		return 404, errorValidation, nil
	}

	if *zoneType == ZoneGT {
		enumSalesZone = model.GroupTerritory
	} else if *zoneType == ZoneArea {
		enumSalesZone = model.Area
	} else if *zoneType == ZoneRegion {
		enumSalesZone = model.Region
	} else {
		enumSalesZone = model.District
	}

	httpStatus, errorMessage, data := service.userZoneRepository.GetListUserByZoneId(zoneId, &enumSalesZone)
	if httpStatus == 500 && errorMessage != nil {
		return 500, errorMessage, nil
	}

	return httpStatus, errorMessage, data
}

func (service *userZoneService) GetUserZoneByUserIDZoneIDZoneType(userId *string, zoneType *string, zoneId *string) (uint, error, *model.UserZone) {
	var enumSalesZone model.SalesZoneType
	salesZoneIDUint64, errConvertStringToUint := strconv.ParseUint(*zoneId, 10, 32)
	if errConvertStringToUint != nil {
		return 500, errConvertStringToUint, nil
	}

	// validation zone from dry function
	httpStatusValidation, errorValidation := service.dryValidationZone(int(salesZoneIDUint64), *zoneType)
	if httpStatusValidation == 500 && errorValidation != nil {
		return 500, errorValidation, nil
	} else if httpStatusValidation == 404 && errorValidation != nil {
		return 404, errorValidation, nil
	}

	if *zoneType == ZoneGT {
		enumSalesZone = model.GroupTerritory
	} else if *zoneType == ZoneArea {
		enumSalesZone = model.Area
	} else if *zoneType == ZoneRegion {
		enumSalesZone = model.Region
	} else {
		enumSalesZone = model.District
	}
	data, errorMessage := service.userZoneRepository.GetUserZoneByUserIDZoneIDZoneType(userId, &enumSalesZone, zoneId)
	if errorMessage != nil {
		return 500, errorMessage, nil
	} else if data.ID == nil || data.Users.ID == nil {
		return 404, errors.New(fmt.Sprintf("The data with User ID %v with Sales Zone ID %v and Type %v doesn't exist", *userId, *zoneId, *zoneType)), nil
	}

	return 200, nil, &data
}

func (s *userZoneService) GetChildByUserId(userId string, roleName string) (*[]model.UserResponse, error, int) {

	var (
		result        *[]model.UserZone
		resultErr     error
		statusCode    int
		salesAreaId   []uint
		salesRegionId []uint
		salesGTId     []uint
		areaModel     []areaModel.Areas
		gtModel       []gt.GroupTerritories
		resultUsers   []model.UserResponse
	)
	/*
		Cek User exist in table user
	*/
	_, err, userData := s.userRepository.GetById(&userId)
	if err != nil {
		return nil, err, 500
	}
	if userData.ID == nil {
		fmt.Println("---User ID Not Found---")
		return nil, errors.New("User ID"), 404
	}
	/*
		Get user in table user_zone by user id
		return salesZoneType and salesZoneId
	*/
	user, err := s.userZoneRepository.GetBySalesZoneUserID(userId)
	if err != nil {
		return nil, err, 500
	}

	/*
		Check user is exist on table user_zone
		if not exist user role name can be head officer (tt, director and super admin)
	*/
	if user.ID == nil {
		/*
			if user role_id = 7,6,5,4 (ff,asm,sm,nsm)
			show error  "user not asigned to any zone"
		*/
		for i := 0; i < len(RoleIds); i++ {
			if RoleIds[i] == int(*userData.RoleId) {
				fmt.Println("---Zone of user not found---")
				return nil, nil, 200
			}
		}
		/*
			get role_id from table role by logged-in user's roleName
		*/
		role, err := s.userRoleRepository.GetByName(roleName)
		if err != nil {
			return nil, err, 500
		}

		if role.ID == nil {
			return nil, errors.New("Role"), 404
		}
		/*
			get users by role_id
			return array of users
			exmp : all user with role_id = 3 (team tread)
		*/
		statusCode, errUser, dataUser := s.userRepository.GetByRoleId(*role.ID)
		if err != nil {
			return nil, errUser, int(statusCode)
		}

		for _, value := range *dataUser {
			resultUser := model.UserResponse{
				UserId:   *value.ID,
				UserName: value.Name,
				Nip:      value.Nip,
				RoleID:   value.RoleId,
				RoleName: value.Role.Label,
			}
			resultUsers = append(resultUsers, resultUser)
		}

		return &resultUsers, nil, 200
	}

	/*
		IF user id exist on table user_zone
	*/
	salesZoneId := user.SalesZoneId
	salesZoneType := user.SalesZoneType
	/*
		Check sales zone type is regions, areas or others
	*/

	stringSalesZone := fmt.Sprintf("%v", salesZoneType)
	switch stringSalesZone {
	case ZoneDistric:
		/*
			get data from table region by sales_zone_id
		*/
		districID := int(*salesZoneId)
		districs, err := s.districtRepository.FindById(districID)
		if err != nil {
			return nil, err, 500
		}
		if districs.ID == 0 {
			fmt.Println("----District ID Not Found----")
			return nil, nil, 200
		}
		/*
			get area id and append in array
		*/

		for _, value := range districs.RegionsUnderDistrict {
			salesRegionId = append(salesRegionId, value.ID)
			areaModel = append(areaModel, value.AreasUnderRegion...)
		}
		for _, value := range areaModel {
			salesAreaId = append(salesAreaId, value.ID)
			gtModel = append(gtModel, value.GtUnderArea...)
		}

		switch roleName {
		case Sm:
			resultZone, err := s.userZoneRepository.GetBySalesZoneIDMultiple(salesAreaId, ZoneRegion, roleName)
			result = &resultZone
			resultErr = err
			statusCode = 200
		case Asm:
			resultZone, err := s.userZoneRepository.GetBySalesZoneIDMultiple(salesAreaId, ZoneArea, roleName)
			result = &resultZone
			resultErr = err
			statusCode = 200
		case FieldForce:
			for _, value := range gtModel {
				salesGTId = append(salesGTId, value.ID)
			}
			resultZone, err := s.userZoneRepository.GetBySalesZoneIDMultiple(salesGTId, ZoneGT, roleName)
			result = &resultZone
			resultErr = err
			statusCode = 200
		default:
			return nil, errors.New("role name is not child of user id"), 400
		}
	case ZoneRegion:
		/*
			get data from table region by sales_zone_id
		*/
		regionId := int(*salesZoneId)
		region, err := s.regionRepository.FindById(regionId)
		if err != nil {
			return nil, err, 500
		}
		if region.ID == 0 {
			fmt.Println("----Region ID Not Found----")
			return nil, nil, 200
		}
		/*
			get area id and append in array
		*/
		for _, value := range region.AreasUnderRegion {
			salesAreaId = append(salesAreaId, value.ID)
			gtModel = append(gtModel, value.GtUnderArea...)
		}

		switch roleName {
		case Asm:
			resultZone, err := s.userZoneRepository.GetBySalesZoneIDMultiple(salesAreaId, ZoneArea, roleName)
			result = &resultZone
			resultErr = err
			statusCode = 200
		case FieldForce:
			//tampilkan array gt
			for _, value := range gtModel {
				salesGTId = append(salesGTId, value.ID)
			}
			resultZone, err := s.userZoneRepository.GetBySalesZoneIDMultiple(salesGTId, ZoneGT, roleName)
			result = &resultZone
			resultErr = err
			statusCode = 200
		default:
			return nil, errors.New("role name is not child of user id"), 400
		}

	case ZoneArea:
		//get data dari table area berdasarkan salesZoneId
		areaId := int(*salesZoneId)
		area, err := s.areaRepository.FindById(areaId)
		if err != nil {
			return nil, err, 500
		}
		if area.ID == 0 {
			fmt.Println("----Area ID Not Found----")
			return nil, nil, 200
		}
		//tampilkan array gt
		switch roleName {
		case FieldForce:
			//ambil array gt
			for _, value := range area.GtUnderArea {
				salesGTId = append(salesGTId, value.ID)
			}
			resultZone, err := s.userZoneRepository.GetBySalesZoneIDMultiple(salesGTId, ZoneGT, roleName)
			result = &resultZone
			resultErr = err
			statusCode = 200
		default:
			return nil, errors.New("role name is not child of user id"), 400
		}
	}

	if result == nil {
		return nil, nil, 200
	}
	if len(*result) == 0 {
		return nil, nil, 200
	}

	for _, value := range *result {
		resultUser := model.UserResponse{
			UserId:        *value.UserId,
			UserName:      value.Users.Name,
			Nip:           value.Users.Nip,
			RoleID:        value.Users.RoleId,
			RoleName:      value.Users.Role.Label,
			SalesZoneType: string(value.SalesZoneType),
			SalesZoneId:   value.SalesZoneId,
		}
		resultUsers = append(resultUsers, resultUser)
	}

	/*
		Delete duplicate data
	*/
	uniq := make(map[string]bool)
	var uniqUser []model.UserResponse

	for _, user := range resultUsers {
		key := *user.UserName + string(user.UserId)
		if !uniq[key] {
			uniq[key] = true
			uniqUser = append(uniqUser, user)
		}
	}

	return &uniqUser, resultErr, statusCode

}

func (s *userZoneService) GetChildVacantByUserId(userId string) (*[]model.UserResponse, error, int) {
	var (
		result        *[]model.UserZone
		resultErr     error
		statusCode    int
		salesAreaId   []uint
		salesRegionId []uint
		salesGTId     []uint
		areaModel     []areaModel.Areas
		gtModel       []gt.GroupTerritories
		resultData    []model.UserZone
		resultUsers   []model.UserResponse
	)
	/*
		Cek User exist in table user
	*/
	_, err, userData := s.userRepository.GetById(&userId)
	if err != nil {
		return nil, err, 500
	}
	if userData.ID == nil {
		fmt.Println("---User ID Not Found---")
		return nil, errors.New("User ID"), 404
	}
	/*
		Get user in table user_zone by user id
		return salesZoneType and salesZoneId
	*/
	user, err := s.userZoneRepository.GetBySalesZoneUserID(userId)
	if err != nil {
		return nil, err, 500
	}

	/*
		Check user is exist on table user_zone
		if not exist user role name can be head officer (nsm, tt, director and super admin)
	*/
	if user.ID == nil {
		/*
			if user role_id = 7,6,5,4 (ff,asm,sm,nsm)
			show error  "user not asigned to any zone"
		*/
		for i := 0; i < len(RoleIds); i++ {
			if RoleIds[i] == int(*userData.RoleId) {
				fmt.Println("---Zone of user not found---")
				return nil, nil, 200
			}
		}
		/*
			get all user vacant
		*/
		userVacant, err := s.userZoneRepository.GetAllUserVacant()
		if err != nil {
			return nil, err, 500
		}

		if len(userVacant) == 0 {
			fmt.Println("---The user has no subordinates---")
			return nil, nil, 200
		}
		/*
			return array of users
		*/
		for _, value := range userVacant {
			resultUser := model.UserResponse{
				UserId:        *value.UserId,
				UserName:      value.Users.Name,
				Nip:           value.Users.Nip,
				RoleID:        value.Users.RoleId,
				RoleName:      value.Users.Role.Label,
				SalesZoneType: string(value.SalesZoneType),
				SalesZoneId:   value.SalesZoneId,
			}
			resultUsers = append(resultUsers, resultUser)
		}

		return &resultUsers, nil, 200
	}

	/*
		IF user id exist on table user_zone
	*/
	salesZoneId := user.SalesZoneId
	salesZoneType := user.SalesZoneType
	/*
		Check sales zone type is regions, areas or others
	*/
	stringSalesZone := fmt.Sprintf("%v", salesZoneType)
	switch stringSalesZone {
	case ZoneDistric:
		/*
			get data from table Distric by sales_zone_id
		*/
		districtID := int(*salesZoneId)
		districs, err := s.districtRepository.FindById(districtID)
		if err != nil {
			return nil, err, 500
		}
		if districs.ID == 0 {
			fmt.Println("----District ID Not Found----")
			return nil, nil, 200
		}
		/*
			get area id, region id and gt id then append in array
		*/

		//Region
		for _, value := range districs.RegionsUnderDistrict {
			salesRegionId = append(salesRegionId, value.ID)
			areaModel = append(areaModel, value.AreasUnderRegion...)
		}
		resultZone, err := s.userZoneRepository.GetUserVacantBySalesZoneIDMultiple(salesRegionId, ZoneRegion)
		resultData = append(resultData, resultZone...)

		//Area
		for _, value := range areaModel {
			salesAreaId = append(salesAreaId, value.ID)
		}
		resultArea, err := s.userZoneRepository.GetUserVacantBySalesZoneIDMultiple(salesAreaId, ZoneArea)
		resultData = append(resultData, resultArea...)

		//GT
		for _, value := range gtModel {
			salesGTId = append(salesGTId, value.ID)
		}
		resultGT, err := s.userZoneRepository.GetUserVacantBySalesZoneIDMultiple(salesGTId, ZoneGT)
		resultData = append(resultData, resultGT...)

		//return response
		result = &resultData
		resultErr = err
		statusCode = 200

	case ZoneRegion:
		/*
			get data from table region by sales_zone_id
		*/
		regionId := int(*salesZoneId)
		region, err := s.regionRepository.FindById(regionId)
		if err != nil {
			return nil, err, 500
		}
		if region.ID == 0 {
			fmt.Println("----Region ID Not Found----")
			return nil, nil, 200
		}
		/*
			get area id and GT id append in array
		*/
		//Area
		for _, value := range region.AreasUnderRegion {
			salesAreaId = append(salesAreaId, value.ID)
			gtModel = append(gtModel, value.GtUnderArea...)
		}
		resultZone, err := s.userZoneRepository.GetUserVacantBySalesZoneIDMultiple(salesAreaId, ZoneArea)
		resultData = append(resultData, resultZone...)

		//GT
		for _, value := range gtModel {
			salesGTId = append(salesGTId, value.ID)
		}
		resultGT, err := s.userZoneRepository.GetUserVacantBySalesZoneIDMultiple(salesGTId, ZoneGT)
		resultData = append(resultData, resultGT...)

		//Return result
		result = &resultData
		resultErr = err
		statusCode = 200

	case ZoneArea:
		//get data dari table area berdasarkan salesZoneId
		areaId := int(*salesZoneId)
		area, err := s.areaRepository.FindById(areaId)
		if err != nil {
			return nil, err, 500
		}
		if area.ID == 0 {
			fmt.Println("----Araea ID Not Found----")
			return nil, nil, 200
		}
		//GT
		for _, value := range area.GtUnderArea {
			salesGTId = append(salesGTId, value.ID)
		}
		resultGT, err := s.userZoneRepository.GetUserVacantBySalesZoneIDMultiple(salesGTId, ZoneGT)
		resultData = append(resultData, resultGT...)

		//Return Result
		result = &resultData
		resultErr = err
		statusCode = 200

	default:
		fmt.Println("---The user has no subordinates---")
		return nil, nil, 200
	}

	if len(*result) == 0 {
		return nil, nil, 200
	}

	for _, value := range *result {
		resultUser := model.UserResponse{
			UserId:        *value.UserId,
			UserName:      value.Users.Name,
			Nip:           value.Users.Nip,
			RoleID:        value.Users.RoleId,
			RoleName:      value.Users.Role.Label,
			SalesZoneType: string(value.SalesZoneType),
			SalesZoneId:   value.SalesZoneId,
		}
		resultUsers = append(resultUsers, resultUser)
	}

	/*
		Delete duplicate data
	*/
	uniq := make(map[string]bool)
	var uniqUser []model.UserResponse

	for _, user := range resultUsers {
		key := *user.UserName + string(rune(user.UserId))
		if !uniq[key] {
			uniq[key] = true
			uniqUser = append(uniqUser, user)
		}
	}

	return &uniqUser, resultErr, statusCode

}

func (service *userZoneService) GetSubordinateEmployeesByUserIDZoneIDZoneType(userId *string, zoneType *string, zoneId *string, roleName *string) (uint, error, *[]model.UserZoneWithEpochEntity) {
	var enumSalesZone model.SalesZoneType
	salesZoneIDUint64, errConvertStringToUint := strconv.ParseUint(*zoneId, 10, 32)
	if errConvertStringToUint != nil {
		return 500, errConvertStringToUint, nil
	}

	// validation zone from dry function
	httpStatusValidation, errorValidation := service.dryValidationZone(int(salesZoneIDUint64), *zoneType)
	if httpStatusValidation == 500 && errorValidation != nil {
		return 500, errorValidation, nil
	} else if httpStatusValidation == 404 && errorValidation != nil {
		return 404, errorValidation, nil
	}

	if *zoneType == ZoneGT {
		enumSalesZone = model.GroupTerritory
	} else if *zoneType == ZoneArea {
		enumSalesZone = model.Area
	} else if *zoneType == ZoneRegion {
		enumSalesZone = model.Region
	} else {
		enumSalesZone = model.District
	}
	httpStatus, errorMessage, data := service.userZoneRepository.GetSubordinateEmployeesByUserIDZoneIDZoneType(userId, &enumSalesZone, zoneId, roleName)

	if httpStatus == 500 && errorMessage != nil {
		return 500, errorMessage, nil
	} else if len(*data) == 0 {
		return 404, errors.New(fmt.Sprintf("The data with User ID %v with Sales Zone ID %v, Type %v, with Role Name %v doesn't exist", *userId, *zoneId, *zoneType, *roleName)), nil
	}

	return 200, nil, data
}

func (s *userZoneService) GetChildNonVacantByUserId(userId string) (*[]model.UserResponse, error, int) {
	var (
		result        *[]model.UserZone
		resultErr     error
		statusCode    int
		salesAreaId   []uint
		salesRegionId []uint
		salesGTId     []uint
		areaModel     []areaModel.Areas
		gtModel       []gt.GroupTerritories
		resultData    []model.UserZone
		resultUsers   []model.UserResponse
	)
	/*
		Cek User exist in table user
	*/
	_, err, userData := s.userRepository.GetById(&userId)
	if err != nil {
		return nil, err, 500
	}
	if userData.ID == nil {
		fmt.Println("---User ID Not Found---")
		return nil, errors.New("User ID"), 404
	}
	/*
		Get user in table user_zone by user id
		return salesZoneType and salesZoneId
	*/
	user, err := s.userZoneRepository.GetBySalesZoneUserID(userId)
	if err != nil {
		return nil, err, 500
	}

	/*
		Check user is exist on table user_zone
		if not exist user role name can be head officer (tt, director and super admin)
	*/
	if user.ID == nil {
		/*
			if user role_id = 7,6,5,4 (ff,asm,sm,nsm)
			show error  "user not asigned to any zone"
		*/

		for i := 0; i < len(RoleIds); i++ {
			if RoleIds[i] == int(*userData.RoleId) {
				fmt.Println("---Zone of user not found---")
				return nil, nil, 200
			}
		}
		/*
			get all user non vacant
		*/
		userNonVacant, err := s.userZoneRepository.GetAllUserNonVacant()
		if err != nil {
			return nil, err, 500
		}

		if len(userNonVacant) == 0 {
			fmt.Println("---The user has no subordinates---")
			return nil, nil, 200
		}
		/*
			return array of users
		*/
		for _, value := range userNonVacant {
			resultUser := model.UserResponse{
				UserId:        *value.UserId,
				UserName:      value.Users.Name,
				Nip:           value.Users.Nip,
				RoleID:        value.Users.RoleId,
				RoleName:      value.Users.Role.Label,
				SalesZoneType: string(value.SalesZoneType),
				SalesZoneId:   value.SalesZoneId,
			}
			resultUsers = append(resultUsers, resultUser)
		}

		return &resultUsers, nil, 200
	}

	/*
		IF user id exist on table user_zone
	*/
	salesZoneId := user.SalesZoneId
	salesZoneType := user.SalesZoneType
	/*
		Check sales zone type is regions, areas or others
	*/
	stringSalesZone := fmt.Sprintf("%v", salesZoneType)
	switch stringSalesZone {
	case ZoneDistric:
		/*
			get data from table Distric by sales_zone_id
		*/
		districtID := int(*salesZoneId)
		districs, err := s.districtRepository.FindById(districtID)
		if err != nil {
			return nil, err, 500
		}
		if districs.ID == 0 {
			fmt.Println("----District ID Not Found----")
			return nil, nil, 200
		}
		/*
			get area id, region id and gt id then append in array
		*/

		//Region
		for _, value := range districs.RegionsUnderDistrict {
			salesRegionId = append(salesRegionId, value.ID)
			areaModel = append(areaModel, value.AreasUnderRegion...)
		}
		resultZone, err := s.userZoneRepository.GetUserNonVacantBySalesZoneIDMultiple(salesRegionId, ZoneRegion)
		resultData = append(resultData, resultZone...)

		//Area
		for _, value := range areaModel {
			salesAreaId = append(salesAreaId, value.ID)
		}
		resultArea, err := s.userZoneRepository.GetUserNonVacantBySalesZoneIDMultiple(salesAreaId, ZoneArea)
		resultData = append(resultData, resultArea...)

		//GT
		for _, value := range gtModel {
			salesGTId = append(salesGTId, value.ID)
		}
		resultGT, err := s.userZoneRepository.GetUserNonVacantBySalesZoneIDMultiple(salesGTId, ZoneGT)
		resultData = append(resultData, resultGT...)

		//return response
		result = &resultData
		resultErr = err
		statusCode = 200

	case ZoneRegion:
		/*
			get data from table region by sales_zone_id
		*/
		regionId := int(*salesZoneId)
		region, err := s.regionRepository.FindById(regionId)
		if err != nil {
			return nil, err, 500
		}
		if region.ID == 0 {
			fmt.Println("----Region ID Not Found----")
			return nil, nil, 200
		}
		/*
			get area id and GT id and append in array
		*/

		//Area
		for _, value := range region.AreasUnderRegion {
			salesAreaId = append(salesAreaId, value.ID)
			gtModel = append(gtModel, value.GtUnderArea...)
		}

		resultZone, err := s.userZoneRepository.GetUserNonVacantBySalesZoneIDMultiple(salesAreaId, ZoneArea)
		resultData = append(resultData, resultZone...)

		//GT
		for _, value := range gtModel {
			salesGTId = append(salesGTId, value.ID)
		}
		resultGT, err := s.userZoneRepository.GetUserNonVacantBySalesZoneIDMultiple(salesGTId, ZoneGT)
		resultData = append(resultData, resultGT...)

		//Return result
		result = &resultData
		resultErr = err
		statusCode = 200

	case ZoneArea:
		/*
			get data from table area by sales_zone_id
		*/
		areaId := int(*salesZoneId)
		area, err := s.areaRepository.FindById(areaId)
		if err != nil {
			return nil, err, 500
		}
		if area.ID == 0 {
			fmt.Println("----Area ID Not Found----")
			return nil, nil, 200
		}
		//GT
		for _, value := range area.GtUnderArea {
			salesGTId = append(salesGTId, value.ID)
		}

		resultGT, err := s.userZoneRepository.GetUserNonVacantBySalesZoneIDMultiple(salesGTId, ZoneGT)
		resultData = append(resultData, resultGT...)
		result = &resultData
		resultErr = err
		statusCode = 200
	default:
		fmt.Println("---The user has no subordinates---")
		return nil, nil, 200
	}

	if len(*result) == 0 {
		fmt.Println("---Data Not Found---")
		return nil, nil, 200
	}

	for _, value := range *result {
		resultUser := model.UserResponse{
			UserId:        *value.UserId,
			UserName:      value.Users.Name,
			Nip:           value.Users.Nip,
			RoleName:      value.Users.Role.Label,
			RoleID:        value.Users.RoleId,
			SalesZoneType: string(value.SalesZoneType),
			SalesZoneId:   value.SalesZoneId,
		}
		resultUsers = append(resultUsers, resultUser)
	}

	/*
		Delete duplicate data
	*/
	uniq := make(map[string]bool)
	var uniqUser []model.UserResponse

	for _, user := range resultUsers {
		key := *user.UserName + string(user.UserId)
		if !uniq[key] {
			uniq[key] = true
			uniqUser = append(uniqUser, user)
		}
	}

	return &uniqUser, resultErr, statusCode

}

func (service *userZoneService) GetZoneByUserID(salesZoneType *string, userID *string) (uint, error, []model.GetZoneByUserIDResponse) {
	var result []model.GetZoneByUserIDResponse
	var resultErr error
	var statusCode uint
	/*
		Cek User exist in table user
	*/
	_, err, userData := service.userRepository.GetById(userID)
	if err != nil {
		return 500, err, nil
	}
	if userData.ID == nil {
		return 404, errors.New(fmt.Sprintf("user id from access token isn't found %v", userID)), nil
	}
	/*
		Get user in table user_zone by user id
		return salesZoneType and salesZoneId
	*/
	user, err := service.userZoneRepository.GetBySalesZoneUserID(*userID)
	if err != nil {
		return 500, err, nil
	}
	fmt.Println(user)
	/*
		IF params sales_zone_type is nil
	*/
	if salesZoneType == nil {
		var enumSaleZoneType string
		if user.SalesZoneType == "districts" {
			enumSaleZoneType = "districts"
		} else if user.SalesZoneType == "regions" {
			enumSaleZoneType = "regions"
		} else if user.SalesZoneType == "areas" {
			enumSaleZoneType = "areas"
		} else if user.SalesZoneType == "group_territories" {
			enumSaleZoneType = "group_territories"
		}
		salesZoneType = &enumSaleZoneType
	}
	if user.SalesZoneType == "regions" {
		if *salesZoneType == "districts" {
			return 400, errors.New(fmt.Sprintf("user role is %v only access sales zone type regions, areas, and group_territories", *userData.Role.Name)), nil
		}
	} else if user.SalesZoneType == "areas" {
		if *salesZoneType == "districts" || *salesZoneType == "regions" {
			return 400, errors.New(fmt.Sprintf("user role is %v only access sales zone type areas, and group_territories", *userData.Role.Name)), nil
		}
	} else if user.SalesZoneType == "group_territories" {
		if *salesZoneType == "districts" || *salesZoneType == "regions" || *salesZoneType == "areas" {
			return 400, errors.New(fmt.Sprintf("user role is %v only access sales zone type group_territories", *userData.Role.Name)), nil
		}
	}
	/*
		IF user id exist on table user_zone
	*/
	userSalesZoneId := user.SalesZoneId
	userSalesZoneType := user.SalesZoneType
	fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, userSalesZoneId: ", *userSalesZoneId)
	fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, userSalesZoneType: ", userSalesZoneType)
	fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, salesZoneType: ", *salesZoneType)
	/*
		Check sales zone type is regions, areas or others
	*/
	switch userSalesZoneType {
	case "districts":
		districtID := int(*userSalesZoneId)
		dataDistrict, errDistrict := service.districtRepository.FindById(districtID)
		if errDistrict != nil {
			return 500, errDistrict, nil
		}

		var DataDistrict []model.GetZoneByUserIDResponse
		if dataDistrict.ID == 0 {
			DataDistrict = []model.GetZoneByUserIDResponse{}
		} else {
			DataDistrict = append(DataDistrict, model.GetZoneByUserIDResponse{
				SalesZoneID:   dataDistrict.ID,
				SalesZoneType: *salesZoneType,
				SalesZoneName: dataDistrict.Name,
			})
		}

		// check user request by params
		if *salesZoneType == "districts" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Districts and Type Filter is Districts")
			// given back to response
			result = DataDistrict
			resultErr = nil
			statusCode = 200
		} else if *salesZoneType == "regions" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Districts and Type Filter is Regions")
			var dataRegionUnderDistrict []model.GetZoneByUserIDResponse
			if len(dataDistrict.RegionsUnderDistrict) == 0 {
				dataRegionUnderDistrict = []model.GetZoneByUserIDResponse{}
			} else {
				for _, valueRegionUnderDistrict := range dataDistrict.RegionsUnderDistrict {
					dataRegionUnderDistrict = append(dataRegionUnderDistrict, model.GetZoneByUserIDResponse{
						SalesZoneID:   valueRegionUnderDistrict.ID,
						SalesZoneType: *salesZoneType,
						SalesZoneName: valueRegionUnderDistrict.Name,
					})
				}
			}
			// given back to response
			result = dataRegionUnderDistrict
			resultErr = nil
			statusCode = 200
		} else if *salesZoneType == "areas" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Districts and Type Filter is Areas")
			var dataRegionUnderDistrictID []int
			if len(dataDistrict.RegionsUnderDistrict) == 0 {
				return 500, errors.New(fmt.Sprintf("ups, RegionID doens't have child on table Area with DisctirctID %v", dataDistrict.ID)), nil
			} else {
				for _, value := range dataDistrict.RegionsUnderDistrict {
					dataRegionUnderDistrictID = append(dataRegionUnderDistrictID, int(value.ID))
				}
			}
			dataAreas, errorDataAreas := service.areaRepository.FindAreaByRegionID(dataRegionUnderDistrictID)
			if errorDataAreas != nil {
				fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, Error when Get Areas 1 - (dataRegionUnderDistrictID)")
				return 500, errorDataAreas, nil
			}
			var dataAreaUnderRegion []model.GetZoneByUserIDResponse
			if len(dataAreas) == 0 {
				dataAreaUnderRegion = []model.GetZoneByUserIDResponse{}
			} else {
				for _, value := range dataAreas {
					dataAreaUnderRegion = append(dataAreaUnderRegion, model.GetZoneByUserIDResponse{
						SalesZoneID:   value.ID,
						SalesZoneType: *salesZoneType,
						SalesZoneName: value.Name,
					})
				}
			}
			// given back to response
			result = dataAreaUnderRegion
			resultErr = nil
			statusCode = 200
		} else if *salesZoneType == "group_territories" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Districts and Type Filter is Group_Territories")
			var dataRegionUnderDistrictID []int
			if len(dataDistrict.RegionsUnderDistrict) == 0 {
				return 500, errors.New(fmt.Sprintf("ups, RegionID doens't have child on table Area with DisctirctID %v", dataDistrict.ID)), nil
			} else {
				for _, value := range dataDistrict.RegionsUnderDistrict {
					dataRegionUnderDistrictID = append(dataRegionUnderDistrictID, int(value.ID))
				}
			}
			dataAreas, errorDataAreas := service.areaRepository.FindAreaByRegionID(dataRegionUnderDistrictID)
			if errorDataAreas != nil {
				fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, Error when Get Areas 2 - (dataRegionUnderDistrictID)")
				return 500, errorDataAreas, nil
			}
			var dataAreaUnderRegionID []int
			if len(dataAreas) == 0 {
				return 500, errors.New(fmt.Sprintf("ups, AreaID doens't have child on table Area with Region ID %v", dataRegionUnderDistrictID)), nil
			} else {
				for _, value := range dataAreas {
					dataAreaUnderRegionID = append(dataAreaUnderRegionID, int(value.ID))
				}
			}

			dataGT, errorDataAreaUnderRegionID := service.gtRepository.FindGTByAreaID(dataAreaUnderRegionID)
			if errorDataAreaUnderRegionID != nil {
				return 500, errorDataAreaUnderRegionID, nil
			}
			var dataGTUnderAreaID []model.GetZoneByUserIDResponse
			if len(dataGT) == 0 {
				dataGTUnderAreaID = []model.GetZoneByUserIDResponse{}
			} else {
				for _, value := range dataGT {
					dataGTUnderAreaID = append(dataGTUnderAreaID, model.GetZoneByUserIDResponse{
						SalesZoneID:   value.ID,
						SalesZoneType: *salesZoneType,
						SalesZoneName: value.Name,
					})
				}
			}
			// given back to response
			result = dataGTUnderAreaID
			resultErr = nil
			statusCode = 200
		}
	case "regions":
		fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Regions")
		regionId := int(*userSalesZoneId)
		dataRegion, errRegion := service.regionRepository.FindById(regionId)
		if errRegion != nil {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, Error when Get Region 1 - (regionId)")
			return 500, errRegion, nil
		}

		var DataRegion []model.GetZoneByUserIDResponse
		if dataRegion.ID == 0 {
			DataRegion = []model.GetZoneByUserIDResponse{}
		} else {
			DataRegion = append(DataRegion, model.GetZoneByUserIDResponse{
				SalesZoneID:   dataRegion.ID,
				SalesZoneType: *salesZoneType,
				SalesZoneName: dataRegion.Name,
			})
		}

		// check user request by params
		if *salesZoneType == "regions" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Regions and Type Filter is Regions")
			// given back to response
			result = DataRegion
			resultErr = nil
			statusCode = 200
		} else if *salesZoneType == "areas" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Regions and Type Filter is Areas")
			var DataAreaUnderRegionID []int
			if len(dataRegion.AreasUnderRegion) == 0 {
				return 500, errors.New(fmt.Sprintf("ups, RegionID doens't have child on table Area with Region ID %v", dataRegion.ID)), nil
			} else {
				for _, value := range dataRegion.AreasUnderRegion {
					DataAreaUnderRegionID = append(DataAreaUnderRegionID, int(value.RegionID))
				}
			}

			var dataAreaByRegionID []model.GetZoneByUserIDResponse

			dataAreaByRegionIDRepo, errorMessageDataAreaByRegionID := service.areaRepository.FindAreaByRegionID(DataAreaUnderRegionID)
			if errorMessageDataAreaByRegionID != nil {
				return 500, errorMessageDataAreaByRegionID, nil
			}

			if len(dataAreaByRegionIDRepo) == 0 {
				dataAreaByRegionID = []model.GetZoneByUserIDResponse{}
			} else {
				for _, valueArea := range dataAreaByRegionIDRepo {
					dataAreaByRegionID = append(dataAreaByRegionID, model.GetZoneByUserIDResponse{
						SalesZoneID:   valueArea.ID,
						SalesZoneType: *salesZoneType,
						SalesZoneName: valueArea.Name,
					})
				}
			}
			// given back to response
			result = dataAreaByRegionID
			resultErr = nil
			statusCode = 200
		} else if *salesZoneType == "group_territories" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Regions and Type Filter is GT")
			var DataAreaUnderRegionID []int
			if len(dataRegion.AreasUnderRegion) == 0 {
				return 500, errors.New(fmt.Sprintf("ups, RegionID doens't have child on table Area with Region ID %v", dataRegion.ID)), nil
			} else {
				for _, value := range dataRegion.AreasUnderRegion {
					DataAreaUnderRegionID = append(DataAreaUnderRegionID, int(value.RegionID))
				}
			}

			var DataAreaID []int
			dataAreaByRegionIDRepo, errorMessageDataAreaByRegionID := service.areaRepository.FindAreaByRegionID(DataAreaUnderRegionID)
			if errorMessageDataAreaByRegionID != nil {
				return 500, errorMessageDataAreaByRegionID, nil
			}

			if len(dataAreaByRegionIDRepo) == 0 {
				return 500, errors.New(fmt.Sprintf("ups, RegionID with List ID: %v on Table Area doesn't exist", DataAreaUnderRegionID)), nil
			} else {
				for _, value := range dataAreaByRegionIDRepo {
					DataAreaID = append(DataAreaID, int(value.ID))
				}
			}

			var dataGTByAreaID []model.GetZoneByUserIDResponse
			dataGTByAreaIDRepo, errorDataGTByAreaIDRepo := service.gtRepository.FindGTByAreaID(DataAreaID)
			if errorDataGTByAreaIDRepo != nil {
				return 500, errorDataGTByAreaIDRepo, nil
			}

			if len(dataGTByAreaIDRepo) == 0 {
				dataGTByAreaID = []model.GetZoneByUserIDResponse{}
			} else {
				for _, value := range dataGTByAreaIDRepo {
					dataGTByAreaID = append(dataGTByAreaID, model.GetZoneByUserIDResponse{
						SalesZoneID:   value.ID,
						SalesZoneType: *salesZoneType,
						SalesZoneName: value.Name,
					})
				}
			}
			// given back to response
			result = dataGTByAreaID
			resultErr = nil
			statusCode = 200
		}
	case "areas":
		fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Areas")
		areaID := int(*userSalesZoneId)
		dataAreas, errAreas := service.areaRepository.FindById(areaID)
		if errAreas != nil {
			return 500, errAreas, nil
		}

		var DataArea []model.GetZoneByUserIDResponse
		if dataAreas.ID == 0 {
			DataArea = []model.GetZoneByUserIDResponse{}
		} else {
			DataArea = append(DataArea, model.GetZoneByUserIDResponse{
				SalesZoneID:   dataAreas.ID,
				SalesZoneType: *salesZoneType,
				SalesZoneName: dataAreas.Name,
			})
		}

		// check user request by params
		if *salesZoneType == "areas" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Areas and Type Filter is Areas")
			// given back to response
			result = DataArea
			resultErr = nil
			statusCode = 200
		} else if *salesZoneType == "group_territories" {
			fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is Areas and Type Filter is GT")
			var DataGTUnderAreaID []int
			if len(dataAreas.GtUnderArea) == 0 {
				return 500, errors.New(fmt.Sprintf("ups, AreaID doens't have child on table GT with Area ID %v", dataAreas.ID)), nil
			} else {
				for _, value := range dataAreas.GtUnderArea {
					DataGTUnderAreaID = append(DataGTUnderAreaID, int(value.AreaID))
				}
			}

			var dataGTByAreaID []model.GetZoneByUserIDResponse
			dataGTByAreaIDRepo, errorDataGTByAreaIDRepo := service.gtRepository.FindGTByAreaID(DataGTUnderAreaID)
			if errorDataGTByAreaIDRepo != nil {
				return 500, errorDataGTByAreaIDRepo, nil
			}

			if len(dataGTByAreaIDRepo) == 0 {
				dataGTByAreaID = []model.GetZoneByUserIDResponse{}
			} else {
				for _, value := range dataGTByAreaIDRepo {
					dataGTByAreaID = append(dataGTByAreaID, model.GetZoneByUserIDResponse{
						SalesZoneID:   value.ID,
						SalesZoneType: *salesZoneType,
						SalesZoneName: value.Name,
					})
				}
			}
			// given back to response
			result = dataGTByAreaID
			resultErr = nil
			statusCode = 200
		}
	case "group_territories":
		fmt.Println("--- USER ZONE SERVICE / GET ZONE BY USER ID AND ZONE TYPE ---, User Zone Type is GT")
		gtID := int(*userSalesZoneId)
		dataGT, errGT := service.gtRepository.FindById(gtID)
		if errGT != nil {
			return 500, errGT, nil
		}

		var DataGT []model.GetZoneByUserIDResponse
		if dataGT.ID == 0 {
			DataGT = []model.GetZoneByUserIDResponse{}
		} else {
			DataGT = append(DataGT, model.GetZoneByUserIDResponse{
				SalesZoneID:   dataGT.ID,
				SalesZoneType: *salesZoneType,
				SalesZoneName: dataGT.Name,
			})
		}
		// given back to response
		result = DataGT
		resultErr = nil
		statusCode = 200
	default:
		result = nil
		resultErr = errors.New("params undefined")
		statusCode = 500
	}

	return statusCode, resultErr, result
}

func (s *userZoneService) GetZoneChildVacantByUserId(userId string) ([]model.ZoneType, error, int) {
	responseCode, err, result := s.userZoneRepository.GetUserZoneByUserId(&userId)
	if err != nil {
		return nil, err, int(responseCode)
	}

	var ResultQuery []model.ZoneType
	if result.SalesZoneId != nil {
		pointerSalesZoneId := int(*result.SalesZoneId)
		ResultQuery, err = s.userZoneRepository.GetZoneChildVacantBySalesZoneData(string(result.SalesZoneType),
			&pointerSalesZoneId)
	} else {
		ResultQuery, err = s.userZoneRepository.GetZoneChildVacantBySalesZoneData(string(result.SalesZoneType),
			nil)
	}

	return ResultQuery, err, http.StatusOK
}

func (s *userZoneService) ImpersonateAccessControlSales(userId string) ([]model.ZoneTypeRole, error, int) {
	responseCode, err, result := s.userZoneRepository.GetUserZoneByUserId(&userId)
	if err != nil {
		return nil, err, int(responseCode)
	}

	var ResultQuery []model.ZoneTypeRole
	if result.SalesZoneId != nil {
		pointerSalesZoneId := int(*result.SalesZoneId)
		ResultQuery, err = s.userZoneRepository.GetZoneChildRoleImpersonate(string(result.SalesZoneType), &pointerSalesZoneId)
	} else {
		ResultQuery, err = s.userZoneRepository.GetZoneChildRoleImpersonate(string(result.SalesZoneType), nil)
	}
	return ResultQuery, nil, http.StatusAccepted
}
