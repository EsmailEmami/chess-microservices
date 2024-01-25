package validations

import (
	"errors"
	"regexp"
	"strings"

	"github.com/esmailemami/chess/shared/consts"
)

// StrongPassword check password to be strong enough but does not return error if password is empty
// It must be used along side with Required validation if empty password is not allowed
// On update, most of time password not changing, so it must not return error
// when it's empty
func StrongPassword() func(value interface{}) error {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		pass, ok := value.(string)
		if !ok {
			return nil
		}

		if pass == "" {
			return nil
		}

		if len(pass) < 8 {
			return errors.New(consts.PasswordIsShort)
		}

		// A-Z
		if ok, _ := regexp.MatchString("[A-Z]{1,}", pass); !ok {
			return errors.New(consts.PasswordIsShort)
		}
		// a-z
		if ok, _ := regexp.MatchString("[a-z]{1,}", pass); !ok {
			return errors.New(consts.PasswordIsShort)
		}
		// 0-9
		if ok, _ := regexp.MatchString("[0-9]{1,}", pass); !ok {
			return errors.New(consts.PasswordIsShort)
		}
		// special char
		sb := strings.Builder{}
		sb.WriteString("!@#$%^&*()-=¡£_+`~.,<>/?;:'|[]{}")
		sb.WriteRune('"')
		pattern := "[" + regexp.QuoteMeta(sb.String()) + "]{1,}"
		if ok, _ := regexp.MatchString(pattern, pass); !ok {
			return errors.New(consts.PasswordIsShort)
		}

		return nil
	}
}
