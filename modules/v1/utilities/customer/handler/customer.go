package handler

/*
   |--------------------------------------------------------------------------
   | Import All package
   |--------------------------------------------------------------------------
*/

import (
	"errors"
	model "ethical-be/modules/v1/utilities/customer/models"
	service "ethical-be/modules/v1/utilities/customer/services"
	res "ethical-be/pkg/api-response"
	"fmt"

	helper "ethical-be/pkg/helpers"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

/*
   |--------------------------------------------------------------------------
   | Type CustomerHandler
   |--------------------------------------------------------------------------
   |
   | CustomerHandler is struct to define CustomerService with type ICustomerServer
   | from interface ICustomerService in file service
*/

type CustomerHandler struct {
	customerService service.ICustomerService
}

/*
   |--------------------------------------------------------------------------
   | New Customer Handler
   |--------------------------------------------------------------------------
   |
   | Create new object for Customer handler and call from router
*/

func NewCustomerHandler(customerService service.ICustomerService) *CustomerHandler {
	return &CustomerHandler{customerService}
}

/*
   |--------------------------------------------------------------------------
   | Index
   |--------------------------------------------------------------------------
   |
   | Get all data and validation
*/
func (h *CustomerHandler) Index(c *gin.Context) {

	limitString := c.Query("limit")
	pageString := c.Query("page")

	limit, errLimit := strconv.Atoi(limitString)
	if errLimit != nil {
		c.JSON(http.StatusBadRequest, res.BadRequest("Limit requiered number"))
		return
	}

	page, errPage := strconv.Atoi(pageString)
	if errPage != nil {
		c.JSON(http.StatusBadRequest, res.BadRequest("Page requiered number"))
		return
	}

	/*
	 Get function findAll() Customer from service
	*/
	customer, count, err := h.customerService.FindAll(page, limit)

	/*
	 Error Handler when data not found return JSON
	*/
	if len(customer) == 0 {
		c.JSON(http.StatusOK, res.Success(nil))
		return
	}

	/*
	 Error Handler when internal server error return JSON
	*/
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ServerError(err))
		return
	}

	/*
	 Looping data Customer for convert data Customer from type Customer to type CustomerResponse
	 and then save to variable CustomersResponse
	*/
	var customersResponse []model.CustomerResponse
	for _, b := range customer {
		customerResponse := responseCustomer(b)
		customersResponse = append(customersResponse, customerResponse)
	}

	endpoint := "v1/customer?"
	metadata := helper.PaginationMetadata(count, limit, &page, endpoint)
	result := res.Pagination{
		MetaData: &metadata,
		Records:  customersResponse,
	}

	/*
	 Return Customer result to JSON if success
	*/
	c.JSON(http.StatusOK, res.Success(result))
}

/*
   |--------------------------------------------------------------------------
   | GetById
   |--------------------------------------------------------------------------
   |
   | Get data by id Customer and validation
*/

func (h *CustomerHandler) GetById(c *gin.Context) {

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	CustomerResult, err := h.customerService.FindByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ServerError(err))
		return
	}

	if CustomerResult.ID == 0 {
		c.JSON(http.StatusNotFound, res.NotFound("ID"))
		return
	}

	CustomerResponse := responseCustomer(CustomerResult)

	c.JSON(http.StatusOK, res.Success(CustomerResponse))
}

/*
   |--------------------------------------------------------------------------
   | Create
   |--------------------------------------------------------------------------
   |
   | Create data and validation
*/

func (h *CustomerHandler) Create(c *gin.Context) {
	var CustomerRequest model.CustomerRequest

	err := c.ShouldBindJSON(&CustomerRequest)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			c.JSON(http.StatusBadRequest, res.BadRequest(errorMessages))
			return
		}
		c.JSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}

	CustomerResult, err := h.customerService.Create(CustomerRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ServerError(err.Error()))
		return
	}
	CustomerResponse := responseCustomer(*CustomerResult)
	c.JSON(http.StatusOK, res.Success(CustomerResponse))
}

/*
   |--------------------------------------------------------------------------
   | Edit
   |--------------------------------------------------------------------------
   |
   | Update data and validation
*/

func (h *CustomerHandler) Edit(c *gin.Context) {
	var CustomerRequest model.CustomerRequest

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	err := c.ShouldBindJSON(&CustomerRequest)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			c.JSON(http.StatusBadRequest, res.BadRequest(errorMessages))
			return
		}
		c.JSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}

	CustomerResult, err, statusCode := h.customerService.Update(id, CustomerRequest)
	if err != nil {
		if statusCode == 404 {
			c.JSON(http.StatusNotFound, res.NotFound("ID"))
			return
		}
		c.JSON(http.StatusInternalServerError, res.ServerError(err.Error()))
		return
	}
	CustomerResponse := responseCustomer(*CustomerResult)
	c.JSON(http.StatusOK, res.Success(CustomerResponse))
}

/*
   |--------------------------------------------------------------------------
   | Delete
   |--------------------------------------------------------------------------
   |
   | Soft delete data Customer
*/

func (h *CustomerHandler) Delete(c *gin.Context) {

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	fmt.Println(id)

	err, statusCode := h.customerService.Delete(id)
	if err != nil {
		if statusCode == 404 {
			c.JSON(http.StatusNotFound, res.NotFound("ID"))
			return
		}
		c.JSON(http.StatusInternalServerError, res.ServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, res.StatusOK("Customer deleted success"))
}

/*
   |--------------------------------------------------------------------------
   | Response Customer
   |--------------------------------------------------------------------------
   |
   | This function is for return Customer's response,
   | This function convert type Customer's to type CustomerResponse.
   | You can call this function when you return data from database to client.
*/

func responseCustomer(b model.Customer) model.CustomerResponse {

	customerResponse := model.CustomerResponse{
		ID:           b.ID,
		CustomerCode: b.CustomerCode,
		CustomerName: b.CustomerName,
		Specialist:   b.Specialist,
		CreatedAt:    uint(b.CreatedAt.Unix()),
	}
	if b.UpdatedAt != nil {
		unix := helper.ConvertDateToUnix(*b.UpdatedAt)
		customerResponse.UpdatedAt = uint(unix)
	}
	return customerResponse
}
