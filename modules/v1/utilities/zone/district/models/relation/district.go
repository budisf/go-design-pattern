package relation

import (
	models "ethical-be/modules/v1/utilities/zone/district/models"
	region "ethical-be/modules/v1/utilities/zone/region/models"
)

type DistritcRelation struct {
	models.Districts
	RegionsUnderDistrict []region.Regions `gorm:"foreignKey:DistrictID;References:ID"`
}

func (DistritcRelation) TableName() string {
	return "districts"
}
