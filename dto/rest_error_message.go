package dto

type RestErrorMessage struct {
	LineNo   int    `json:"lineNo"`
	ErrorMsg string `json:"errorMsg"`
}
