package rpc

func (db *db) InsertSMS(p SMS) (*SMS, error) {
	var sms SMS
	err := db.postgres.QueryRow(bg(), `INSERT INTO
		sms(account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status)
		values($1, '', NOW(), NOW(), $2, $3, $4, $5, $6, $7, $8, 'pending')
		RETURNING id, account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status
		`,
		p.AccountID,
		p.MessageRef,
		p.Country,
		p.Message,
		p.SMSCount,
		p.GSM,
		p.Recipient,
		p.Sender,
	).Scan(&sms.ID, &sms.AccountID, &sms.MessageID, &sms.CreatedAt, &sms.UpdatedAt, &sms.MessageRef, &sms.Country, &sms.Message, &sms.SMSCount, &sms.GSM, &sms.Recipient, &sms.Sender, &sms.Status)
	if err != nil {
		return &SMS{}, err
	}

	return &sms, nil

}

func (db *db) MarkStatus(smsID string, status string) error {
	sql := `
	UPDATE sms
	SET status = $2
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
	SET status = 'sent', message_id = $2
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
	SET status = 'failed'
	WHERE id = $1
	`
	_, err := db.Exec(sql, smsID)
	if err != nil {
		return err
	}
	return nil
}

func (db *db) FindSMSByMessageID(messageID string) (*SMS, error) {
	var sms SMS
	err := db.postgres.QueryRow(bg(), `SELECT id, account_id, message_id, created_at, updated_at, message_ref, country, message, sms_count, gsm, recipient, sender, status
		FROM sms
		WHERE message_id = $1
		`,
		messageID,
	).Scan(&sms.ID, &sms.AccountID, &sms.MessageID, &sms.CreatedAt, &sms.UpdatedAt, &sms.MessageRef, &sms.Country, &sms.Message, &sms.SMSCount, &sms.GSM, &sms.Recipient, &sms.Sender, &sms.Status)
	if err != nil {
		return &SMS{}, err
	}

	return &sms, nil
}
