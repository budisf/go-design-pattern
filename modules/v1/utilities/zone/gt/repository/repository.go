package repository

import (
	models "ethical-be/modules/v1/utilities/zone/gt/models"
	relations "ethical-be/modules/v1/utilities/zone/gt/models/relations"

	"gorm.io/gorm"
)

type GtRepository interface {
	FindAll() ([]models.GroupTerritories, error)
	FindById(ID int) (relations.GroupTerritoriesRelation, error)
	FindById_ByIsVacant(ID int, IsVacant bool) (relations.GroupTerritoriesRelation, error)
	Create(gt models.GroupTerritories) (models.GroupTerritories, error)
	Update(gt models.GroupTerritories) (models.GroupTerritories, error)
	Delete(ID int) error
	FindGTByAreaID(SliceAreaID []int) ([]models.GroupTerritories, error)
}

type repositorygt struct {
	db *gorm.DB
}

func NewGtrespository(db *gorm.DB) *repositorygt {
	return &repositorygt{db}
}

func (r *repositorygt) FindAll() ([]models.GroupTerritories, error) {
	var gt []models.GroupTerritories
	err := r.db.Find(&gt).Error
	return gt, err
}

func (r *repositorygt) FindById(ID int) (relations.GroupTerritoriesRelation, error) {
	var gtRelations relations.GroupTerritoriesRelation

	err := r.db.Preload("OutletUnderGt").Where("group_territories.id = ?", ID).Find(&gtRelations).Error

	return gtRelations, err
}

func (r *repositorygt) FindGTByAreaID(SliceAreaID []int) ([]models.GroupTerritories, error) {
	var gt []models.GroupTerritories

	err := r.db.Where("area_id IN ?", SliceAreaID).Order("id asc").Find(&gt).Error

	return gt, err
}

func (r *repositorygt) FindById_ByIsVacant(ID int, IsVacant bool) (relations.GroupTerritoriesRelation, error) {
	var Isvacant relations.GroupTerritoriesRelation

	err := r.db.Preload("OutletUnderGt").Where("group_territories.id = ?", ID).Where("group_territories.is_vacant= ?", IsVacant).Find(&Isvacant).Error
	return Isvacant, err
}

func (r *repositorygt) CountGroupterritories() (int, error) {
	var gt []models.GroupTerritories
	var count int64
	err := r.db.Find(&gt).Count(&count).Error
	countInt := int(count)
	return countInt, err
}

func (r *repositorygt) Create(gt models.GroupTerritories) (models.GroupTerritories, error) {
	err := r.db.Create(&gt).Error
	return gt, err
}

func (r *repositorygt) Update(gt models.GroupTerritories) (models.GroupTerritories, error) {
	err := r.db.Save(&gt).Error
	return gt, err
}

func (r *repositorygt) Delete(ID int) error {
	var gt models.GroupTerritories
	err := r.db.Delete(&gt, ID).Error
	return err
}
