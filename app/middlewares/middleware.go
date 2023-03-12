package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ethical-be/app/config"
	driver "ethical-be/driver"
	res "ethical-be/pkg/api-response"
	jwt "ethical-be/pkg/jwt"

	"github.com/gin-gonic/gin"
)

var (
	env, _              = config.Init()
	authorizationHeader = "Authorization"
	apiKeyHeader        = "X-API-key"
	cronExecutedHeader  = "X-Appengine-Cron"
	valName             = "FIREBASE_ID_TOKEN"
)

func AuthJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		authHeader := c.Request.Header.Get(authorizationHeader)
		token := strings.Replace(authHeader, "Bearer ", "", 1)
		// idToken, _ := client.VerifyIDToken(c, token)

		validate_token, err := jwt.ValidateToken(token, env.App.Secret_key)

		if err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			c.JSON(http.StatusUnauthorized, res.UnAuthorized(errorMessage))
			c.Abort()
			return
		}

		uuid, errExtract := jwt.ExtractTokenUUID(validate_token)
		if errExtract != nil {
			errorMessage := fmt.Sprintf("%v", err)
			c.JSON(http.StatusUnauthorized, res.UnAuthorized(errorMessage))
			c.Abort()
			return
		}

		log.Println("auth time", startTime)
		c.Set("uuid", uuid)
		c.Set("token", token)
		c.Set(valName, validate_token)
		c.Next()

	}
}

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser := c.Request.Header.Get("Auth-User-Id")
		userRepository := driver.UserRepository
		intAuthUser, err := strconv.Atoi(authUser)
		if err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			c.JSON(http.StatusUnauthorized, res.UnAuthorized(errorMessage))
			c.Abort()
			return
		}
		_, err, userLogin := userRepository.GetUserByAuthServerId(intAuthUser)
		if err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			c.JSON(http.StatusUnauthorized, res.UnAuthorized(errorMessage))
			c.Abort()
			return
		}
		if userLogin.ID == nil {
			fmt.Println("unauthorized user not found")
			c.JSON(http.StatusUnauthorized, res.UnAuthorized("unauthorized user not found"))
			c.Abort()
			return
		}
		c.Set("user", *userLogin)
		c.Set("user_id", *userLogin.ID)
		c.Set("user_name", *userLogin.Name)
		c.Set("user_nip", *userLogin.Nip)
		c.Set("user_role_id", *userLogin.RoleId)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Auth-User-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, HEAD, PATCH, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
