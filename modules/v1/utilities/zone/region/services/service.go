package services

import (
	models "ethical-be/modules/v1/utilities/zone/region/models"
	relation "ethical-be/modules/v1/utilities/zone/region/models/relations"
	repos "ethical-be/modules/v1/utilities/zone/region/repository"

	"fmt"
)

type RegionsService interface {
	FindAll() ([]models.Regions, error)
	FindById(ID int) (relation.RegionRelation, error)
	Create(regionsRequest models.RegionsRequest) (models.Regions, error)
	Update(ID int, regionsRequest models.RegionsRequest) (models.Regions, error)
	Delete(ID int) error
}

type service struct {
	repositoryregions repos.RegionsRepository
}

func NewRegionsService(repositoryregions repos.RegionsRepository) *service {
	return &service{repositoryregions}
}

func (service *service) FindAll() ([]models.Regions, error) {
	booking, err := service.repositoryregions.FindAll()
	return booking, err
}

func (service *service) FindById(ID int) (relation.RegionRelation, error) {
	regions, err := service.repositoryregions.FindById(ID)
	return regions, err
}

func (service *service) Create(regionsRequest models.RegionsRequest) (models.Regions, error) {
	regions := models.Regions{
		Name:       regionsRequest.Name,
		DistrictID: regionsRequest.DistrictID,
	}
	newRegions, err := service.repositoryregions.Create(regions)
	fmt.Println(newRegions)
	return newRegions, err
}

func (service *service) Update(ID int, regionRequest models.RegionsRequest) (models.Regions, error) {
	regions, err := service.repositoryregions.FindById(ID)
	regions.Name = regionRequest.Name
	regions.DistrictID = regionRequest.DistrictID

	newRegion, err := service.repositoryregions.Update(regions.Regions)
	return newRegion, err
}

func (service *service) Delete(ID int) error {
	regions, errby := service.repositoryregions.FindById(ID)
	if errby != nil {
		return errby
	}
	err := service.repositoryregions.Delete(regions.Regions)
	return err
}
