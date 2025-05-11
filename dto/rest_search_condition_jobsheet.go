package dto

type RestSearchConditionJobSheet struct {
	Client         int    `json:"client"`
	Business       int    `json:"business"`
	BusinessSystem int    `json:"businessSystem"`
	Inquiry        int    `json:"inquiry"`
	Contact        string `json:"contact"`
	Deal           string `json:"deal"`
	OccurDateFrom  string `json:"occurDateFrom"`
	OccurDateTo    string `json:"occurDateTo"`
	CompleteSign   int    `json:"completeSign"`
	LimitDate      string `json:"limitDate"`
	Keyword        string `json:"keyword"`
}
