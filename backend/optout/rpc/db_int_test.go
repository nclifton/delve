// +build integration

package rpc

import (
	"context"
	"fmt"
	"os"
	"testing"

	types "github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
)

func initDB(t *testing.T) *db {
	db, err := NewDB(os.Getenv("OPTOUT_TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func Test_DB_FindOptOutByLinkID(t *testing.T) {
	db := initDB(t)

	testOptOuts := []types.OptOut{
		types.OptOut{ID: "11111111-1111-1111-1111-111111111111", AccountID: "11111111-1111-1111-1111-111111111112", MessageID: "11111111-1111-1111-1111-111111111113", MessageType: "SMS", LinkID: "dVYHEhq6", Sender: "1701"},
	}

	tests := []struct {
		name           string
		linkID         string
		expectedOptOut types.OptOut
		expectedErr    error
	}{
		{
			name:           "test happy path",
			linkID:         "dVYHEhq6",
			expectedOptOut: testOptOuts[0],
			expectedErr:    nil,
		},
		{
			name:           "test link id not found",
			linkID:         "00000000",
			expectedOptOut: types.OptOut{},
			expectedErr:    fmt.Errorf("optOut not found"),
		},
	}

	for _, test := range tests {
		cleanup(t, db)
		insertOptOuts(t, db, testOptOuts)

		t.Run(test.name, func(t *testing.T) {
			res, err := db.FindOptOutByLinkID(context.Background(), test.linkID)
			if err == nil && !types.OptOutEqual(test.expectedOptOut, *res) {
				t.Errorf("unexpected result %+v, \nbut got %+v", test.expectedOptOut, res)
			}

			if err == nil && test.expectedOptOut.LinkID != res.LinkID {
				t.Errorf("unexpected result %+v, \nbut got %+v", test.expectedOptOut, res)
			}

			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("unexpected error, got error %v", err)
			}
		})
	}
}

func Test_DB_InsertOptOut(t *testing.T) {
	db := initDB(t)

	testOptOuts := []types.OptOut{
		types.OptOut{ID: "11111111-1111-1111-1111-111111111111", AccountID: "11111111-1111-1111-1111-111111111112", MessageID: "11111111-1111-1111-1111-111111111113", MessageType: "SMS", LinkID: "dVYHEhq6", Sender: "1701"},
	}

	tests := []struct {
		name           string
		accountID      string
		messageID      string
		messageType    string
		sender         string
		expectedOptOut types.OptOut
		expectedErr    error
	}{
		{
			name:           "test happy path",
			accountID:      "11111111-1111-1111-1111-111111111122",
			messageID:      "11111111-1111-1111-1111-111111111133",
			messageType:    "SMS",
			sender:         "1701",
			expectedOptOut: types.OptOut{AccountID: "11111111-1111-1111-1111-111111111122", MessageID: "11111111-1111-1111-1111-111111111133", MessageType: "SMS", Sender: "1701"},
			expectedErr:    nil,
		},
		{
			name:        "test with error - empty value",
			accountID:   "",
			messageID:   "",
			messageType: "",
			expectedErr: fmt.Errorf("ERROR: invalid input syntax for type uuid: \"\" (SQLSTATE 22P02)"),
		},
	}

	for _, test := range tests {
		cleanup(t, db)
		insertOptOuts(t, db, testOptOuts)

		t.Run(test.name, func(t *testing.T) {
			res, err := db.InsertOptOut(context.Background(), test.accountID, test.messageID, test.messageType, test.sender)
			if err == nil && !types.OptOutEqual(test.expectedOptOut, *res) {
				t.Errorf("unexpected result %+v, \nbut got %+v", test.expectedOptOut, res)
			}

			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("unexpected error, got error %v", err)
			}
		})
	}
}

///////////////////
// insert & cleanup
///////////////////

func insertOptOuts(t *testing.T, db *db, optOuts []types.OptOut) {
	for _, optOut := range optOuts {
		if _, err := db.postgres.Exec(context.Background(),
			`INSERT INTO opt_out(id, account_id, message_id, message_type, sender, link_id, created_at, updated_at)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
			optOut.ID,
			optOut.AccountID,
			optOut.MessageID,
			optOut.MessageType,
			optOut.Sender,
			optOut.LinkID,
			optOut.CreatedAt,
			optOut.UpdatedAt,
		); err != nil {
			t.Fatal(err)
		}
	}
}

func cleanup(t *testing.T, db *db) {
	if _, err := db.postgres.Exec(context.Background(), "TRUNCATE TABLE opt_out CASCADE"); err != nil {
		t.Fatal(err)
	}
}
