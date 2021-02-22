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
			return fmt.Errorf("is not one of %s", strings.Join(enumValues, "|"))
		},
		ExcludeKinds: []reflect.Kind{reflect.Array, reflect.Slice},
	}
}
