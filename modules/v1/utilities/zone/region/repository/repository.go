package repository

import (
	"ethical-be/app/config"
	model "ethical-be/modules/v1/utilities/zone/region/models"
	relations "ethical-be/modules/v1/utilities/zone/region/models/relations"

	"gorm.io/gorm"
)

var (
	conf, _ = config.Init()
)

type RegionsRepository interface {
	FindAll() ([]model.Regions, error)
	FindById(ID int) (relations.RegionRelation, error)
	FindById_ByIsVacant(ID int, IsVacant bool) (relations.RegionRelation, error)
	Create(regions model.Regions) (model.Regions, error)
	Update(regions model.Regions) (model.Regions, error)
	Delete(regions model.Regions) error
	FindRegionByDistrictID(SliceDistrictID []int) ([]model.Regions, error)
}

type repositoryregions struct {
	db *gorm.DB
}

func NewRegionsRepository(db *gorm.DB) *repositoryregions {
	return &repositoryregions{db}
}

func (r *repositoryregions) FindAll() ([]model.Regions, error) {
	var region []model.Regions
	err := r.db.Find(&region).Error
	return region, err
}

func (r *repositoryregions) FindById(ID int) (relations.RegionRelation, error) {

	var regionRelation relations.RegionRelation

	err := r.db.Preload("AreasUnderRegion.GtUnderArea").Preload("AreasUnderRegion").Where("regions.id = ?", ID).Find(&regionRelation).Error

	return regionRelation, err
}

func (r *repositoryregions) FindById_ByIsVacant(ID int, IsVacant bool) (relations.RegionRelation, error) {
	var Isvacant relations.RegionRelation

	err := r.db.Preload("AreasUnderRegion").Where("regions.id = ?", ID).Where("regions.is_vacant= ?", IsVacant).Find(&Isvacant).Error
	return Isvacant, err
}

func (r *repositoryregions) FindRegionByDistrictID(SliceDistrictID []int) ([]model.Regions, error) {
	var region []model.Regions

	err := r.db.Where("district_id IN ?", SliceDistrictID).Order("id asc").Find(&region).Error

	return region, err
}

func (r *repositoryregions) Create(regions model.Regions) (model.Regions, error) {
	err := r.db.Create(&regions).Error
	return regions, err
}

func (r *repositoryregions) Update(regions model.Regions) (model.Regions, error) {
	err := r.db.Save(&regions).Error
	return regions, err
}

func (r *repositoryregions) Delete(regions model.Regions) error {
	err := r.db.Delete(&regions).Error
	return err
}
