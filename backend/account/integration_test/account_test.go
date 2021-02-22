package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/accountpb"
	"github.com/burstsms/mtmo-tp/backend/account/rpc/db"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_FindAccountByAPIKey(t *testing.T) {
	s := newSetup(t, tfx)

	time1 := time.Date(2021, 02, 10, 10, 30, 22, 0, time.UTC)

	testAccounts := []db.Account{
		{ID: "11111111-1111-1111-1111-111111111111", CreatedAt: time1, UpdatedAt: time1, Name: "test account 1", AlarisUsername: "user1", AlarisPassword: "pwd", AlarisURL: "http://alaris.com"},
	}

	testAccountAPIKeys := []db.AccountAPIKey{
		{ID: "11111111-1111-1111-1111-111111111123", AccountID: "11111111-1111-1111-1111-111111111111", Description: "123", Key: "key-123"},
	}

	tests := []struct {
		name              string
		key               string
		expectedAccountPB *accountpb.Account
		expectedErr       error
	}{
		{
			name: "happy path",
			key:  "key-123",
			expectedAccountPB: &accountpb.Account{
				Id:             testAccounts[0].ID,
				CreatedAt:      timestamppb.New(time1),
				UpdatedAt:      timestamppb.New(time1),
				Name:           testAccounts[0].Name,
				AlarisUsername: testAccounts[0].AlarisUsername,
				AlarisPassword: testAccounts[0].AlarisPassword,
				AlarisUrl:      testAccounts[0].AlarisURL,
			},
			expectedErr: nil,
		},
		{
			name:        "test with error - account not found",
			key:         "randomkey",
			expectedErr: errors.New("rpc error: code = NotFound desc = account not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup(t, s.DBConn)
			insertAccounts(t, s.DBConn, testAccounts)
			insertAccountAPIKeys(t, s.DBConn, testAccountAPIKeys)

			r, err := s.accountClient.FindAccountByAPIKey(context.Background(), &accountpb.FindAccountByAPIKeyParams{
				Key: test.key,
			})

			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected error %+v, \nbut got %+v", test.expectedErr, err)
			}

			if r != nil && !AccountPBEqual(test.expectedAccountPB, r.Account) {
				t.Errorf("expected result %+v, \nbut got %+v", test.expectedAccountPB, r.Account)
			}
		})
	}
}

func Test_FindAccountByID(t *testing.T) {
	s := newSetup(t, tfx)

	time1 := time.Date(2021, 02, 10, 10, 30, 22, 0, time.UTC)

	testAccounts := []db.Account{
		{ID: "11111111-1111-1111-1111-111111111111", CreatedAt: time1, UpdatedAt: time1, Name: "test account 1", AlarisUsername: "user1", AlarisPassword: "pwd", AlarisURL: "http://alaris.com"},
	}

	tests := []struct {
		name              string
		id                string
		expectedAccountPB *accountpb.Account
		expectedErr       error
	}{
		{
			name: "happy path",
			id:   "11111111-1111-1111-1111-111111111111",
			expectedAccountPB: &accountpb.Account{
				Id:             testAccounts[0].ID,
				CreatedAt:      timestamppb.New(time1),
				UpdatedAt:      timestamppb.New(time1),
				Name:           testAccounts[0].Name,
				AlarisUsername: testAccounts[0].AlarisUsername,
				AlarisPassword: testAccounts[0].AlarisPassword,
				AlarisUrl:      testAccounts[0].AlarisURL,
			},
			expectedErr: nil,
		},
		{
			name:        "test with error - account not found",
			id:          "11111111-1111-1111-1111-111111111789",
			expectedErr: errors.New("rpc error: code = NotFound desc = account not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup(t, s.DBConn)
			insertAccounts(t, s.DBConn, testAccounts)

			r, err := s.accountClient.FindAccountByID(context.Background(), &accountpb.FindAccountByIDParams{
				Id: test.id,
			})

			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected error %+v, \nbut got %+v", test.expectedErr, err)
			}

			if r != nil && !AccountPBEqual(test.expectedAccountPB, r.Account) {
				t.Errorf("expected result %+v, \nbut got %+v", test.expectedAccountPB, r.Account)
			}
		})
	}
}
