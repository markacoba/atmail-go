package http

import (
	"atmail/backend/helpers/response"
	"atmail/backend/middleware"
	"atmail/backend/model"
	userAdmin "atmail/backend/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userUC userAdmin.UseCase
}

// Controller Layer
func NewUserHandler(router *gin.Engine, userUC userAdmin.UseCase) {
	handler := &userHandler{userUC}

	route := router.Group("/users")
	route.Use(middleware.BasicAuthMiddleware())
	route.Use(middleware.RBACMiddleware())

	// User Routes
	route.GET("/:id", handler.FetchById)
	route.POST("/", handler.Store)
	route.PUT("/:id", handler.Update)
	route.DELETE("/:id", handler.Delete)
}

func (handler *userHandler) FetchById(ctx *gin.Context) {
	id := ctx.Param("id")
	i, _ := strconv.Atoi(id)

	user, err := handler.userUC.FetchById(uint(i))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err)
		return
	}

	response.JSON(ctx, user)
}

func (handler *userHandler) Store(ctx *gin.Context) {
	var user model.User
	err := ctx.BindJSON(&user)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	user = user.CleanData()
	err = user.ValidateUser()
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	u, err := handler.userUC.Store(user)
	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			response.Error(ctx, http.StatusConflict, err)
		} else {
			response.Error(ctx, http.StatusBadRequest, err)
		}

		return
	}

	response.JSON201(ctx, u)
}

func (handler *userHandler) Update(ctx *gin.Context) {
	roleId := ctx.Params.ByName("id")
	id, _ := strconv.Atoi(roleId)

	var user model.User
	err := ctx.BindJSON(&user)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	user = user.CleanData()
	err = user.ValidateUser()
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	if id > 0 {
		u, err := handler.userUC.Update(uint(id), user)
		if err != nil {
			if strings.Contains(err.Error(), "already exist") {
				response.Error(ctx, http.StatusConflict, err)
			} else {
				response.Error(ctx, http.StatusBadRequest, err)
			}
			return
		}
		response.JSON(ctx, u)
	} else {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}
}

func (handler *userHandler) Delete(ctx *gin.Context) {
	eventID := ctx.Params.ByName("id")
	id, _ := strconv.Atoi(eventID)

	_, err := handler.userUC.FetchById(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err)
		return
	}

	if err := handler.userUC.Delete(uint(id)); err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	response.JSON204(ctx)
}
