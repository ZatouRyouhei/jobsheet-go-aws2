package dto

import "jobsheet-go-aws2/database/model"

type RestInquiry struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRestInquiry(inquiry model.Inquiry) RestInquiry {
	restInquiry := new(RestInquiry)
	restInquiry.ID = inquiry.ID
	restInquiry.Name = inquiry.Name
	return *restInquiry
}
