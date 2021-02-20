package dao

import (
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/google/uuid"
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

func UserByAccessToken(token uuid.UUID) (*db.User, error) {
	user := &db.User{}

	err := db.DB.
		Where("access_token = ?", token).
		First(user).
		Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UserUpdateAccessToken(userID, token uuid.UUID) error {
	return db.DB.
		Model(&db.User{}).
		Where("id = ?", userID).
		Update("access_token", token).
		Error
}

func UserRemoveAccessToken(userID, token uuid.UUID) error {
	return db.DB.
		Model(&db.User{}).
		Where("id = ?", userID).
		Update("access_token", nil).
		Error
}
