package middleware

import (
	"atmail/backend/helpers/response"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Should use a better auth
		user, pass, err := c.Request.BasicAuth()
		if !err || (os.Getenv("USERNAME") != user || os.Getenv("PASSWORD") != pass) {
			response.Error(c, http.StatusUnauthorized, fmt.Errorf("username or password is incorrect"))
			c.Abort()
			return
		}

		c.Next()
	}
}
