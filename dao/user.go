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

func UserUpdateAccessToken(user *db.User, token uuid.UUID) error {
	user.AccessToken = &token

	return db.DB.
		Model(user).
		Update("access_token", user.AccessToken).
		Error
}

func UserRemoveAccessToken(user *db.User) error {
	user.AccessToken = nil

	return db.DB.
		Model(user).
		Update("access_token", user.AccessToken).
		Error
}

func UserUpdateEmailAndLanguage(user *db.User, email, language string) error {
	err := db.DB.
		Model(user).
		Updates(
			map[string]interface{}{
				"Email":    email,
				"Language": language,
			},
		).
		Error
	if err != nil {
		return err
	}

	user.Email = email
	user.Language = language

	return nil
}
