package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/routes"
)

func main() {
	router := gin.Default()

	routes.SetupPublicRoutes(router)
	routes.SetupProtectedRoutes(router)

	// since this /test registered after setup protected routes method, it is protected as well (tested)
	router.GET("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"msg": "success"}) })

	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Server Error: %s\n", err.Error())
	}

}
