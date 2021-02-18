package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	tagName      = "valid"
	paramsRegexp = regexp.MustCompile(`\((.*)\)$`)
)

func Validate(s interface{}) error {
	if s == nil {
		return fmt.Errorf("can not validate nil")
	}

	var errs Errors

	// defref the given value if pointer
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// we only accept structs
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("can only validate a struct, got %s", val.Kind())
	}

	// iterate using the number of fields in the struct
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		// skip private fields
		if typeField.PkgPath != "" {
			continue
		}

		// explicitly skipped fields
		if typeField.Tag.Get(tagName) == "-" {
			continue
		}

		// lets now just validate as a field
		if err := validateField(valueField, typeField, val, []validator{}); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func validateField(value reflect.Value, field reflect.StructField, parent reflect.Value, validators []validator) error {
	if !value.IsValid() {
		return fmt.Errorf("invalid value (%s) in field (%s)", value.String(), field.Name)
	}

	// you can define a validator for a slice or map type and hve it apply to each value
	validators = append(validators, getValidators(field.Tag.Get(tagName))...)

	// check the field kind to determine if we need to recurse or iterate
	kind := value.Kind()
	isStruct := kind == reflect.Struct
	isStructPtr := kind == reflect.Ptr && value.Elem().Kind() == reflect.Struct

	if isStruct || isStructPtr {
		// recurse !
		if err := Validate(value.Interface()); err != nil {
			err = PrependPathToErrors(err, field.Name)
			return err
		}

		return nil
	}

	// if we have an interface or pointer type we use the underlying element
	if kind == reflect.Ptr || kind == reflect.Interface {
		// check the value pointed to
		if !value.IsNil() {
			if err := validateField(value.Elem(), field, parent, validators); err != nil {
				return err
			}
		}

		return nil
	}

	if kind == reflect.Array || kind == reflect.Slice {
		// validate as an array (required, length?)
		if err := validate(value, field, parent, validators); err != nil {
			return err
		}

		// iterate the items
		for i := 0; i < value.Len(); i++ {
			if err := validateField(value.Index(i), field, parent, validators); err != nil {
				return err
			}
		}

		return nil
	}

	if kind == reflect.Map {
		// checking all the map values
		for _, k := range value.MapKeys() {
			if err := validateField(value.MapIndex(k), field, parent, validators); err != nil {
				return err
			}
		}

		return nil
	}

	if kind == reflect.Chan || kind == reflect.Func || kind == reflect.UnsafePointer {
		return fmt.Errorf("can not validate field kind: %s", kind)
	}

	// just a regular field in this case
	return validate(value, field, parent, validators)
}

func validate(value reflect.Value, field reflect.StructField, parent reflect.Value, validators []validator) error {
	for _, v := range validators {
		fn, exists := TagMap[v.name]
		if !exists {
			return fmt.Errorf("%s is not a defined validator function", v.name)
		}

		if err := fn(value.Interface(), parent.Interface(), v.params); err != nil {
			return Error{jsonName(field.Tag.Get("json"), field.Name), err, v.name, []string{}}
		}
	}

	return nil
}

func getValidators(tag string) []validator {
	validators := []validator{}

	if tag == "" {
		return validators
	}

	for _, name := range strings.Split(tag, ",") {
		v := validator{}
		v.name = strings.TrimSpace(name)
		v.name = paramsRegexp.ReplaceAllString(v.name, "")

		matches := paramsRegexp.FindStringSubmatch(name)
		if len(matches) == 2 {
			v.params = strings.Split(matches[1], "|")
		}

		validators = append(validators, v)
	}

	return validators
}

func jsonName(tag, fieldName string) string {
	if tag == "" {
		return fieldName
	}

	// JSON name always comes first. If there's no options then split[0] is
	// JSON name, if JSON name is not set, then split[0] is an empty string.
	split := strings.SplitN(tag, ",", 2)
	name := split[0]

	// However it is possible that the field is skipped when
	// (de-)serializing from/to JSON, in which case assume that there is no
	// tag name to use
	if name == "-" {
		return fieldName
	}

	return name
}
