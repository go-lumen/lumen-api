package server

import (
	"github.com/go-lumen/lumen-api/config"
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
func (a *API) SetupRouter() {
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

	dbType := a.Config.GetString("db_type")
	switch dbType {
	case "mongo":
		router.Use(middlewares.StoreMongoMiddleware(a.MongoDatabase))
	case "postgresql":
		router.Use(middlewares.StorePostgreMiddleware(a.PostgreDatabase))
	case "mysql":
		router.Use(middlewares.StoreMySQLMiddleware(a.MySQLDatabase))
	}

	router.Use(middlewares.ConfigMiddleware(a.Config))

	router.Use(middlewares.EmailMiddleware(a.EmailSender))
	router.Use(middlewares.TextMiddleware(a.TextSender))

	authMiddleware := middlewares.AuthMiddleware() //User
	adminMiddleware := middlewares.AdminMiddleware()

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			userController := controllers.NewUserController()
			authentication.POST("/", authController.TokensGeneration)
			authentication.POST("/renew", authController.TokenRenewal)
			authentication.Use(authMiddleware)
			authentication.GET("/", userController.GetUserMe)
		}

		users := v1.Group("/users")
		{
			userController := controllers.NewUserController()
			users.POST("/resetPassword/:email", userController.ResetPasswordRequest)
			users.POST("/update", userController.UpdateUser)
			users.POST("/", userController.CreateUser)
			users.Use(authMiddleware)
			users.GET("/:id", userController.GetUser)
			users.Use(adminMiddleware)
			users.DELETE("/:id", userController.DeleteUser)
			users.GET("/", userController.GetUsers)
		}
	}
}
