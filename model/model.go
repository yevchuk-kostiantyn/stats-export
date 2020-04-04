package model

type University struct {
	Code   float64
	Name   string
	Number float64
	USIDs  []string
}

type EDBOOfferResponse struct {
	Offers             []interface{} `json:"offers"`
	OffersSubjects     interface{}   `json:"offers_subjects"`
	OffersRequestsInfo interface{}   `json:"offers_requests_info"`
	Facultets          interface{}   `json:"facultets"`
	Specialities       interface{}   `json:"specialities"`
	Subjects           interface{}   `json:"subjects"`
	Specializations    interface{}   `json:"specializations"`
}

type ExportResponse struct {
	Region         string  `json:"region,omitempty"`
	Speciality     string  `json:"speciality,omitempty"`
	LicenseVolume  float64 `json:"license_volume,omitempty"`
	BudgetVolume   float64 `json:"budget_volume,omitempty"`
	ContractVolume float64 `json:"contract_volume,omitempty"`
}
