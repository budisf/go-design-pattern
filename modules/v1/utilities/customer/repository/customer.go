package repository

import (
	model "ethical-be/modules/v1/utilities/customer/models"

	"gorm.io/gorm"
)

type ICustomerRepository interface {
	FindAll(limit int, skip int) ([]model.Customer, *int64, error)
	FindByID(ID int) (model.Customer, error)
	Create(customer model.Customer) (model.Customer, error)
	Update(customer model.Customer) (model.Customer, error)
	Delete(customer model.Customer) (model.Customer, error)
}

type repository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll(limit int, skip int) ([]model.Customer, *int64, error) {
	var customer []model.Customer
	var count int64

	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, nil, err
	}

	if err := tx.Model(&customer).Where("is_deleted", false).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := r.db.Where("is_deleted", false).Offset(skip).Limit(limit).Order(" id asc ").Find(&customer).Error; err != nil {
		return nil, nil, err
	}
	return customer, &count, tx.Commit().Error
}

func (r *repository) FindByID(ID int) (model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("is_deleted", false).Find(&customer, ID).Error
	return customer, err
}

func (r *repository) Create(customer model.Customer) (model.Customer, error) {
	err := r.db.Create(&customer).Error
	return customer, err
}

func (r *repository) Update(customer model.Customer) (model.Customer, error) {
	err := r.db.Save(&customer).Error
	return customer, err
}

func (r *repository) Delete(customer model.Customer) (model.Customer, error) {
	err := r.db.Save(&customer).Error
	return customer, err
}
