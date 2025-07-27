package main

import (
	router "mercor/internal/domain/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.InitRoutes(r)
	r.Run(":8080")
}
