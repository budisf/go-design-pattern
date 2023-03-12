package relations

import (
	area "ethical-be/modules/v1/utilities/zone/area/models"
	models "ethical-be/modules/v1/utilities/zone/region/models"
)

type RegionRelation struct {
	models.Regions
	AreasUnderRegion []area.Areas `gorm:"foreignKey:RegionID;References:ID"`
}

func (RegionRelation) TableName() string {
	return "regions"
}
