package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/gocarina/gocsv"
	"github.com/vincent-petithory/dataurl"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func (s *senderImpl) CreateSendersFromCSVDataURL(ctx context.Context, r *senderpb.CreateSendersFromCSVDataURLParams) (*senderpb.CreateSendersFromCSVDataURLReply, error) {

	replySenders := []*senderpb.Sender{}

	csvSenders, err := s.unmarshalSenderCSVDataUrl(r.CSV)
	if err != nil {
		return nil, err
	}

	if len(csvSenders) > 0 {

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

type SenderCSV struct {
	AccountId      string `csv:"account_id"`
	Address        string `csv:"address"`
	Country        string `csv:"country"`
	Channels       Array  `csv:"channels"`
	MMSProviderKey string `csv:"mms_provider_key"`
	Comment        string `csv:"comment"`
}

func (s *senderImpl) unmarshalSenderCSVDataUrl(csvDataUrl []byte) (csvSenders []SenderCSV, err error) {

	data, err := dataurl.DecodeString(string(csvDataUrl))
	if err != nil {
		return nil, err
	}

	csvSenders = []SenderCSV{}

	log.Printf("csvSenders: \n%s", string(data.Data))

	reader := gocsv.LazyCSVReader(bytes.NewReader(data.Data))
	err = gocsv.UnmarshalCSV(reader, &csvSenders)
	if err != nil {
		return nil, err
	}

	return csvSenders, nil

}

type Array []string

// Convert the internal string array to JSON string
func (a *Array) MarshalCSV() (string, error) {
	str, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// Convert the CSV JSON string to string array
func (a *Array) UnmarshalCSV(csv string) error {
	err := json.Unmarshal([]byte(csv), &a)
	return err
}

func (a *Array) String() []string {
	array := make([]string, len(*a))
	for _, str := range *a {
		array = append(array, str)
	}
	return array
}
