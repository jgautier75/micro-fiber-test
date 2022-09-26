package api

import (
	"micro-fiber-test/pkg/model"
)

type UserDaoInterface interface {
	Create(cnxParams string, user model.UserInterface) (int64, error)
}
