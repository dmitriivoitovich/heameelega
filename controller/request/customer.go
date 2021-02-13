package request

import (
	"time"
)

type CreateCustomer struct {
	FirstName string `form:"firstName" json:"firstName" valid:"type(string),printableascii,maxstringlength(100),required"`
	LastName  string `form:"lastName" json:"lastName" valid:"type(string),printableascii,maxstringlength(100),required"`
	BirthDate string `form:"birthDate" json:"birthDate" valid:"type(string),age(18|60),required"`
	Gender    string `form:"gender" json:"gender" valid:"type(string),gender,required"`
	Email     string `form:"email" json:"email" valid:"type(string),email,maxstringlength(255),required"`
	Address   string `form:"address" json:"address" valid:"type(string),printableascii,maxstringlength(200),optional"`
}

func (c *CreateCustomer) BirthDateTime() (*time.Time, error) {
	date, err := time.Parse(dateFormat, c.BirthDate)
	if err != nil {
		return nil, err
	}

	return &date, nil
}

func (c *CreateCustomer) Validate() (bool, []string) {
	return validateStruct(*c)
}
