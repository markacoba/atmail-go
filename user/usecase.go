package user

import (
	"atmail/backend/model"
)

type UseCase interface {
	FetchById(id uint) (*model.User, error)
	Store(user model.User) (*model.User, error)
	Update(id uint, user model.User) (*model.User, error)
	Delete(id uint) error
}
