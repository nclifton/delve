package valid

import (
	"errors"
	"fmt"
	"log"
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
	"contains":    Contains,

	"webhook_url": IsWebhookURL,
}

var reservedIPNets []*net.IPNet

func init() {
	// construct lookup table of reserved IP addresses for webhook URL validation
	reservedCIDRs := []string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"240.0.0.0/4",
		"255.255.255.255/32",
	}
	for _, cidr := range reservedCIDRs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("CIDR parse error: %s", err)
		}

		reservedIPNets = append(reservedIPNets, ipnet)
	}
}

/**
 *     - if the value is a string it cannot be blank (""), an empty string
 *	- if the value is an array it cannot be empty, length 0
 *	- if an array of strings, the strings in the array cannot be blank, an empty string
 *	- if not a string, the value or values if an array, cannot be the type's zero value
 */
func IsRequired(i interface{}, parent interface{}, params []string) error {
	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return errors.New("required")
		}
	case reflect.Map, reflect.Slice:
		if v.IsNil() || v.Len() == 0 {
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

/**
 *  - Contains parameter is an array of strings.
 *  - Validation is true if the value equals one of the strings in the params string array.
 *  - This is effectively a "one-of" validation and not a string contains string validation
 */

func Contains(i interface{}, parent interface{}, params []string) error {
	if len(params) == 0 {
		return errors.New("expected at least 1 param to compare against")
	}

	str, ok := i.(string)
	if !ok {
		return errors.New("expected string type")
	}

	match := false
	for _, candidate := range params {
		if candidate == str {
			match = true
		}
	}
	if !match {
		return fmt.Errorf("%s did not match any of %s", i, strings.Join(params, ","))
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

// IsWebhookURL checks that the given value is a valid url to be used as a webhook.
// It first checks that i is a valid url via IsURL
// It then looks up the ip addresses of the host and ensures that the ip does not
// match the range reservedIPNets defined in the init of this package
// this makes sure that a user cannot set a webhook that could possibly access areas
// of our internal system
func IsWebhookURL(i interface{}, parent interface{}, params []string) error {

	err := IsURL(i, parent, params)
	if err != nil {
		return err
	}

	key, ok := i.(string)
	if !ok {
		return errors.New("expected string type")
	}

	parsedURL, err := url.Parse(key)
	if err != nil {
		return errors.New("invalid URL")
	}

	addrs, err := net.LookupHost(parsedURL.Hostname())
	if err != nil {
		return errors.New("failed to lookup host")
	}

	for _, ipResolved := range addrs {
		ip := net.ParseIP(ipResolved)

		for _, ipReserved := range reservedIPNets {
			if ipReserved.Contains(ip) {
				return errors.New("URL resolves to reserved IP")
			}
		}
	}

	return nil
}
