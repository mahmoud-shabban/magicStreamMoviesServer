package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/controllers"
)

func main() {
	router := gin.Default()

	router.GET("/movies", controllers.GetMovies())
	router.GET("/movie/:id", controllers.GetMovieByID())
	router.POST("/addmovie", controllers.AddMovie())
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Server Error: %s\n", err.Error())
	}
}
