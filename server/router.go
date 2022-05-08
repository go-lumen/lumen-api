package server

import (
	"github.com/go-lumen/lumen-api/config"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"github.com/go-lumen/lumen-api/controllers"
	"github.com/go-lumen/lumen-api/middlewares"

	"github.com/gin-gonic/gin"
)

// Index is the default place
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the " + config.GetString(c, "mail_sender_name") + " API."})
}

// SetupRouter is the main routing point
func (a *API) SetupRouter() (mongoDB *mongo.Database, config *viper.Viper) {
	router := a.Router

	router.Use(middlewares.ErrorMiddleware())

	router.Use(middlewares.CorsMiddleware(middlewares.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.Use(middlewares.ConfigMiddleware(a.Config))

	dbType := a.Config.GetString("db_type")
	switch dbType {
	case "mongo":
		router.Use(middlewares.StoreMongoMiddleware(a.MongoDatabase))
		/*case "postgresql":
		router.Use(middlewares.StorePostgreMiddleware(a.PostgreDatabase))*/
	}

	router.Use(middlewares.EmailMiddleware(a.EmailSender))
	router.Use(middlewares.TextMiddleware(a.TextSender))

	authenticationMiddleware := middlewares.AuthenticationMiddleware()
	//authorizationMiddleware := middlewares.AuthorizationMiddleware()

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			userController := controllers.NewUserController()
			authentication.POST("/", authController.TokensGeneration)
			authentication.POST("/renew", authController.TokenRenewal)
			authentication.Use(authenticationMiddleware)
			authentication.GET("/", userController.GetUserMe)
			//https://skarlso.github.io/2016/06/12/google-signin-with-go/
			//https://github.com/zalando/gin-oauth2/blob/master/google/google.go
		}

		users := v1.Group("/users")
		{
			userController := controllers.NewUserController()
			users.POST("/resetPassword/:email", userController.ResetPasswordRequest)
			users.POST("/update", userController.UpdateUser)
			users.Use(authenticationMiddleware)
			users.POST("/", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.POST("/changeGroup/:groupID", userController.ChangeUserGroup)
			users.DELETE("/:id", userController.DeleteUser)
			users.GET("/", userController.GetUsers)
		}

		organizations := v1.Group("/organizations")
		{
			organizationController := controllers.NewOrganizationController()
			organizations.GET("/token/:id", organizationController.GetOrganizationByAppKey)
			//organizations.GET("/token/:id/trackers", organizationController.ListOrgDevices)
			organizations.Use(authenticationMiddleware)
			organizations.POST("/", organizationController.CreateOrganization)
			organizations.GET("/index", organizationController.GetOrganizationsIndex)
			organizations.GET("/", organizationController.GetOrganizations)
			organizations.GET("/me", organizationController.GetUserOrganization)
			organizations.GET("/:id", organizationController.GetOrganization)
			organizations.GET("/:id/groups", organizationController.GetOrganizationGroups)
		}

		groups := v1.Group("/groups")
		{
			groupController := controllers.NewGroupController()
			groups.Use(authenticationMiddleware)
			groups.POST("/", groupController.CreateGroup)
			groups.GET("/index", groupController.GetGroupsIndex)
			groups.GET("/", groupController.GetGroups)
			groups.GET("/:id", groupController.GetGroup)
			groups.GET("/me", groupController.GetUserGroup)
		}

	}
	return a.MongoDatabase, a.Config
}
