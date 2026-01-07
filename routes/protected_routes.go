package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/controllers"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/middlewares"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middlewares.AuthMiddleware())

	router.GET("/movie/:id", controllers.GetMovieByID())
	router.POST("/addmovie", controllers.AddMovie())
}
