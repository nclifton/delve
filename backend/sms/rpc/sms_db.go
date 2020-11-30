package rpc

import (
	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	"github.com/jackc/pgx/v4"
)

func (db *db) InsertSMS(p types.SMS) (*types.SMS, error) {
	var sms types.SMS
	err := db.postgres.QueryRow(bg(), `INSERT INTO
		sms(id, account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status, track_links)
		values($1, $2, '', NOW(), NOW(), $3, $4, $5, $6, $7, $8, $9, 'pending', $10)
		RETURNING id, account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status, track_links
		`,
		p.ID,
		p.AccountID,
		p.MessageRef,
		p.Country,
		p.Message,
		p.SMSCount,
		p.GSM,
		p.Recipient,
		p.Sender,
		p.TrackLinks,
	).Scan(&sms.ID, &sms.AccountID, &sms.MessageID, &sms.CreatedAt, &sms.UpdatedAt, &sms.MessageRef, &sms.Country, &sms.Message, &sms.SMSCount, &sms.GSM, &sms.Recipient, &sms.Sender, &sms.Status, &sms.TrackLinks)
	if err != nil {
		return &types.SMS{}, err
	}

	return &sms, nil

}

func (db *db) MarkStatus(smsID string, status string) error {
	sql := `
	UPDATE sms
	SET status = $2, updated_at = NOW()
	WHERE id = $1
	`
	_, err := db.Exec(sql, smsID, status)
	if err != nil {
		return err
	}

	return nil
}

func (db *db) MarkSent(smsID string, messageID string) error {
	sql := `
	UPDATE sms
	SET status = 'sent', message_id = $2, updated_at = NOW()
	WHERE id = $1
	`
	_, err := db.Exec(sql, smsID, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (db *db) MarkFailed(smsID string) error {
	sql := `
	UPDATE sms
	SET status = 'failed', updated_at = NOW()
	WHERE id = $1
	`
	_, err := db.Exec(sql, smsID)
	if err != nil {
		return err
	}
	return nil
}

func (db *db) FindSMSByID(ID string, AccountID string) (*types.SMS, error) {
	var sms types.SMS
	err := db.postgres.QueryRow(bg(), `
		SELECT id, account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status, track_links
		FROM sms
		WHERE id = $1 AND account_id = $2`,
		ID,
		AccountID,
	).Scan(&sms.ID, &sms.AccountID, &sms.MessageID, &sms.CreatedAt, &sms.UpdatedAt, &sms.MessageRef, &sms.Country, &sms.Message, &sms.SMSCount, &sms.GSM, &sms.Recipient, &sms.Sender, &sms.Status, &sms.TrackLinks)
	if err != nil {
		return &types.SMS{}, err
	}

	return &sms, nil
}

func (db *db) FindSMSByMessageID(messageID string) (*types.SMS, error) {
	var sms types.SMS
	err := db.postgres.QueryRow(bg(), `
		SELECT id, account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status, track_links
		FROM sms
		WHERE message_id = $1`,
		messageID,
	).Scan(&sms.ID, &sms.AccountID, &sms.MessageID, &sms.CreatedAt, &sms.UpdatedAt, &sms.MessageRef, &sms.Country, &sms.Message, &sms.SMSCount, &sms.GSM, &sms.Recipient, &sms.Sender, &sms.Status, &sms.TrackLinks)
	if err != nil {
		return &types.SMS{}, err
	}

	return &sms, nil
}

func (db *db) FindSMSRelatedToMO(accountID string, mosender string, morecipient string) (*types.SMS, error) {
	var sms types.SMS
	err := db.postgres.QueryRow(bg(), `
		SELECT id, account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status, track_links
		FROM sms
		WHERE account_id = $1 AND sender = $3 AND recipient = $2
		AND created_at BETWEEN NOW() - INTERVAL '72 HOURS' AND NOW()
		ORDER BY updated_at DESC
		LIMIT 1
		`,
		accountID,
		mosender,
		morecipient,
	).Scan(&sms.ID, &sms.AccountID, &sms.MessageID, &sms.CreatedAt, &sms.UpdatedAt, &sms.MessageRef, &sms.Country, &sms.Message, &sms.SMSCount, &sms.GSM, &sms.Recipient, &sms.Sender, &sms.Status, &sms.TrackLinks)
	if err != nil && err != pgx.ErrNoRows {
		return &types.SMS{}, err
	}

	return &sms, nil
}
