package request

import (
	"github.com/asaskevich/govalidator"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"

	dateFormat = "2006-01-02"
)

func init() {
	govalidator.TagMap["gender"] = func(str string) bool {
		return str == GenderMale || str == GenderFemale
	}

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
	govalidator.ParamTagRegexMap["age"] = regexp.MustCompile("^age\\((\\d+)\\|(\\d+)\\)$")
}

func validateStruct(s interface{}) (bool, []string) {
	res, err := govalidator.ValidateStruct(s)

	if err != nil {
		errs := err.(govalidator.Errors).Errors()
		fields := make([]string, 0, len(errs))

		for i := range errs {
			m := errs[i].Error()
			fields = append(fields, m[:strings.IndexByte(m, ':')])
		}

		return res, fields
	}

	return res, nil
}
