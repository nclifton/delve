package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/burstsms/mtmo-tp/backend/lib/valid"
)

/**
 * Address string `valid:"address_unique_in_upload"`
 *
 * **NOTE:** this is only intended for the address field.
 * Reflection would be required to be able to use this to validate other fields
 */
func addressOccurrenceValidator(csvSenders []SenderCSV) valid.CustomValidator {
	return valid.CustomValidator{
		Name: "address_unique_in_upload",
		Fn: func(i interface{}, parent interface{}, params []string) error {
			cnt := 0
			for _, csvSender := range csvSenders {
				if csvSender.Address == i.(string) {
					cnt++
				}
			}
			if cnt > 1 {
				return errors.New("multiple occurrence in upload")
			}
			return nil
		},
		ExcludeKinds: []reflect.Kind{reflect.Array, reflect.Slice},
	}

}

/**
 * Address string `valid:"address_new"`
 *
 * **NOTE:** this is only intended for the address field.
 * Checks that the specified address does not exists in the sender database sender table
 */
func (s *senderImpl) addressDbValidator(ctx context.Context) valid.CustomValidator {
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

/**
 * `valid:sender_enum(<database enum type name>)`
 *
 * Checks that the field value value matches one of the defined values for the named database enum type.
 * This validator does not assume the role of the "required" validator, if the field is empty this validator will not validate the field.
 *
 */
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

/**
 * ... `valid:"required_if(<field name>[|<rule>[|rule parameters]])`
 *
 * conditional required. The field being validated is validated as required if another field passes a specified validation
 *
 * example:
 *
 * MMSProviderKey string `valid:required_id(Channels|contains|mms)`
 *
 * The field `MMSProviderKey` will be "required" if the `Channels` field "contains" `mms`.
 * This rule uses the `contains` rule to check that the `Channels` field contains the value `mms`.
 * The check field, in this case `Channels`, may be an array field
 *
 * reflect is used by this validator function.
 *
 */
func requiredIfValidator() valid.CustomValidator {
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
