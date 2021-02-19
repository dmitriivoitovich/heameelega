package request

import (
	"regexp"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"

	dateFormat = "2006-01-02"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.TagMap["gender"] = func(str string) bool {
		return str == GenderMale || str == GenderFemale
	}

	govalidator.ParamTagMap["inrange"] = func(str string, params ...string) bool {
		value, err := strconv.Atoi(str)
		if err != nil {
			return false
		}

		min, err := strconv.Atoi(params[0])
		if err != nil {
			return false
		}

		max, err := strconv.Atoi(params[1])
		if err != nil {
			return false
		}

		return value >= min && value <= max
	}
	govalidator.ParamTagRegexMap["inrange"] = regexp.MustCompile(`^inrange\((\d+)\|(\d+)\)$`)

	govalidator.ParamTagMap["age"] = func(str string, params ...string) bool {
		date, err := time.Parse(dateFormat, str)
		if err != nil {
			return false
		}

		minAge, err := strconv.Atoi(params[0])
		if err != nil {
			return false
		}

		maxAge, err := strconv.Atoi(params[1])
		if err != nil {
			return false
		}

		min := time.Now().AddDate(-maxAge, 0, 0)
		max := time.Now().AddDate(-minAge, 0, 0)

		return date.After(min) && date.Before(max)
	}
	govalidator.ParamTagRegexMap["age"] = regexp.MustCompile(`^age\((\d+)\|(\d+)\)$`)
}

func validateStruct(s interface{}) []string {
	if _, err := govalidator.ValidateStruct(s); err != nil {
		errs := govalidator.ErrorsByField(err)

		fields := make([]string, 0, len(errs))
		for i := range errs {
			fields = append(fields, i)
		}

		return fields
	}

	return nil
}
