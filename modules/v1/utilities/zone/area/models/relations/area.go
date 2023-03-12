package relations

import (
	models "ethical-be/modules/v1/utilities/zone/area/models"
	gt "ethical-be/modules/v1/utilities/zone/gt/models/relations"
)

type AreaRelation struct {
	models.Areas
	GtUnderArea []gt.GroupTerritoriesRelation `gorm:"foreignKey:AreaID;References:ID"`
}

func (AreaRelation) TableName() string {
	return "areas"
}
