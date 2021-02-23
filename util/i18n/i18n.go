package i18n

import (
	"errors"
	"fmt"
)

const (
	LanguageEnglish = "EN"
	LanguageRussian = "RU"

	KeyEmailTaken         Key = "email-taken"
	KeyCredentialsInvalid Key = "credentials-invalid"
	KeyDataCollision      Key = "data-collision"
)

type Key string

var (
	translations = map[Key]map[string]string{
		KeyEmailTaken: {
			LanguageEnglish: "Email taken",
			LanguageRussian: "Email занят",
		},
	}

	errKeyNotDefined  = errors.New("translation key not defined")
	errLangNotDefined = errors.New("language not defined for key")
)

func Translate(key Key, language string) (string, error) {
	m, ok := translations[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", errKeyNotDefined, key)
	}

	res, ok := m[language]
	if !ok {
		return "", fmt.Errorf("%w: %s %s", errLangNotDefined, language, key)
	}

	return res, nil
}
