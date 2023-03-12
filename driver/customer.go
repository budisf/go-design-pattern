package driver

import (
	customer "ethical-be/modules/v1/utilities/customer/handler"
	repo "ethical-be/modules/v1/utilities/customer/repository"
	service "ethical-be/modules/v1/utilities/customer/services"
)

var (
	CustomerRepository = repo.NewCustomerRepository(DB)
	CustomerService    = service.NewCustomerService(CustomerRepository)
	CustomerHandler    = customer.NewCustomerHandler(CustomerService)
)
