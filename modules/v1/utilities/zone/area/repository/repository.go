package repository

import (
	"ethical-be/app/config"
	model "ethical-be/modules/v1/utilities/zone/area/models"
	relations "ethical-be/modules/v1/utilities/zone/area/models/relations"

	"gorm.io/gorm"
)

var (
	conf, _ = config.Init()
)

type AreasReporitory interface {
	FindAll() ([]model.Areas, error)
	FindById(ID int) (relations.AreaRelation, error)
	FindById_ByIsVacant(ID int, IsVacant bool) (relations.AreaRelation, error)
	Create(areas model.Areas) (model.Areas, error)
	Update(areas model.Areas) (model.Areas, error)
	Delete(ID int) error
	FindAreaByRegionID(SliceRegionID []int) ([]model.Areas, error)
}

type repositoryareas struct {
	db *gorm.DB
}

func NewAreasRepository(db *gorm.DB) *repositoryareas {
	return &repositoryareas{db}
}

func (r *repositoryareas) FindAll() ([]model.Areas, error) {
	var area []model.Areas
	err := r.db.Find(&area).Error
	return area, err
}

func (r *repositoryareas) FindById(ID int) (relations.AreaRelation, error) {
	var areaRelation relations.AreaRelation
	err := r.db.Preload("GtUnderArea").Where("areas.id = ?", ID).Find(&areaRelation).Error

	return areaRelation, err
}

func (r *repositoryareas) FindAreaByRegionID(SliceRegionID []int) ([]model.Areas, error) {
	var area []model.Areas

	err := r.db.Where("region_id IN ?", SliceRegionID).Order("id asc").Find(&area).Error

	return area, err
}

func (r *repositoryareas) FindById_ByIsVacant(ID int, IsVacant bool) (relations.AreaRelation, error) {
	var Isvacant relations.AreaRelation

	err := r.db.Preload("GtUnderArea").Where("areas.id= ?", ID).Where("areas.is_vacant= ?", IsVacant).Find(&Isvacant).Error
	return Isvacant, err
}

func (r *repositoryareas) Create(areas model.Areas) (model.Areas, error) {
	err := r.db.Create(&areas).Error
	return areas, err
}

func (r *repositoryareas) Update(area model.Areas) (model.Areas, error) {
	err := r.db.Save(&area).Error
	return area, err
}

func (r *repositoryareas) Delete(ID int) error {
	var area model.Areas
	err := r.db.Delete(&area, ID).Error
	return err
}
