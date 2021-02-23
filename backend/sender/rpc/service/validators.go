package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/burstsms/mtmo-tp/backend/lib/valid"
)

func (s *senderImpl) addressValidator(ctx context.Context) valid.CustomValidator {
	return valid.CustomValidator{
		Name: "address_new",
		Fn: func(i interface{}, parent interface{}, params []string) error {
			exists, err := s.db.SenderAddressExists(ctx, fmt.Sprintf("%v", i))
			if err != nil {
				return err
			}
			if exists {
				return errors.New("is not new")
			}
			return nil
		},
		ExcludeKinds: []reflect.Kind{reflect.Array, reflect.Slice},
	}
}

func (s *senderImpl) senderEnumValidator(ctx context.Context) valid.CustomValidator {
	return valid.CustomValidator{
		Name: "sender_enum",
		Fn: func(i interface{}, parent interface{}, params []string) error {
			if len(params) == 0 {
				return errors.New("rule error: enum type not specified")
			}
			v := reflect.ValueOf(i)
			if v.Kind() != reflect.String {
				return errors.New("value is not a string")
			}
			value := fmt.Sprintf("%s", i)
			if value == "" {
				return nil
			}
			var err error
			enums, err := s.db.GetSenderEnums(ctx)
			if err != nil {
				return err
			}
			enumValues, inEnums := enums[params[0]]
			if !inEnums {
				return fmt.Errorf("rule error: enum type %s is not defined", params[0])
			}
			for _, enumValue := range enumValues {
				if value == enumValue {
					return nil
				}
			}
			return fmt.Errorf("%s did not match any of %s", i.(string), strings.Join(enumValues, ","))
		},
		ExcludeKinds: []reflect.Kind{reflect.Array, reflect.Slice},
	}
}

func requiredIf() valid.CustomValidator {
	return valid.CustomValidator{
		Name: "required_if",
		Fn: func(i interface{}, parent interface{}, params []string) error {

			// parse params
			var checkFieldName, checkFieldRule, checkRuleInMessage string
			var checkRuleArgs []string
			switch len(params) {
			case 0:
				return fmt.Errorf("rule error: required_if: not enough parameters")
			case 1:
				checkFieldName = params[0]
				checkFieldRule = "required"
				checkRuleInMessage = "is present"
			default:
				checkFieldName = params[0]
				checkFieldRule = params[1]
				checkRuleInMessage = params[1]
				checkRuleArgs = params[2:]
			}

			// find the check field in the parent
			// parent needs to be a reflect Value and a struct
			parentValue := reflect.ValueOf(parent)
			checkField := parentValue.FieldByName(checkFieldName)
			if !checkField.IsValid() {
				return fmt.Errorf("rule error: required_if: field %s is not defined", checkFieldName)
			}
			checkFieldFn, defined := valid.TagMap[checkFieldRule]
			if !defined {
				return fmt.Errorf("rule error: required_if: sub-rule %s is not defined", checkFieldRule)
			}

			// the check field may be an array or a struct or a pointer ... (this is normally handled by valid::validateField, which is un-exposed so we repeat it here)
			checkFieldKind := checkField.Kind()

			if checkFieldKind == reflect.Ptr || checkFieldKind == reflect.Interface {
				checkField = checkField.Elem()
				checkFieldKind = checkField.Kind()
			}

			if err := validateField(checkField, parent, checkFieldKind, checkFieldRule, checkFieldFn, checkRuleArgs); err != nil {
				return nil // failed required if check condition, therefore field is not required
			}

			// now apply the required rule to the field because the check field condition tells us that the field is required
			requiredFn, defined := valid.TagMap["required"]
			if !defined {
				return fmt.Errorf(`rule error: required_if: "required" rule is not defined`)
			}
			field := reflect.ValueOf(i)
			kind := field.Kind()
			if err := validateField(field, parent, kind, "required", requiredFn, []string{}); err != nil {
				return fmt.Errorf(`required if %s %s %s`, checkFieldName, checkRuleInMessage, strings.Join(checkRuleArgs, ","))
			}

			return nil

		},
	}

}

func validateField(field reflect.Value, parent interface{}, fieldKind reflect.Kind, fieldRule string, ruleFn valid.ValidatorFunc, ruleArgs []string) error {

	excludeKinds, has := valid.RuleExcludeKinds[fieldRule]
	excludeRule := false
	if has {
		for _, excludeKind := range excludeKinds {
			if excludeKind == fieldKind {
				excludeRule = true
			}
		}
	}
	if !excludeRule {
		if err := ruleFn(field.Interface(), parent, ruleArgs); err != nil {
			return err
		}
	}
	if fieldKind == reflect.Array || fieldKind == reflect.Slice {

		// check each item in the array or slice
		for i := 0; i < field.Len(); i++ {
			if err := validateField(field.Index(i), parent, field.Index(i).Kind(), fieldRule, ruleFn, ruleArgs); err != nil {
				return err
			}
		}
	}

	return nil
}
