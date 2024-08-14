package http

import (
	"atmail/backend/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_userRepo "atmail/backend/user/repository"
	_userUC "atmail/backend/user/usecase"
)

type ResponseUser struct {
	Data model.User `json:"data"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func fetchUser(r *gin.Engine, handler userHandler) *gin.Engine {
	r.GET("/users/:id", handler.FetchById)
	return r
}

func postUser(r *gin.Engine, handler userHandler) *gin.Engine {
	r.POST("/users", handler.Store)
	return r
}

func putUser(r *gin.Engine, handler userHandler) *gin.Engine {
	r.PUT("/users/:id", handler.Update)
	return r
}

func deleteUser(r *gin.Engine, handler userHandler) *gin.Engine {
	r.DELETE("/users/:id", handler.Delete)
	return r
}

func Test_userHandler_CRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// TODO: Use a Test DB
	dbString := "atmail_user:TWAKQ6meOtFeir7z@tcp(127.0.0.1:3306)/atmail?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbString), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}

	router := setupRouter()
	userRepo := _userRepo.NewUserRepo(db)
	userUC := _userUC.NewUserUsecase(userRepo)
	handler := &userHandler{userUC}
	router = postUser(router, *handler)

	mockUser := model.User{
		UserName: "user_test",
		Email:    "mark@test.com",
		Age:      30,
		Role:     "admin",
	}

	jsonData, _ := json.Marshal(mockUser)

	// Create User Test
	fmt.Println("Create User Test")
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(jsonData)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	readData, _ := io.ReadAll(w.Body)
	assert.Equal(t, http.StatusCreated, w.Code)

	var created ResponseUser
	if err := json.Unmarshal([]byte(readData), &created); err != nil {
		t.Errorf("err (%v)", err)
	}

	assert.Equal(t, mockUser.UserName, created.Data.UserName)

	id := created.Data.ID

	// Fetch User By Id Test
	fmt.Println("Fetch User By Id Test")
	router = fetchUser(router, *handler)
	req, _ = http.NewRequest("GET", fmt.Sprintf("/users/%d", id), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	readData, _ = io.ReadAll(w.Body)

	if err := json.Unmarshal([]byte(readData), &created); err != nil {
		t.Errorf("err (%v)", err)
	}

	assert.Equal(t, mockUser.UserName, created.Data.UserName)

	// Update User Test
	fmt.Println("Update User By Id Test")
	router = putUser(router, *handler)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/users/%d", id), strings.NewReader(string(jsonData)))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	readData, _ = io.ReadAll(w.Body)

	if err := json.Unmarshal([]byte(readData), &created); err != nil {
		t.Errorf("err (%v)", err)
	}

	assert.Equal(t, mockUser.UserName, created.Data.UserName)

	// Delete User Test
	// TODO: Hard delete the entry
	fmt.Println("Delete User By Id Test")
	router = deleteUser(router, *handler)
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/users/%d", id), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
