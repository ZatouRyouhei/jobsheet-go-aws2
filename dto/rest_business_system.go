package dto

import (
	"jobsheet-go-aws2/database/business"
	"jobsheet-go-aws2/database/model"
	"log/slog"
)

type RestBusinessSystem struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Business RestBusiness `json:"business"`
}

func NewRestBusinessSystem(businessSystem model.BusinessSystem) RestBusinessSystem {
	restBusinessSystem := new(RestBusinessSystem)
	restBusinessSystem.ID = businessSystem.ID
	restBusinessSystem.Name = businessSystem.Name
	targetBusiness, err := business.GetItem(businessSystem.BusinessID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	} else {
		restBusinessSystem.Business = NewRestBusiness(*targetBusiness)
	}
	return *restBusinessSystem
}
