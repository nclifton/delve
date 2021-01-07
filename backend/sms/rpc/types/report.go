package types

type ReportLine struct {
	CurrentMonth  int64
	PreviousMonth int64
	TwoMonthsAgo  int64
}

type AccountSMSUsage map[string]*ReportLine

type GenerateAccountSMSUsageParams struct {
	AccountID string
}

type GenerateAccountSMSUsageResponse struct {
	Report AccountSMSUsage
}
