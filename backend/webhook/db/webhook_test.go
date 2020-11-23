package db_test

import (
	"fmt"
	"log"
	"testing"
)

var (
	testAccountID = "12345"
)

func TestInsert(t *testing.T) {
	testUtil.Setup()
	t.Run("can insert a webhook", func(t *testing.T) {
		w, err := webhookDB.Insert(testAccountID,
			"sms_inbound",
			"my first webby",
			"https://google.com",
			5,
		)
		if err != nil {
			t.Fatalf("failed to insert a webhook: %s", err)
		}

		if w.Event != "sms_inbound" {
			t.Fatalf("failed to insert a webhook - expected 'sms_inbound' event, got: '%s'", w.Event)
		}
	})
}

func TestFind(t *testing.T) {
	testUtil.Setup()
	_, err := testUtil.Postgres().Exec(bg(), fmt.Sprintf(`
	INSERT INTO webhook (account_id, event, name, url, rate_limit, created_at, updated_at)
		VALUES
			(%[1]s, 'sms_inbound', 'webhook the first', 'https://google.com', 5, current_timestamp, current_timestamp),
			(%[1]s, 'sms_inbound', 'webhook the second', 'https://google.com', 6, current_timestamp, current_timestamp),
			(%[1]s, 'sms_inbound', 'webhook the third', 'https://google.com', 7, current_timestamp, current_timestamp);
	`, testAccountID))
	if err != nil {
		log.Fatal(err)
	}
	t.Run("can retrieve webhooks", func(t *testing.T) {
		whs, err := webhookDB.Find(testAccountID)
		if err != nil {
			t.Fatalf("failed to find webhooks by account id: %s", err)
		}
		if len(whs) != 3 {
			t.Fatalf("expected 3 webhooks, got: %d", len(whs))
		}
	})
}

func TestDelete(t *testing.T) {
	testUtil.Setup()
	_, err := testUtil.Postgres().Exec(bg(), fmt.Sprintf(`
	INSERT INTO webhook (account_id, event, name, url, rate_limit, created_at, updated_at)
		VALUES
			(%[1]s, 'sms_inbound', 'webhook the first', 'https://google.com', 5, current_timestamp, current_timestamp),
			(%[1]s, 'sms_inbound', 'webhook the second', 'https://google.com', 6, current_timestamp, current_timestamp),
			(%[1]s, 'sms_inbound', 'webhook the third', 'https://google.com', 7, current_timestamp, current_timestamp);
	`, testAccountID))
	if err != nil {
		log.Fatal(err)
	}
	t.Run("can delete webhooks by id", func(t *testing.T) {
		err := webhookDB.Delete(testAccountID, "1")

		if err != nil {
			t.Fatalf("failed to delete webhooks: %s", err)
		}
	})
	t.Run("can't delete non-existant webhooks", func(t *testing.T) {
		err := webhookDB.Delete(testAccountID, "5678")

		if err == nil {
			t.Fatalf("expected delete webhook error")
		}
	})
}

func TestUpdate(t *testing.T) {
	testUtil.Setup()
	_, err := testUtil.Postgres().Exec(bg(), fmt.Sprintf(`
	INSERT INTO webhook (account_id, event, name, url, rate_limit, created_at, updated_at)
		VALUES
			(%[1]s, 'sms_inbound', 'webhook the first', 'https://google.com', 5, current_timestamp, current_timestamp);
	`, testAccountID))
	if err != nil {
		log.Fatal(err)
	}

	t.Run("can update a webhook", func(t *testing.T) {
		w, err := webhookDB.Update("1", testAccountID,
			"link_hit",
			"my first webby",
			"https://google.com",
			5,
		)
		if err != nil {
			t.Fatalf("failed to update a webhook: %s", err)
		}

		if w.Event != "link_hit" {
			t.Fatalf("failed to insert a webhook - expected 'link_hit' event, got: '%s'", w.Event)
		}
	})

	t.Run("can not update not existant webhook", func(t *testing.T) {
		_, err := webhookDB.Update("987", testAccountID,
			"link_hit",
			"my first webby",
			"https://google.com",
			5,
		)
		if err == nil {
			t.Fatalf("expected update webhook error")
		}
	})
}

func TestFindByEvent(t *testing.T) {
	testUtil.Setup()
	_, err := testUtil.Postgres().Exec(bg(), fmt.Sprintf(`
	INSERT INTO webhook (account_id, event, name, url, rate_limit, created_at, updated_at)
		VALUES
			(%[1]s, 'sms_inbound', 'webhook the first', 'https://google.com', 5, current_timestamp, current_timestamp),
			(%[1]s, 'sms_inbound', 'webhook the second', 'https://google.com', 6, current_timestamp, current_timestamp),
			(%[1]s, 'contact_update', 'webhook the third', 'https://google.com', 7, current_timestamp, current_timestamp);
	`, testAccountID))
	if err != nil {
		log.Fatal(err)
	}
	t.Run("can retrieve webhooks", func(t *testing.T) {
		whs, err := webhookDB.FindByEvent(testAccountID, "sms_inbound")
		if err != nil {
			t.Fatalf("failed to find webhooks by account id: %s", err)
		}
		if len(whs) != 2 {
			t.Fatalf("expected 2 webhooks, got: %d", len(whs))
		}
	})
}
