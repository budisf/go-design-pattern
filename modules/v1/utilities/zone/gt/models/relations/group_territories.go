package relation

import (
	models "ethical-be/modules/v1/utilities/zone/gt/models"
)

type GroupTerritoriesRelation struct {
	models.GroupTerritories
}

func (GroupTerritoriesRelation) TableName() string {
	return "group_territories"
}
