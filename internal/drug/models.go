package drug

// Recall is a single openFDA drug enforcement report
// (https://open.fda.gov/apis/drug/enforcement/).
type Recall struct {
	RecallNumber            string `json:"recall_number"`
	ReasonForRecall         string `json:"reason_for_recall"`
	Status                  string `json:"status"`
	City                    string `json:"city"`
	State                   string `json:"state"`
	Country                 string `json:"country"`
	Classification          string `json:"classification"`
	EventID                 string `json:"event_id"`
	RecallingFirm           string `json:"recalling_firm"`
	ReportDate              string `json:"report_date"`
	VoluntaryMandated       string `json:"voluntary_mandated"`
	InitialFirmNotification string `json:"initial_firm_notification"`
	DistributionPattern     string `json:"distribution_pattern"`
	ProductDescription      string `json:"product_description"`
	ProductQuantity         string `json:"product_quantity"`
	CodeInfo                string `json:"code_info"`
	TerminationDate         string `json:"termination_date"`
	RecallInitiationDate    string `json:"recall_initiation_date"`
}

// RecallResultsMeta is the pagination/count block of an openFDA response.
type RecallResultsMeta struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// RecallMeta is the meta block of an openFDA response.
type RecallMeta struct {
	Disclaimer  string            `json:"disclaimer"`
	Terms       string            `json:"terms"`
	License     string            `json:"license"`
	LastUpdated string            `json:"last_updated"`
	Results     RecallResultsMeta `json:"results"`
}

// RecallResponse is the openFDA drug enforcement API response envelope.
type RecallResponse struct {
	Meta    RecallMeta `json:"meta"`
	Results []Recall   `json:"results"`
}
