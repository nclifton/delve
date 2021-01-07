package rpc

type SMSUsage struct {
	AccountID string
	Total     int
	Month     int
	MonthName string
	Year      int
}

func (db *db) GenerateAccountSMSUsage(AccountID string) ([]SMSUsage, error) {

	usage := []SMSUsage{}

	// Sum the last 3 months of SMS totals
	rows, err := db.postgres.Query(
		bg(),
		`select account_id, total, month, year from (
				select account_id,
					sum(sms_count) as total,
					extract(month from created_at) as month,
					extract(year from created_at) as year
        from sms
        where status not in ('failed', 'UNKNOWN')
					and CASE WHEN $1 = '' THEN account_id is not null ELSE account_id = $1::uuid END
        group by account_id, month, year
    ) as t where month > (extract(month from NOW())-3)`,
		AccountID,
	)
	if err != nil {
		return usage, err
	}
	defer rows.Close()

	for rows.Next() {
		sms := SMSUsage{}
		err := rows.Scan(
			&sms.AccountID,
			&sms.Total,
			&sms.Month,
			&sms.Year,
		)
		if err != nil {
			return usage, err
		}
		usage = append(usage, sms)
	}

	return usage, nil
}
