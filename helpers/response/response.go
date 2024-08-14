package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, gin.H{"error": err.Error()})
}

func JSON(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func JSON201(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{"data": data})
}

func JSON204(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, nil)
}
