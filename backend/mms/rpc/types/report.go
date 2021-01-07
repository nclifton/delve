package types

type ReportLine struct {
	CurrentMonth  int64
	PreviousMonth int64
	TwoMonthsAgo  int64
}

type AccountMMSUsage map[string]*ReportLine

type GenerateAccountMMSUsageParams struct {
	AccountID string
}

type GenerateAccountMMSUsageResponse struct {
	Report AccountMMSUsage
}
