package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	engine.GET("/health", Health())

	return engine
}

func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}
}
