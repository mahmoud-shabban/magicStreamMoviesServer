package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/controllers"
)

func SetupPublicRoutes(router *gin.Engine) {

	router.GET("/movies", controllers.GetMovies())
	router.POST("/register", controllers.RegisterUser())
	router.POST("/login", controllers.LoginUser())

}
