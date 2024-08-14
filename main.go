package main

import (
	"atmail/backend/model"
	"atmail/backend/settings"
	"log"
	"os"

	_userHttp "atmail/backend/user/http"
	_userRepo "atmail/backend/user/repository"
	_userUC "atmail/backend/user/usecase"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Connect to MySQL DB
	dbString := os.Getenv("DB_STRING")
	if len(dbString) == 0 {
		dbString = "user:pass@tcp(127.0.0.1:3306)/atmail?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dbString), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}

	// Migrate DB Schema
	db.AutoMigrate(
		&model.User{},
	)

	// Setup Router
	router := settings.InitRouter()

	// Routes
	userRepo := _userRepo.NewUserRepo(db)
	userUC := _userUC.NewUserUsecase(userRepo)
	_userHttp.NewUserHandler(router, userUC)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = ":8000"
	} else {
		PORT = ":" + PORT
	}
	router.Run(PORT)
}
