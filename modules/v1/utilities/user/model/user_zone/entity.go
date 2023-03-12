package model

import (
	"database/sql/driver"
	"ethical-be/modules/v1/utilities/user/model"
	"fmt"
	"time"
)

type SalesZoneType string

const (
	GroupTerritory SalesZoneType = "group_territories"
	Area           SalesZoneType = "areas"
	Region         SalesZoneType = "regions"
	District       SalesZoneType = "districts"
)

func (s *SalesZoneType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*s = SalesZoneType(value.([]byte))
	case string:
		*s = SalesZoneType(value.(string))
	default:
		return fmt.Errorf("cannot sql.Scan() SalesZoneType from: %#v", v)
	}
	return nil
}

func (s SalesZoneType) Value() (driver.Value, error) {
	return string(s), nil
}

type UserZone struct {
	ID            *uint `gorm:"primaryKey"`
	UserId        *uint
	SalesZoneId   *uint
	SalesZoneType SalesZoneType `gorm:"type:sales_zone_type"`
	AssignedDate  *time.Time
	FinishedDate  *time.Time
	CreatedAt     *time.Time `gorm:"DEFAULT:current_timestamp"`
	UpdatedAt     *time.Time
	Users         model.Users `gorm:"foreignKey:UserId"`
}

func (UserZone) TableName() string {
	return "user_zones"
}

type GetBySalesZoneIdUserIdRawQuery struct {
	ID            *uint
	UserId        *uint
	Username      *string
	UserNip       *string
	SalesZoneId   *uint
	NameSalesZone *string
	SalesZoneType SalesZoneType
	AssignedDate  *uint64
	FinishedDate  *uint64
	CreatedAt     *uint64
	UpdatedAt     *uint64
}

type GetListUserByZoneRawQuery struct {
	GetBySalesZoneIdUserIdRawQuery
	UserRoleId *uint
}

type UserIDByZoneData struct {
	UserZoneType        *SalesZoneType
	UserIDByZoneIDSlice *[]string
}

type UserZoneWithEpochEntity struct {
	UserZoneID    *uint   `json:"user_zone_id"`
	UserId        *uint   `json:"user_id"`
	SalesZoneId   *uint   `json:"sales_zone_id"`
	SalesZoneType *string `json:"sales_zone_type"`
	AssignedDate  *uint64 `json:"assigned_date"`
	CreatedAt     *uint64 `json:"created_at"`
	UpdatedAt     *uint64 `json:"updated_at"`
}
