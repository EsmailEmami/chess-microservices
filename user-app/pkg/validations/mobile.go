package validations

import (
	"errors"
	"regexp"

	"github.com/esmailemami/chess/user/pkg/consts"
)

func IsValidMobileNumber() func(value interface{}) error {
	return func(value interface{}) error {
		mobile, ok := value.(string)
		if !ok {
			return errors.New(consts.InvalidMobileNumber)
		}

		if match, _ := regexp.MatchString("^09[0-9]{9}$", mobile); match {
			return nil
		}
		return errors.New(consts.InvalidMobileNumber)
	}
}
