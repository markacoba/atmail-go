package usecase

import (
	"atmail/backend/model"
	"atmail/backend/user"
	"fmt"
)

type userUC struct {
	userRepo user.Repository
}

// Usecase or Services Layer
func NewUserUsecase(userRepo user.Repository) user.UseCase {
	return &userUC{userRepo}
}

/*
* FetchById
* @param {uint} id
* @returns {*model.User, error}
 */
func (uc *userUC) FetchById(id uint) (*model.User, error) {
	return uc.userRepo.FetchById(id)
}

/*
* Store
* @param {model.User} user
* @returns {*model.User, error}
 */
func (uc *userUC) Store(user model.User) (*model.User, error) {
	// Validate if username already exist
	_, err := uc.userRepo.FetchByUsername(user.UserName)
	if err == nil {
		return nil, fmt.Errorf("username already exist")
	}

	return uc.userRepo.Store(user)
}

/*
* Update
* @param {uint} id
* @returns {*model.User, error}
 */
func (uc *userUC) Update(id uint, user model.User) (*model.User, error) {
	u, err := uc.userRepo.FetchById(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Validate if new username doesn't exist
	if u.UserName != user.UserName {
		_, err := uc.userRepo.FetchByUsername(user.UserName)
		if err == nil {
			return nil, fmt.Errorf("username already exist")
		}
	}

	return uc.userRepo.Update(id, user)
}

/*
* Update
* @param {uint} id
* @returns {error}
 */
func (uc *userUC) Delete(id uint) error {
	return uc.userRepo.Delete(id)
}
