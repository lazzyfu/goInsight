package dao

import (
	"fmt"
	"regexp"
)

var identifierPattern = regexp.MustCompile(`^[0-9A-Za-z_$]+$`)

func validateIdentifier(value, field string) error {
	if !identifierPattern.MatchString(value) {
		return fmt.Errorf("%s包含非法字符", field)
	}
	return nil
}
