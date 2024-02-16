package validations

import (
	"errors"
	"regexp"
	"strings"

	"github.com/esmailemami/chess/shared/util"
	"github.com/esmailemami/chess/user/internal/consts"
)

func IsValidMobileNumber() func(value interface{}) error {
	return func(value interface{}) error {
		if util.IsNil(value) {
			return nil
		}

		mobile, ok := util.Value(value).(string)

		if strings.TrimSpace(mobile) == "" {
			return nil
		}

		if !ok {
			return errors.New(consts.InvalidMobileNumber)
		}

		if match, _ := regexp.MatchString("^09[0-9]{9}$", mobile); match {
			return nil
		}
		return errors.New(consts.InvalidMobileNumber)
	}
}
