package rpc

import (
	"log"
	"time"

	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
)

func (s *SMSService) GenerateAccountSMSUsage(p types.GenerateAccountSMSUsageParams, r *types.GenerateAccountSMSUsageResponse) error {

	// Get the usage numbers
	rep, err := s.db.GenerateAccountSMSUsage(p.AccountID)
	if err != nil {
		log.Printf("[Generate SMS Usage] Could not generate total sms counts: %s", p.AccountID)
		return err
	}

	// Generate the report format
	currentMonth := int(time.Now().Month())
	report := make(types.AccountSMSUsage)

	for _, sms := range rep {
		monthdiff := currentMonth - sms.Month
		if monthdiff < 0 {
			monthdiff += 12
		}
		if _, found := report[sms.AccountID]; !found {
			report[sms.AccountID] = &types.ReportLine{}
		}
		switch monthdiff {
		case 0:
			report[sms.AccountID].CurrentMonth = int64(sms.Total)
		case 1:
			report[sms.AccountID].PreviousMonth = int64(sms.Total)
		case 2:
			report[sms.AccountID].TwoMonthsAgo = int64(sms.Total)
		}

	}

	r.Report = report

	return nil
}
