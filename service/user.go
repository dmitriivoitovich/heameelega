package service

import (
	"errors"
	"strings"

	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/util/apperror"
	"github.com/dmitriivoitovich/heameelega/util/i18n"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserAuth(token string) (*db.User, *apperror.Error) {
	if strings.TrimSpace(token) == "" {
		return nil, nil
	}

	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return nil, nil
	}

	user, err := dao.UserByAccessToken(tokenUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, apperror.Internal(err, "failed to load user by access token")
	}

	return user, nil
}

func UserLogin(req request.LoginUser) (*db.User, *apperror.Error) {
	user, err := dao.UserByEmail(req.Email)
	if err != nil {
		if dao.IsErrRecordNotFound(err) {
			return nil, apperror.Validation(i18n.KeyErrorCredentialsInvalid)
		}

		return nil, apperror.Internal(err, "failed to load user by email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.Validation(i18n.KeyErrorCredentialsInvalid)
	}

	if err := dao.UserUpdateAccessToken(user, uuid.New()); err != nil {
		return nil, apperror.Internal(err, "failed to update user access token")
	}

	return user, nil
}

func UserLogout(user *db.User) *apperror.Error {
	if err := dao.UserRemoveAccessToken(user); err != nil {
		return apperror.Internal(err, "failed to remove user's access token")
	}

	return nil
}

func UserRegister(req request.RegisterUser) (*db.User, *apperror.Error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal(err, "failed to generate password hash")
	}

	user := &db.User{
		ID:       uuid.New(),
		Email:    strings.ToLower(req.Email),
		Password: string(hash),
		Language: db.LanguageCodeEnglish,
	}

	if err := dao.CreateUser(user); err != nil {
		if dao.IsErrDuplicateKey(err) {
			return nil, apperror.Validation(i18n.KeyErrorEmailTaken)
		}

		return nil, apperror.Internal(err, "failed to create new user")
	}

	if err := dao.UserUpdateAccessToken(user, uuid.New()); err != nil {
		return nil, apperror.Internal(err, "failed to update user access toekn")
	}

	return user, nil
}

func UserUpdate(user *db.User, req request.EditUser) *apperror.Error {
	if err := dao.UserUpdateEmailAndLanguage(user, req.Email, req.Language); err != nil {
		if dao.IsErrDuplicateKey(err) {
			return apperror.Validation(i18n.KeyErrorEmailTaken)
		}

		return apperror.Internal(err, "failed to update user")
	}

	return nil
}
