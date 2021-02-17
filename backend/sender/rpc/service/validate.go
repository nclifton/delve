package service

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
)

func (s *senderImpl) validateCSVSenders(ctx context.Context, csvSenders []SenderCSV) (validSenders []SenderCSV, results []SenderCSV, err error) {

	validSenders = make([]SenderCSV, 0, len(csvSenders))
	results = make([]SenderCSV, 0, len(csvSenders))
	for _, csvSender := range csvSenders {
		valid, validatedNewSender, err := s.validateCSVSender(ctx, csvSender)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, validatedNewSender)
		if valid {
			validSenders = append(validSenders, validatedNewSender)
		}
	}

	return validSenders, results, nil
}

func (s *senderImpl) validateCSVSender(ctx context.Context, csvSender SenderCSV) (bool, SenderCSV, error) {

	validator := &validator{
		db:     s.db,
		errors: []string{},
	}

	err := validator.validateAddress(ctx, csvSender.Address)
	if err != nil {
		return false, csvSender, err
	}

	csvSender.Status = validator.getStatus()
	csvSender.Error = validator.getFirstError()

	return len(validator.errors) == 0, csvSender, nil
}

type validator struct {
	db     db.DB
	errors []string
}

func (v *validator) validateAddress(ctx context.Context, address string) error {

	if len(address) == 0 {
		v.errors = append(v.errors, `Field "address" cannot be empty`)
		return nil
	}
	senders, err := v.db.FindSendersByAddress(ctx, address)
	if err != nil {
		return err
	}
	if len(senders) > 0 {
		v.errors = append(v.errors, `Field "address" must be unique`)
		return nil
	}

	return nil

}

func (v *validator) getStatus() string {
	if len(v.errors) == 0 {
		return "ok"
	}
	return "skipped"
}

func (v *validator) getFirstError() string {
	if len(v.errors) == 0 {
		return ""
	}
	return v.errors[0]
}
