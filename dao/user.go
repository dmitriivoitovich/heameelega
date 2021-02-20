package dao

import (
	"github.com/dmitriivoitovich/heameelega/dao/db"
)

func CreateUser(user *db.User) error {
	return db.DB.Create(user).Error
}

func UserByEmail(email string) (*db.User, error) {
	user := &db.User{}

	err := db.DB.
		Where("email = ?", email).
		First(user).
		Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
