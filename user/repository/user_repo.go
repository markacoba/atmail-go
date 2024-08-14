package repository

import (
	"atmail/backend/model"
	"atmail/backend/user"
	"fmt"

	"gorm.io/gorm"
)

type userRepo struct {
	DB *gorm.DB
}

// Repository or Data Access Layer
func NewUserRepo(DB *gorm.DB) user.Repository {
	return &userRepo{DB}
}

/*
* FetchById
* @param {uint} id
* @returns {*model.User, error}
 */
func (dbr *userRepo) FetchById(id uint) (*model.User, error) {
	var user model.User
	if err := dbr.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

/*
* FetchByUsername
* @param {string} username
* @returns {*model.User, error}
 */
func (dbr *userRepo) FetchByUsername(username string) (*model.User, error) {
	var user model.User
	if err := dbr.DB.First(&user, "user_name = ?", username).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

/*
* Store
* @param {model.User} user
* @returns {*model.User, error}
 */
func (dbr *userRepo) Store(user model.User) (*model.User, error) {
	if err := dbr.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

/*
* Update
* @param {uint} id
* @returns {*model.User, error}
 */
func (dbr *userRepo) Update(id uint, user model.User) (*model.User, error) {
	userSave, _ := dbr.FetchById(id)
	userSave.Age = user.Age
	userSave.Email = user.Email
	userSave.Role = user.Role
	// There is no instructions to validate username so we'll update username every request
	userSave.UserName = user.UserName
	if err := dbr.DB.Save(&userSave).Error; err != nil {
		return nil, err
	}

	return userSave, nil
}

/*
* Update
* @param {uint} id
* @returns {error}
 */
func (dbr *userRepo) Delete(id uint) error {
	db := dbr.DB

	var user model.User
	if err := db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
