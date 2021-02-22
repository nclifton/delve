package service

import (
	"context"
	"errors"
	"fmt"

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
	}
}
