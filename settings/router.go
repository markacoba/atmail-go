package settings

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(DefaultSecurity())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	// config.AllowedOrigins = []string{
	// 	os.Getenv("ADMINPANEL_URL"),
	// }
	router.Use(cors.New(config))
	// router.Use(Limit())
	router.NoRoute(noRoute())
	return router
}

func noRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add additional logic here for 404
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}
