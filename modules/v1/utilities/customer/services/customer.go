package services

import (
	"errors"
	model "ethical-be/modules/v1/utilities/customer/models"
	repo "ethical-be/modules/v1/utilities/customer/repository"
	"time"
)

type ICustomerService interface {
	FindAll(page int, limit int) ([]model.Customer, *int64, error)
	FindByID(id int) (model.Customer, error)
	Create(customerRequest model.CustomerRequest) (*model.Customer, error)
	Update(id int, customerRequest model.CustomerRequest) (*model.Customer, error, int)
	Delete(id int) (error, int)
}

type customerService struct {
	repository repo.ICustomerRepository
}

func NewCustomerService(repository repo.ICustomerRepository) *customerService {
	return &customerService{repository}
}

func (s *customerService) FindAll(page int, limit int) ([]model.Customer, *int64, error) {

	pageA := page - 1
	if page == 0 {
		pageA = 0
	}
	skip := pageA * limit
	customer, count, err := s.repository.FindAll(limit, skip)

	return customer, count, err
}

func (s *customerService) FindByID(id int) (model.Customer, error) {
	customer, err := s.repository.FindByID(id)
	return customer, err
}

func (s *customerService) Create(CustomerRequest model.CustomerRequest) (*model.Customer, error) {

	Customer := model.Customer{
		CustomerName: CustomerRequest.CustomerName,
		Specialist:   CustomerRequest.Specialist,
		CustomerCode: CustomerRequest.CustomerCode,
	}

	newCustomer, err := s.repository.Create(Customer)
	return &newCustomer, err
}

func (s *customerService) Update(id int, CustomerRequest model.CustomerRequest) (*model.Customer, error, int) {

	Customer, errGet := s.repository.FindByID(id)
	if errGet != nil {
		return nil, errGet, 500
	}
	if Customer.ID == 0 {
		return nil, errors.New("ID Not Found"), 404
	}
	updatedDate := time.Now()

	Customer.CustomerName = CustomerRequest.CustomerName
	Customer.Specialist = CustomerRequest.Specialist
	Customer.CustomerCode = CustomerRequest.CustomerCode
	Customer.UpdatedAt = &updatedDate

	newCustomer, err := s.repository.Update(Customer)
	if err != nil {
		return nil, nil, 500
	}
	return &newCustomer, nil, 200
}

func (s *customerService) Delete(id int) (error, int) {

	Customer, errGet := s.repository.FindByID(id)
	if errGet != nil {
		return errGet, 500
	}
	if Customer.ID == 0 {
		return nil, 404
	}
	deletedDate := time.Now()

	Customer.IsDeleted = true
	Customer.DeletedAt = &deletedDate

	_, err := s.repository.Delete(Customer)
	if err != nil {
		return err, 500
	}
	return nil, 200
}
