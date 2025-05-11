package dto

import "jobsheet-go-aws2/database/model"

type RestClient struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRestClient(client model.Client) RestClient {
	restClient := new(RestClient)
	restClient.ID = client.ID
	restClient.Name = client.Name
	return *restClient
}
