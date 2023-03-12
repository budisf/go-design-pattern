package repository

import (
	"ethical-be/app/config"
	model "ethical-be/modules/v1/utilities/zone/district/models"
	relations "ethical-be/modules/v1/utilities/zone/district/models/relation"

	"gorm.io/gorm"
)

var (
	conf, _ = config.Init()
)

type RepositoryDistrict interface {
	FindAll() ([]model.Districts, error)
	FindById(ID int) (relations.DistritcRelation, error)
	FindById_ByIsVacant(ID int, IsVacant bool) (relations.DistritcRelation, error)
	Create(district model.Districts) (model.Districts, error)
	Update(district model.Districts) (model.Districts, error)
	Delete(district model.Districts) error
}

type repositorydistrict struct {
	db *gorm.DB
}

func NewDistrictRepository(db *gorm.DB) *repositorydistrict {
	return &repositorydistrict{db}
}

func (r *repositorydistrict) FindAll() ([]model.Districts, error) {
	var district []model.Districts
	err := r.db.Find(&district).Error
	return district, err
}

func (r *repositorydistrict) FindById(ID int) (relations.DistritcRelation, error) {
	var district relations.DistritcRelation
	err := r.db.Preload("RegionsUnderDistrict.AreasUnderRegion.GtUnderArea").Preload("RegionsUnderDistrict.AreasUnderRegion").Preload("RegionsUnderDistrict").Where("districts.id = ?", ID).Find(&district).Error
	return district, err
}

func (r *repositorydistrict) FindById_ByIsVacant(ID int, IsVacant bool) (relations.DistritcRelation, error) {
	var Isvacant relations.DistritcRelation

	err := r.db.Preload("RegionsUnderDistrict").Where("districts.id =?", ID).Where("districts.is_vacant= ?", IsVacant).Find(&Isvacant).Error
	return Isvacant, err
}

func (r *repositorydistrict) Create(district model.Districts) (model.Districts, error) {
	err := r.db.Create(&district).Error
	return district, err
}

func (r *repositorydistrict) Update(district model.Districts) (model.Districts, error) {
	err := r.db.Save(&district).Error
	return district, err
}

func (r *repositorydistrict) Delete(district model.Districts) error {
	err := r.db.Delete(&district).Error
	return err
}
