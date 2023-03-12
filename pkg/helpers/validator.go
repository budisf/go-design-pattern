package helpers

import (
	"ethical-be/app/config"
	roleRepository "ethical-be/modules/v1/utilities/role/repository"
	"ethical-be/modules/v1/utilities/user/repository"
	userService "ethical-be/modules/v1/utilities/user/service"
	helperDatabases "ethical-be/pkg/helpers/databases"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type IValidationHelper interface {
	UserExists(fl validator.FieldLevel) bool
}

type helperInjection struct {
	userService.IUserService
}

func InitHelperInjection(userService userService.IUserService) *helperInjection {
	return &helperInjection{
		userService,
	}
}

type validationHelper struct {
	IValidationHelper
}

func InitValidationHelper(iValidationHelper IValidationHelper) *validationHelper {
	return &validationHelper{
		iValidationHelper,
	}
}

func HelperInit(db *gorm.DB, conf *config.Conf) *validationHelper {
	var (
		helperDatabase       = helperDatabases.InitHelperDatabase(db)
		userRepository       = repository.InitUserRepository(db, helperDatabase, conf)
		roleRepository       = roleRepository.InitRoleRepository(db, helperDatabase, conf)
		userService          = userService.InitUserRepository(userRepository, roleRepository)
		helperInjectionObj   = InitHelperInjection(userService)
		initValidationHelper = InitValidationHelper(helperInjectionObj)
	)

	return initValidationHelper
}

/*
   |--------------------------------------------------------------------------
   | Return Error Message For Validation
   |--------------------------------------------------------------------------
   |
   | This function is for return error message for validation,
   | Its can return array string depends on the errors.
   | This function need go-playground-validator package.
*/

func ErrorMessage(err interface{}) []string {
	errorMessages := []string{}
	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e.ActualTag())
		switch e.ActualTag() {
		case "Enum":
			replacer := *strings.NewReplacer("_", ",")
			errorMessage := fmt.Sprintf("Error on field %s, must be one of: %s", e.Field(), replacer.Replace(e.Param()))
			errorMessages = append(errorMessages, errorMessage)
		case "EnumVersionTwo":
			replacer := *strings.NewReplacer("&", ", ")
			errorMessage := e.Field() + " must be one of " + replacer.Replace(e.Param())
			errorMessages = append(errorMessages, errorMessage)
		case "UserExists":
			errorMessage := fmt.Sprintf("Error on field %s, condition: User with ID %v is not exists", e.Field(), e.Value())
			errorMessages = append(errorMessages, errorMessage)
		case "min":
			errorMessage := fmt.Sprintf("Error on field %s, condition: Should Be At Least %v Character", e.Field(), e.Param())
			errorMessages = append(errorMessages, errorMessage)
		case "e164":
			errorMessage := fmt.Sprintf("Error on field %s, condition: Must Use Country Code Like: +62", e.Field())
			errorMessages = append(errorMessages, errorMessage)
		case "email":
			errorMessage := fmt.Sprintf("Error on field %s, condition: Must Use The Correct Email Format", e.Field())
			errorMessages = append(errorMessages, errorMessage)
		case "gte":
			errorMessage := fmt.Sprintf("Error on field %s, condition: Must Grather Than Equals %v", e.Field(), e.Param())
			errorMessages = append(errorMessages, errorMessage)
		default:
			errorMessage := fmt.Sprintf("Error on field %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}
	}
	return errorMessages
}

// /*
//    |--------------------------------------------------------------------------
//    | Custom Enum Validation without underscore requirement
//    |--------------------------------------------------------------------------
//    |
//    | This function Make Custom Binding Validation For Enum data Type and read the whitespace using delimiter '&'
// */

func (helper *helperInjection) UserExists(
	fl validator.FieldLevel,
) bool {
	value := fl.Field().Int()
	valueString := strconv.Itoa(int(value))
	_, _, result := helper.GetById(&valueString)
	if result != nil {
		return true
	} else {
		return false
	}
}

/*
   |--------------------------------------------------------------------------
   | Custom Enum Validation without underscore requirement
   |--------------------------------------------------------------------------
   |
   | This function Make Custom Binding Validation For Enum data Type and read the whitespace using delimiter '&'
*/

func EnumVersionTwo(
	fl validator.FieldLevel,
) bool {
	enumString := fl.Param()     // get string male_female
	value := fl.Field().String() // the actual field
	fmt.Println(fl.Field())
	enumSlice := strings.Split(enumString, "&") // convert to slice
	fmt.Println(enumSlice)
	for _, v := range enumSlice {
		if value == v {
			return true
		}
	}
	return false
}

/*
   |--------------------------------------------------------------------------
   | Custom Enum Validation
   |--------------------------------------------------------------------------
   |
   | This function Make Custom Binding Validation For Enum data Type
*/

func Enum(
	fl validator.FieldLevel,
) bool {
	enumString := fl.Param()                    // get string male_female
	value := fl.Field().String()                // the actual field
	enumSlice := strings.Split(enumString, "_") // convert to slice
	for _, v := range enumSlice {
		if value == v {
			return true
		}
	}
	return false
}
