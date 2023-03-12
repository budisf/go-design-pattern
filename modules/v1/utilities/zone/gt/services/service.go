package services

import (
	models "ethical-be/modules/v1/utilities/zone/gt/models"
	relation "ethical-be/modules/v1/utilities/zone/gt/models/relations"
	repos "ethical-be/modules/v1/utilities/zone/gt/repository"
)

type GtServices interface {
	FindAll() ([]models.GroupTerritories, error)
	FindById(ID int) (relation.GroupTerritoriesRelation, error)
	Create(gtsRequest models.AllRequest) (models.GroupTerritories, error)
	Update(ID int, gtRequest models.GtRequest) (models.GroupTerritories, error)
	UpdateAreas(ID int, areaRequest models.AreaRequest) (models.GroupTerritories, error)
	Delete(ID int) error
}

type service struct {
	repositorygt repos.GtRepository
}

func NewGtService(repositorygt repos.GtRepository) *service {
	return &service{repositorygt}
}

func (service *service) FindAll() ([]models.GroupTerritories, error) {
	booking, err := service.repositorygt.FindAll()
	return booking, err
}

func (service *service) FindById(ID int) (relation.GroupTerritoriesRelation, error) {
	areas, err := service.repositorygt.FindById(ID)
	return areas, err
}

func (service *service) Create(gtRequest models.AllRequest) (models.GroupTerritories, error) {
	gt := models.GroupTerritories{
		Name:   gtRequest.Name,
		AreaID: gtRequest.AreaID,
	}
	newGt, err := service.repositorygt.Create(gt)
	return newGt, err
}

func (service *service) Update(ID int, gtRequest models.GtRequest) (models.GroupTerritories, error) {
	gt, _ := service.repositorygt.FindById(ID)

	gt.Name = gtRequest.Name

	newGt, err := service.repositorygt.Update(gt.GroupTerritories)
	return newGt, err
}

func (service *service) UpdateAreas(ID int, areaRequest models.AreaRequest) (models.GroupTerritories, error) {
	area, _ := service.repositorygt.FindById(ID)

	area.AreaID = areaRequest.AreaID

	newGt, err := service.repositorygt.Update(area.GroupTerritories)
	return newGt, err
}

func (service *service) Delete(ID int) error {
	err := service.repositorygt.Delete(ID)

	return err
}
