package request

import "strings"

type RegisterUser struct {
	Email    string `form:"email" json:"email" valid:"email,maxstringlength(255),required"`
	Password string `form:"password" json:"password" valid:"printableascii,maxstringlength(32),required"`
}

type LoginUser struct {
	Email    string `form:"email" json:"email" valid:"email,maxstringlength(255),required"`
	Password string `form:"password" json:"password" valid:"printableascii,maxstringlength(32),required"`
}

func (u *RegisterUser) Validate() []string {
	return validateStruct(*u)
}

func (u *RegisterUser) Sanitized() RegisterUser {
	return RegisterUser{
		Email:    strings.ToLower(strings.TrimSpace(u.Email)),
		Password: strings.TrimSpace(u.Password),
	}
}

func (u *LoginUser) Validate() []string {
	return validateStruct(*u)
}
