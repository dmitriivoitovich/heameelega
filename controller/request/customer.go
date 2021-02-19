package request

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultSearchOrderField = "first_name"
	DefaultSearchDirection  = "asc"
)

type SearchCustomers struct {
	Filter    *string `query:"s" json:"s" valid:"printableascii,maxstringlength(201),optional"`
	Order     *string `query:"o" json:"o" valid:"printableascii,in(first_name|last_name|birth_date|gender|email|address),optional"`
	Direction *string `query:"d" json:"d" valid:"printableascii,in(asc|desc),optional"`
	Page      *uint32 `query:"p" json:"p" valid:"numeric,inrange(1|99999),optional"`
}

type SearchCustomersNormalized struct {
	Filter    string
	Order     string
	Direction string
	Page      uint32
}

type CreateCustomer struct {
	FirstName string `form:"firstName" json:"firstName" valid:"printableascii,maxstringlength(100),required"`
	LastName  string `form:"lastName" json:"lastName" valid:"printableascii,maxstringlength(100),required"`
	BirthDate string `form:"birthDate" json:"birthDate" valid:"age(18|60),required"`
	Gender    string `form:"gender" json:"gender" valid:"gender,required"`
	Email     string `form:"email" json:"email" valid:"email,maxstringlength(255),required"`
	Address   string `form:"address" json:"address" valid:"printableascii,maxstringlength(200),optional"`
}

type EditCustomer struct {
	ID        uuid.UUID `param:"id" json:"id" valid:"required"`
	FirstName string    `form:"firstName" json:"firstName" valid:"printableascii,maxstringlength(100),required"`
	LastName  string    `form:"lastName" json:"lastName" valid:"printableascii,maxstringlength(100),required"`
	BirthDate string    `form:"birthDate" json:"birthDate" valid:"age(18|60),required"`
	Gender    string    `form:"gender" json:"gender" valid:"gender,required"`
	Email     string    `form:"email" json:"email" valid:"email,maxstringlength(255),required"`
	Address   string    `form:"address" json:"address" valid:"printableascii,maxstringlength(200),optional"`
	LoadedAt  time.Time `form:"loadedAt" json:"loadedAt" valid:"required"`
}

func (c *CreateCustomer) BirthDateTime() (*time.Time, error) {
	date, err := time.Parse(dateFormat, c.BirthDate)
	if err != nil {
		return nil, err
	}

	return &date, nil
}

func (c *CreateCustomer) Validate() []string {
	return validateStruct(*c)
}

func (c *EditCustomer) Validate() []string {
	return validateStruct(*c)
}

func (c *EditCustomer) BirthDateTime() (time.Time, error) {
	return time.Parse(dateFormat, c.BirthDate)
}

func (c *SearchCustomers) Validate() []string {
	fields := validateStruct(*c)

	if c.Page != nil && *c.Page == 0 {
		fields = append(fields, "p")
	}

	return fields
}

func (c *SearchCustomers) Normalized() SearchCustomersNormalized {
	invalidFields := c.Validate()

	invalidFieldsMap := make(map[string]bool, len(invalidFields))
	for i := range invalidFields {
		invalidFieldsMap[invalidFields[i]] = true
	}

	filter := ""
	if c.Filter != nil && !invalidFieldsMap["s"] {
		filter = strings.TrimSpace(*c.Filter)
	}

	order := DefaultSearchOrderField
	if c.Order != nil && !invalidFieldsMap["o"] {
		order = strings.TrimSpace(*c.Order)
	}

	direction := DefaultSearchDirection
	if c.Direction != nil && !invalidFieldsMap["d"] {
		direction = strings.TrimSpace(*c.Direction)
	}

	page := uint32(1)
	if c.Page != nil && !invalidFieldsMap["p"] {
		page = *c.Page
	}

	res := SearchCustomersNormalized{
		Filter:    filter,
		Order:     order,
		Direction: direction,
		Page:      page,
	}

	return res
}
