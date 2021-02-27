package i18n

import (
	"errors"
	"fmt"
	"strings"
)

const (
	LanguageEnglish = "EN"
	LanguageRussian = "RU"

	keyTitleLogin          Key = "title.login"
	keyTitleRegister       Key = "title.register"
	keyTitleSettings       Key = "title.settings"
	keyTitleDashboard      Key = "title.dashboard"
	keyTitleCustomers      Key = "title.customers"
	keyTitleCreateCustomer Key = "title.create-customer"

	keyButtonLogin          Key = "button.login"
	keyButtonRegister       Key = "button.register"
	keyButtonLogout         Key = "button.logout"
	keyButtonUserSettings   Key = "button.user-settings"
	keyButtonAddCustomer    Key = "button.add-customer"
	keyButtonCreateCustomer Key = "button.create-customer"
	keyButtonEditCustomer   Key = "button.edit-customer"
	keyButtonSaveCustomer   Key = "button.save-customer"
	keyButtonSaveUser       Key = "button.save-user"
	keyButtonCancel         Key = "button.cancel"

	keyFieldEmail        Key = "field.email"
	keyFieldPassword     Key = "field.password"
	keyFieldSearch       Key = "field.search"
	keyFieldFirstName    Key = "field.first-name"
	keyFieldLastName     Key = "field.last-name"
	keyFieldBirthDate    Key = "field.birth-date"
	keyFieldGender       Key = "field.gender"
	keyFieldGenderMale   Key = "field.gender-male"
	keyFieldGenderFemale Key = "field.gender-female"
	keyFieldAddress      Key = "field.address"
	keyFieldLanguage     Key = "field.language"
	keyFieldLanguageEn   Key = "field.language-en"
	keyFieldLanguageRu   Key = "field.language-ru"

	keyLinkLogin    Key = "link.login"
	keyLinkRegister Key = "link.register"

	keyColumnName      Key = "column.name"
	keyColumnEmail     Key = "column.email"
	keyColumnBirthDate Key = "column.birth-date"
	keyColumnGender    Key = "column.gender"
	keyColumnAddress   Key = "column.address"
	keyColumnJanuary   Key = "column.month-1"
	keyColumnFebruary  Key = "column.month-2"
	keyColumnMarch     Key = "column.month-3"
	keyColumnApril     Key = "column.month-4"
	keyColumnMay       Key = "column.month-5"
	keyColumnJune      Key = "column.month-6"
	keyColumnJuly      Key = "column.month-7"
	keyColumnAugust    Key = "column.month-8"
	keyColumnSeptember Key = "column.month-9"
	keyColumnOctober   Key = "column.month-10"
	keyColumnNovember  Key = "column.month-11"
	keyColumnDecember  Key = "column.month-12"

	KeyErrorEmailTaken         Key = "error.email-taken"
	KeyErrorCredentialsInvalid Key = "error.credentials-invalid"
	KeyErrorDataCollision      Key = "error.data-collision"

	keyNoSearchResults Key = "message.no-search-results"
)

type Key string

var (
	translations = map[Key]map[string]string{
		// page titles
		keyTitleLogin: {
			LanguageEnglish: "Sign in",
			LanguageRussian: "Вход",
		},
		keyTitleRegister: {
			LanguageEnglish: "Create account",
			LanguageRussian: "Создание аккаунта",
		},
		keyTitleDashboard: {
			LanguageEnglish: "Dashboard",
			LanguageRussian: "Главная",
		},
		keyTitleSettings: {
			LanguageEnglish: "Settings",
			LanguageRussian: "Настройки",
		},
		keyTitleCustomers: {
			LanguageEnglish: "Customers",
			LanguageRussian: "Клиенты",
		},
		keyTitleCreateCustomer: {
			LanguageEnglish: "New customer",
			LanguageRussian: "Новый клиент",
		},

		// buttons
		keyButtonLogin: {
			LanguageEnglish: "Sign in",
			LanguageRussian: "Войти",
		},
		keyButtonRegister: {
			LanguageEnglish: "Create account",
			LanguageRussian: "Создать аккаунт",
		},
		keyButtonLogout: {
			LanguageEnglish: "Sign out",
			LanguageRussian: "Выйти",
		},
		keyButtonUserSettings: {
			LanguageEnglish: "Settings",
			LanguageRussian: "Настройки",
		},
		keyButtonAddCustomer: {
			LanguageEnglish: "Add",
			LanguageRussian: "Создать",
		},
		keyButtonCreateCustomer: {
			LanguageEnglish: "Create",
			LanguageRussian: "Создать",
		},
		keyButtonEditCustomer: {
			LanguageEnglish: "Edit",
			LanguageRussian: "Редактировать",
		},
		keyButtonSaveCustomer: {
			LanguageEnglish: "Save",
			LanguageRussian: "Сохранить",
		},
		keyButtonCancel: {
			LanguageEnglish: "Cancel",
			LanguageRussian: "Отмена",
		},
		keyButtonSaveUser: {
			LanguageEnglish: "Save",
			LanguageRussian: "Сохранить",
		},

		// form fields
		keyFieldEmail: {
			LanguageEnglish: "Email",
			LanguageRussian: "Имейл",
		},
		keyFieldPassword: {
			LanguageEnglish: "Password",
			LanguageRussian: "Пароль",
		},
		keyFieldSearch: {
			LanguageEnglish: "Search",
			LanguageRussian: "Поиск",
		},
		keyFieldFirstName: {
			LanguageEnglish: "First name",
			LanguageRussian: "Имя",
		},
		keyFieldLastName: {
			LanguageEnglish: "Last name",
			LanguageRussian: "Фамилия",
		},
		keyFieldBirthDate: {
			LanguageEnglish: "Birth date",
			LanguageRussian: "Дата рождения",
		},
		keyFieldGender: {
			LanguageEnglish: "Gender",
			LanguageRussian: "Пол",
		},
		keyFieldGenderMale: {
			LanguageEnglish: "Male",
			LanguageRussian: "Мужской",
		},
		keyFieldGenderFemale: {
			LanguageEnglish: "Female",
			LanguageRussian: "Женский",
		},
		keyFieldAddress: {
			LanguageEnglish: "Address",
			LanguageRussian: "Адрес",
		},
		keyFieldLanguage: {
			LanguageEnglish: "Language",
			LanguageRussian: "Язык",
		},
		keyFieldLanguageEn: {
			LanguageEnglish: "English",
			LanguageRussian: "Английский",
		},
		keyFieldLanguageRu: {
			LanguageEnglish: "Russian",
			LanguageRussian: "Русский",
		},

		// links
		keyLinkLogin: {
			LanguageEnglish: "Already have an account?",
			LanguageRussian: "Уже есть аккаунт?",
		},
		keyLinkRegister: {
			LanguageEnglish: "Don't have an account? Create one now.",
			LanguageRussian: "Нет аккаунта? Создать учетную запись.",
		},

		// table columns
		keyColumnName: {
			LanguageEnglish: "Name",
			LanguageRussian: "Имя",
		},
		keyColumnEmail: {
			LanguageEnglish: "Email",
			LanguageRussian: "Имейл",
		},
		keyColumnBirthDate: {
			LanguageEnglish: "Birth date",
			LanguageRussian: "Дата рождения",
		},
		keyColumnGender: {
			LanguageEnglish: "Gender",
			LanguageRussian: "Пол",
		},
		keyColumnAddress: {
			LanguageEnglish: "Address",
			LanguageRussian: "Адрес",
		},

		keyColumnJanuary: {
			LanguageEnglish: "January",
			LanguageRussian: "Январь",
		},
		keyColumnFebruary: {
			LanguageEnglish: "February",
			LanguageRussian: "Февраль",
		},
		keyColumnMarch: {
			LanguageEnglish: "March",
			LanguageRussian: "Март",
		},
		keyColumnApril: {
			LanguageEnglish: "April",
			LanguageRussian: "Апрель",
		},
		keyColumnMay: {
			LanguageEnglish: "May",
			LanguageRussian: "Май",
		},
		keyColumnJune: {
			LanguageEnglish: "June",
			LanguageRussian: "Июнь",
		},
		keyColumnJuly: {
			LanguageEnglish: "July",
			LanguageRussian: "Июль",
		},
		keyColumnAugust: {
			LanguageEnglish: "August",
			LanguageRussian: "Август",
		},
		keyColumnSeptember: {
			LanguageEnglish: "September",
			LanguageRussian: "Сентябрь",
		},
		keyColumnOctober: {
			LanguageEnglish: "October",
			LanguageRussian: "Октябрь",
		},
		keyColumnNovember: {
			LanguageEnglish: "November",
			LanguageRussian: "Ноябрь",
		},
		keyColumnDecember: {
			LanguageEnglish: "December",
			LanguageRussian: "Декабрь",
		},

		// error messages
		KeyErrorEmailTaken: {
			LanguageEnglish: "Email address is already in use",
			LanguageRussian: "Имейл адрес уже используется",
		},
		KeyErrorCredentialsInvalid: {
			LanguageEnglish: "Credentials invalid",
			LanguageRussian: "Неверный email или пароль",
		},
		KeyErrorDataCollision: {
			LanguageEnglish: "Data collision",
			LanguageRussian: "Коллизия данных",
		},

		// other
		keyNoSearchResults: {
			LanguageEnglish: "No customers",
			LanguageRussian: "Нет клиентов",
		},
	}

	errLangNotSupported = errors.New("target translation language not supported")
	errKeyNotDefined    = errors.New("translation key not defined")
	errLangNotDefined   = errors.New("language not defined for key")
)

func Translate(key Key, language string) (string, error) {
	if strings.TrimSpace(language) == "" || (language != LanguageEnglish && language != LanguageRussian) {
		return "", fmt.Errorf("%w: %s", errLangNotSupported, language)
	}

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
