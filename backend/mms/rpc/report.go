package rpc

import (
	"log"
	"time"

	"github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
)

func (s *MMSService) GenerateAccountMMSUsage(p types.GenerateAccountMMSUsageParams, r *types.GenerateAccountMMSUsageResponse) error {

	rep, err := s.db.GenerateAccountMMSUsage(p.AccountID)
	if err != nil {
		log.Printf("[Generate MMS Usage] Could not generate total mms counts: %s", p.AccountID)
		return err
	}
	currentMonth := int(time.Now().Month())
	report := make(types.AccountMMSUsage)
	for _, mms := range rep {

		monthdiff := currentMonth - mms.Month
		if monthdiff < 0 {
			monthdiff += 12
		}
		if _, found := report[mms.AccountID]; !found {
			report[mms.AccountID] = &types.ReportLine{}
		}
		switch monthdiff {
		case 0:
			report[mms.AccountID].CurrentMonth = int64(mms.Total)
		case 1:
			report[mms.AccountID].PreviousMonth = int64(mms.Total)
		case 2:
			report[mms.AccountID].TwoMonthsAgo = int64(mms.Total)
		}
	}

	r.Report = report

	return nil
}
