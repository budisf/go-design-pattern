package services

import (
	models "ethical-be/modules/v1/utilities/zone/district/models"
	repos "ethical-be/modules/v1/utilities/zone/district/repository"
	relation "ethical-be/modules/v1/utilities/zone/district/models/relation"

	"fmt"
)

type RegionsService interface {
	FindAll() ([]models.Districts, error)
	FindById(ID int) (relation.DistritcRelation, error)
	Create(districtRequest models.DistrictRequest) (models.Districts, error)
	Update(ID int, districtRequest models.DistrictRequest) (models.Districts, error)
	Delete(ID int) error
}

type service struct {
	repositorydistrict repos.RepositoryDistrict
}

func NewDistrictService(repositorydistrict repos.RepositoryDistrict) *service {
	return &service{repositorydistrict}
}

func (service *service) FindAll() ([]models.Districts, error) {
	district, err := service.repositorydistrict.FindAll()
	return district, err
}

func (service *service) FindById(ID int) (relation.DistritcRelation, error) {
	district, err := service.repositorydistrict.FindById(ID)
	return district, err
}

func (service *service) Create(regionsRequest models.DistrictRequest) (models.Districts, error) {
	district := models.Districts{
		Name: regionsRequest.Name,
	}
	newDistrict, err := service.repositorydistrict.Create(district)
	fmt.Println(newDistrict)
	return newDistrict, err
}

func (service *service) Update(ID int, districtRequest models.DistrictRequest) (models.Districts, error) {
	district, err := service.repositorydistrict.FindById(ID)
	district.Name = districtRequest.Name

	newDistrict, err := service.repositorydistrict.Update(district.Districts)
	return newDistrict, err
}

func (service *service) Delete(ID int) error {
	district, errby := service.repositorydistrict.FindById(ID)
	if errby != nil {
		return errby
	}
	err := service.repositorydistrict.Delete(district.Districts)
	return err
}
