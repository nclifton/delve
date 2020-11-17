package valid

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidatorFunc func(i interface{}, parent interface{}, params []string) error
type validator struct {
	name   string
	params []string
}

var TagMap = map[string]ValidatorFunc{
	"required":    IsRequired,
	"url":         IsURL,
	"email":       IsEmail,
	"integer":     IsInteger,
	"alpha":       IsAlpha,
	"length":      Length,
	"rune_length": RuneLength,
	"range":       Range,
}

var reservedIPNets []*net.IPNet

func IsRequired(i interface{}, parent interface{}, params []string) error {
	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return errors.New("required")
		}
	case reflect.String:
		if v.IsZero() {
			return errors.New("required")
		}
	}

	return nil
}

func IsURL(i interface{}, parent interface{}, params []string) error {
	maxURLRuneCount := 2083
	minURLRuneCount := 3

	str, ok := i.(string)
	if !ok {
		return errors.New("expected a string type")
	}

	if str == "" {
		return nil
	}

	if strings.HasPrefix(str, ".") {
		return errors.New("can not begin with .")
	}

	length := utf8.RuneCountInString(str)
	if length <= minURLRuneCount || length >= maxURLRuneCount {
		return fmt.Errorf("length must be between %d and %d", minURLRuneCount, maxURLRuneCount)
	}

	strTemp := str
	if strings.Contains(str, ":") && !strings.Contains(str, "://") {
		// support no indicated urlscheme but with colon for port number
		// http:// is appended so url.Parse will succeed, strTemp used so it does not impact rxURL.MatchString
		strTemp = "http://" + str
	}

	u, err := url.Parse(strTemp)
	if err != nil {
		return err
	}

	if strings.HasPrefix(u.Host, ".") {
		return errors.New("host can not begin with .")
	}

	if match := rxURL.MatchString(str); !match {
		return errors.New("invalid as per regexp")
	}

	return nil
}

func IsInteger(i interface{}, parent interface{}, params []string) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("expected string type")
	}

	_, err := strconv.Atoi(str)
	return err
}

// TODO uppercase letters are not supported
func IsEmail(i interface{}, parent interface{}, params []string) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("expected string type")
	}

	if str != "" && !rxEmail.MatchString(str) {
		return errors.New("invalid email address")
	}

	return nil
}

func IsAlpha(i interface{}, parent interface{}, params []string) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("expected string type")
	}

	if str != "" && !rxAlpha.MatchString(str) {
		return errors.New("invalid as per regexp")
	}

	return nil
}

func Length(i interface{}, parent interface{}, params []string) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("expected string type")
	}

	return checkLength(len(str), params)
}

func RuneLength(i interface{}, parent interface{}, params []string) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("expected string type")
	}

	return checkLength(utf8.RuneCountInString(str), params)
}

func checkLength(strLen int, params []string) error {
	if len(params) != 2 {
		return errors.New("expected 2 params to validate length")
	}

	min, err := strconv.Atoi(params[0])
	if err != nil {
		return errors.New("first param did not parse as integer")
	}

	max, err := strconv.Atoi(params[1])
	if err != nil {
		return errors.New("second param did not parse as integer")
	}

	if strLen < min || strLen > int(max) {
		return fmt.Errorf("length must be between %d and %d", min, max)
	}

	return nil
}

func Range(i interface{}, parent interface{}, params []string) error {
	if len(params) != 2 {
		return fmt.Errorf("expected only 2 params for range, was given: %v", params)
	}

	switch t := i.(type) {
	case float32, float64:
		return inRangeFloat64(i.(float64), params)
	case int, int32, int64:
		return inRangeInt(i.(int), params)
	default:
		return fmt.Errorf("can not validate range for type: %T\n", t)
	}
}

func inRangeFloat64(val float64, params []string) error {
	min, err := strconv.ParseFloat(params[0], 64)
	if err != nil {
		return errors.New("first param did not parse as float64")
	}

	max, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return errors.New("second param did not parse as float64")
	}

	if val < min {
		return fmt.Errorf("less than the minimum value of %f", min)
	}

	if val > max {
		return fmt.Errorf("greather than the maximum value of %f", max)
	}

	return nil
}

func inRangeInt(val int, params []string) error {
	min, err := strconv.Atoi(params[0])
	if err != nil {
		return errors.New("first param did not parse as int")
	}

	max, err := strconv.Atoi(params[1])
	if err != nil {
		return errors.New("second param did not parse as int")
	}

	if val < min {
		return fmt.Errorf("less than the minimum value of %d", min)
	}

	if val > max {
		return fmt.Errorf("greather than the maximum value of %d", max)
	}

	return nil
}

func paramsContains(params []string, param string) bool {
	for _, val := range params {
		if val == param {
			return true
		}
	}

	return false
}
