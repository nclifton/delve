package service

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gocarina/gocsv"
	"github.com/vincent-petithory/dataurl"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type SenderCSV struct {
	AccountId      string       `csv:"account_id" validate:"omitempty,uuid"`
	Address        string       `csv:"address" validate:"required"`
	Country        string       `csv:"country" validate:"required"`
	Channels       CSVJSONArray `csv:"channels" validate:"required"` // see custom conversion below
	MMSProviderKey string       `csv:"mms_provider_key" validate:"omitempty"`
	Comment        string       `csv:"comment" validate:"omitempty"`
	Status         string       `csv:"status"`
	Error          string       `csv:"error"`
}

func (s *senderImpl) CreateSendersFromCSVDataURL(ctx context.Context, r *senderpb.CreateSendersFromCSVDataURLParams) (*senderpb.CreateSendersFromCSVDataURLReply, error) {

	replySenders := []*senderpb.Sender{}

	csvSenders, err := unmarshalSenderCSVDataUrl(r.CSV)
	if err != nil {
		return nil, err
	}

	// ignoring results here at the moment
	validCSVSenders, _, err := s.validateCSVSenders(ctx, csvSenders)
	if err != nil {
		return nil, err
	}

	if len(validCSVSenders) > 0 {

		newSenders := make([]db.Sender, 0, len(csvSenders))
		for _, sender := range csvSenders {
			newSenders = append(newSenders, db.Sender{
				AccountID:      sender.AccountId,
				Address:        sender.Address,
				MMSProviderKey: sender.MMSProviderKey,
				Channels:       sender.Channels,
				Country:        sender.Country,
				Comment:        sender.Comment,
			})
		}

		dbSenders, err := s.db.InsertSenders(ctx, newSenders)
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
