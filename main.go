package main

import (
	// config "ethical-be/app/config"
	// database "ethical-be/app/databases"

	"ethical-be/app/firebase"
	middleware "ethical-be/app/middlewares"
	"log"

	//"ethical-be/modules/v1/routes"
	driver "ethical-be/driver"
	"ethical-be/modules/v1/routes"
	"ethical-be/pkg/helpers"
	errors "ethical-be/pkg/http-errors"
	singleton "ethical-be/pkg/value-object"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.DebugMode)
	singleton.GetPaginationValueObject() //Init Pagination Default ValueObject
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }
	/*
	 Define .env configuratian
	*/
	conf := driver.Conf
	if driver.ErrConf != nil {
		log.Fatal(driver.ErrConf)
	}

	// godotenv.Load()

	//ctx := context.Background()

	/*
	 Database connection
	*/
	db := driver.DB

	/*
	 Firebase Config Init
	*/
	_ = godotenv.Load()
	firebase.FirebaseCredentialInit(&conf)

	validationHelperInstance := helpers.HelperInit(db, &conf)

	/*
	 Database migration, uncomment if you want to migrate database
	*/
	// database.Migrate(db)

	router := gin.Default()
	router.RedirectTrailingSlash = false
	// router.RedirectFixedPath = true
	/*
	 Error Handling when page not found and method not allowed
	*/
	errors.Init(router)

	/*
	 CORSMiddleware
	*/
	router.Use(middleware.CORSMiddleware())

	/*
	 Registration Custom Validation ENUM
	*/
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("Enum", helpers.Enum)
		_ = v.RegisterValidation("EnumVersionTwo", helpers.EnumVersionTwo)
		_ = v.RegisterValidation("UserExists", validationHelperInstance.UserExists)
	}

	/*
	 Routes define
	*/

	routes.Customer(router, &conf, *driver.CustomerHandler)
	routes.Regions(router, &conf, driver.RegionsHandler)
	routes.Areas(router, &conf, driver.AreaHandler)
	routes.GroupTerritories(router, &conf, driver.GtHandler)

	routes.Role(router, &conf, driver.RoleHandler)
	routes.User(router, &conf, driver.UserHandler, driver.UserRoleHandler, driver.UserZoneHandler)

	/*
	 Server connection
	*/
	port := &conf.App.Port
	// port := os.Getenv("APP_PORT")
	router.Run(":" + *port)
}
