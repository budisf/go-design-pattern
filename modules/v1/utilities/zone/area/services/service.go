package services

import (
	models "ethical-be/modules/v1/utilities/zone/area/models"
	relation "ethical-be/modules/v1/utilities/zone/area/models/relations"
	repos "ethical-be/modules/v1/utilities/zone/area/repository"
	// reg "ethical-be/modules/v1/utilities/zone/region/services"
)

type AreasService interface {
	FindAll() ([]models.Areas, error)
	FindById(ID int) (relation.AreaRelation, error)
	Create(areasRequest models.AllRequest) (models.Areas, error)
	Update(ID int, areasRequest models.AreasRequest) (models.Areas, error)
	UpdateRegion(ID int, regRequest models.RegionRequest) (models.Areas, error)
	Delete(ID int) error
}

type service struct {
	respositoryareas repos.AreasReporitory
}

func NewAreasService(repositoryareas repos.AreasReporitory) *service {
	return &service{repositoryareas}
}

func (service *service) FindAll() ([]models.Areas, error) {
	booking, err := service.respositoryareas.FindAll()
	return booking, err
}

func (service *service) FindById(ID int) (relation.AreaRelation, error) {
	areas, err := service.respositoryareas.FindById(ID)
	return areas, err
}

func (service *service) Create(areasRequest models.AllRequest) (models.Areas, error) {
	areas := models.Areas{
		Name:     areasRequest.Name,
		RegionID: areasRequest.RegionID,
	}
	newAreas, err := service.respositoryareas.Create(areas)
	return newAreas, err
}

func (service *service) Update(ID int, areaRequest models.AreasRequest) (models.Areas, error) {
	areas, err := service.respositoryareas.FindById(ID)

	areas.Name = areaRequest.Name

	newArea, err := service.respositoryareas.Update(areas.Areas)
	return newArea, err
}

func (service *service) UpdateRegion(ID int, regRequest models.RegionRequest) (models.Areas, error) {
	reg, _ := service.respositoryareas.FindById(ID)
	reg.RegionID = regRequest.RegionID

	newArea, err := service.respositoryareas.Update(reg.Areas)
	return newArea, err
}

func (service *service) Delete(ID int) error {
	err := service.respositoryareas.Delete(ID)

	return err
}
