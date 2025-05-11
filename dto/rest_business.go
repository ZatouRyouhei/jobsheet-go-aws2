package dto

import (
	"jobsheet-go-aws2/database/model"
)

type RestBusiness struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRestBusiness(business model.Business) RestBusiness {
	return RestBusiness{
		ID:   business.ID,
		Name: business.Name,
	}
}
