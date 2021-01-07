package rpc

type MMSUsage struct {
	AccountID string
	Total     int
	Month     int
	MonthName string
	Year      int
}

func (db *db) GenerateAccountMMSUsage(AccountID string) ([]MMSUsage, error) {

	usage := []MMSUsage{}
	// Sum the last 3 months of MMS totals
	rows, err := db.postgres.Query(
		bg(),
		`select account_id, total, month, year from (
				select account_id,
					count(id) as total,
					extract(month from created_at) as month,
					extract(year from created_at) as year
			from mms
			where status not in ('pending')
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
		mms := MMSUsage{}
		err := rows.Scan(
			&mms.AccountID,
			&mms.Total,
			&mms.Month,
			&mms.Year,
		)
		if err != nil {
			return usage, err
		}

		usage = append(usage, mms)

	}

	return usage, nil
}
