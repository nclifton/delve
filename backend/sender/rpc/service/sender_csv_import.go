package service

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gocarina/gocsv"
	"github.com/vincent-petithory/dataurl"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/valid"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type SenderCSV struct {
	AccountId      string       `csv:"account_id" valid:""`
	Address        string       `csv:"address" valid:"required"`
	Country        string       `csv:"country" valid:"required"`
	Channels       CSVJSONArray `csv:"channels" valid:"required"` // see custom CSV Field conversion below
	MMSProviderKey string       `csv:"mms_provider_key"`
	Comment        string       `csv:"comment"`
	Status         string       `csv:"status"`
	Error          string       `csv:"error"`
}

const (
	CSV_STATUS_SKIPPED = "skipped"
	CSV_STATUS_OK      = "ok"
)

func (s *senderImpl) CreateSendersFromCSVDataURL(ctx context.Context, r *senderpb.CreateSendersFromCSVDataURLParams) (*senderpb.CreateSendersFromCSVDataURLReply, error) {

	replySenders := []*senderpb.Sender{}

	csvSenders, err := unmarshalSenderCSVDataUrl(r.CSV)
	if err != nil {
		return nil, err
	}

	// ignoring results here at the moment until we have the validation available and working how we want it
	validSenders, _ := s.validateCSVSenders(ctx, csvSenders)

	if len(validSenders) > 0 {

		dbSenders, err := s.db.InsertSenders(ctx, validSenders)
		if err != nil {
			return nil, err
		}

		for _, dbSender := range dbSenders {
			replySenders = append(replySenders, dbSenderToSender(dbSender))
		}
	}

	return &senderpb.CreateSendersFromCSVDataURLReply{
		Senders: replySenders,
	}, nil
}

func dbSenderToSender(sender db.Sender) *senderpb.Sender {
	return &senderpb.Sender{
		Id:             sender.ID,
		AccountId:      sender.AccountID,
		Address:        sender.Address,
		MMSProviderKey: sender.MMSProviderKey,
		Channels:       sender.Channels,
		Country:        sender.Country,
		Comment:        sender.Comment,
		CreatedAt:      timestamppb.New(sender.CreatedAt),
		UpdatedAt:      timestamppb.New(sender.UpdatedAt),
	}
}

func unmarshalSenderCSVDataUrl(csvDataUrl []byte) (csvSenders []SenderCSV, err error) {

	data, err := dataurl.DecodeString(string(csvDataUrl))
	if err != nil {
		return nil, err
	}

	csvSenders = []SenderCSV{}

	reader := gocsv.LazyCSVReader(bytes.NewReader(data.Data))
	err = gocsv.UnmarshalCSV(reader, &csvSenders)
	if err != nil {
		return nil, err
	}

	return csvSenders, nil

}

func (s *senderImpl) validateCSVSenders(ctx context.Context, csvSenders []SenderCSV) ([]db.Sender, []SenderCSV) {
	validSenders := make([]db.Sender, 0, len(csvSenders))
	validatedCSVSenders := make([]SenderCSV, 0, len(csvSenders))
	for _, csvSender := range csvSenders {
		err := valid.Validate(csvSender)
		if err != nil {
			csvSender.Status = CSV_STATUS_SKIPPED
			csvSender.Error = err.Error()
		} else {
			validSenders = append(validSenders, db.Sender{
				AccountID:      csvSender.AccountId,
				Address:        csvSender.Address,
				MMSProviderKey: csvSender.MMSProviderKey,
				Channels:       csvSender.Channels,
				Country:        csvSender.Country,
				Comment:        csvSender.Comment,
			})
		}
		validatedCSVSenders = append(validatedCSVSenders, csvSender)
	}
	return validSenders, validatedCSVSenders
}

/*
below is for a custom conversion to be used by the CSV marshalling and un-marshalling
*/

type CSVJSONArray []string

// Convert the internal string array to JSON string
func (a *CSVJSONArray) MarshalCSV() (string, error) {
	str, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// Convert the CSV JSON string to string array
func (a *CSVJSONArray) UnmarshalCSV(csv string) error {
	err := json.Unmarshal([]byte(csv), &a)
	return err
}

func (a *CSVJSONArray) String() []string {
	array := make([]string, len(*a))
	for _, str := range *a {
		array = append(array, str)
	}
	return array
}
