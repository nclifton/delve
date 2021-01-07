package adminapi

import (
	"log"
	"net/http"
	"time"

	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
)

type ReportLine struct {
	AccountID            string `json:"account_id"`
	SMSCurrentMonth      int64  `json:"sms_current_month"`
	SMSCurrentMonthName  string `json:"sms_current_month_name"`
	SMSPreviousMonth     int64  `json:"sms_previous_month"`
	SMSPreviousMonthName string `json:"sms_previous_month_name"`
	SMSTwoMonthsAgo      int64  `json:"sms_two_months_ago"`
	SMSTwoMonthsAgoName  string `json:"sms_two_months_ago_name"`
	MMSCurrentMonth      int64  `json:"mms_current_month"`
	MMSCurrentMonthName  string `json:"mms_current_month_name"`
	MMSPreviousMonth     int64  `json:"mms_previous_month"`
	MMSPreviousMonthName string `json:"mms_previous_month_name"`
	MMSTwoMonthsAgo      int64  `json:"mms_two_months_ago"`
	MMSTwoMonthsAgoName  string `json:"mms_two_months_ago_name"`
}

func UsageReportGET(r *Route) {

	type report map[string]*ReportLine
	data := make(report)

	now := time.Now().UTC()
	_, currentMonth, _ := now.Date()
	_, previousMonth, _ := time.Now().AddDate(0, -1, 0).Date()
	_, twoMonthsAgo, _ := time.Now().AddDate(0, -2, 0).Date()

	accountID := r.params.ByName("account_id")

	res, err := r.api.sms.GenerateAccountSMSUsage(sms.GenerateAccountSMSUsageParams{AccountID: accountID})
	if err != nil {
		// handler rpc error
		log.Printf("Could not generate SMS totals: %s", err.Error())
		r.WriteError("Could not generate sms totals", http.StatusInternalServerError)
		return
	}

	for k, v := range res.Report {
		if _, found := data[k]; !found {
			data[k] = &ReportLine{AccountID: k}
		}
		data[k].SMSCurrentMonth = v.CurrentMonth
		data[k].SMSCurrentMonthName = currentMonth.String()
		data[k].SMSPreviousMonth = v.PreviousMonth
		data[k].SMSPreviousMonthName = previousMonth.String()
		data[k].SMSTwoMonthsAgo = v.TwoMonthsAgo
		data[k].SMSTwoMonthsAgoName = twoMonthsAgo.String()
	}

	resMMS, err := r.api.mms.GenerateAccountMMSUsage(mms.GenerateAccountMMSUsageParams{AccountID: accountID})
	if err != nil {
		// handler rpc error
		log.Printf("Could not generate MMS totals: %s", err.Error())
		r.WriteError("Could not generate mms totals", http.StatusInternalServerError)
		return
	}

	for k, v := range resMMS.Report {
		if _, found := data[k]; !found {
			data[k] = &ReportLine{AccountID: k}
		}
		data[k].MMSCurrentMonth = v.CurrentMonth
		data[k].MMSCurrentMonthName = currentMonth.String()
		data[k].MMSPreviousMonth = v.PreviousMonth
		data[k].MMSPreviousMonthName = previousMonth.String()
		data[k].MMSTwoMonthsAgo = v.TwoMonthsAgo
		data[k].MMSTwoMonthsAgoName = twoMonthsAgo.String()

	}

	type line map[string]interface{}
	flattened := []line{}
	for _, v := range data {
		line := make(line)
		// get the account id
		accountName := "unknown"
		res, err := r.api.account.FindByID(v.AccountID)
		if err != nil {
			log.Printf("Could not find associated account: %s", err.Error())
		}
		if res.Account.Name != "" {
			accountName = res.Account.Name
		}
		line["account_id"] = v.AccountID
		line["account_name"] = accountName
		line["country"] = "au"
		line["sms-"+v.SMSCurrentMonthName] = v.SMSCurrentMonth
		line["sms-"+v.SMSPreviousMonthName] = v.SMSPreviousMonth
		line["sms-"+v.SMSTwoMonthsAgoName] = v.SMSTwoMonthsAgo
		line["mms-"+v.MMSCurrentMonthName] = v.MMSCurrentMonth
		line["mms-"+v.MMSPreviousMonthName] = v.MMSPreviousMonth
		line["mms-"+v.MMSTwoMonthsAgoName] = v.MMSTwoMonthsAgo
		flattened = append(flattened, line)
	}

	r.Write(flattened, http.StatusOK)
}
