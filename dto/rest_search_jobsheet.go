package dto

import (
	"jobsheet-go-aws2/database/business"
	"jobsheet-go-aws2/database/business_system"
	"jobsheet-go-aws2/database/client"
	"jobsheet-go-aws2/database/inquiry"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/database/user"
	"log/slog"
)

type RestSearchJobSheet struct {
	ID             string             `json:"id"`
	Client         RestClient         `json:"client"`
	BusinessSystem RestBusinessSystem `json:"businessSystem"`
	Inquiry        RestInquiry        `json:"inquiry"`
	Department     string             `json:"department"`
	Person         string             `json:"person"`
	OccurDate      string             `json:"occurDate"`
	OccurTime      string             `json:"occurTime"`
	Title          string             `json:"title"`
	Content        string             `json:"content"`
	Contact        RestUser           `json:"contact"`
	LimitDate      string             `json:"limitDate"`
	Deal           RestUser           `json:"deal"`
	CompleteDate   string             `json:"completeDate"`
	Support        string             `json:"support"`
	ResponseTime   float64            `json:"responseTime"`
	FileList       []RestAttachment   `json:"fileList"`
}

func NewRestSearchJobSheet(jobSheet model.JobSheet, attachmentList []model.Attachment) RestSearchJobSheet {
	restSearchJobSheet := new(RestSearchJobSheet)
	restSearchJobSheet.ID = jobSheet.ID
	targetClient, err := client.GetItem(jobSheet.ClientID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	} else {
		restSearchJobSheet.Client = NewRestClient(*targetClient)
	}
	targetBusinessSystem, err := business_system.GetItem(jobSheet.BusinessSystemID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	} else {
		targetBusiness, err := business.GetItem(targetBusinessSystem.BusinessID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
		} else {
			restSearchJobSheet.BusinessSystem = RestBusinessSystem{
				ID:   targetBusinessSystem.ID,
				Name: targetBusinessSystem.Name,
				Business: RestBusiness{
					ID:   targetBusiness.ID,
					Name: targetBusiness.Name,
				},
			}
		}
	}
	targetInquiry, err := inquiry.GetItem(jobSheet.InquiryID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	} else {
		restSearchJobSheet.Inquiry = NewRestInquiry(*targetInquiry)
	}
	restSearchJobSheet.Department = jobSheet.Department
	restSearchJobSheet.Person = jobSheet.Person
	restSearchJobSheet.OccurDate = jobSheet.OccurDate
	restSearchJobSheet.OccurTime = jobSheet.OccurTime
	restSearchJobSheet.Title = jobSheet.Title
	restSearchJobSheet.Content = jobSheet.Content
	targetContact, err := user.GetItem(jobSheet.ContactID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	} else {
		restSearchJobSheet.Contact = RestUser{
			Id:       targetContact.Id,
			Password: "",
			Name:     targetContact.Name,
			SeqNo:    targetContact.SeqNo,
		}
	}
	restSearchJobSheet.LimitDate = jobSheet.LimitDate
	targetDeal, err := user.GetItem(jobSheet.DealID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	} else {
		restSearchJobSheet.Deal = RestUser{
			Id:       targetDeal.Id,
			Password: "",
			Name:     targetDeal.Name,
			SeqNo:    targetDeal.SeqNo,
		}
	}
	restSearchJobSheet.CompleteDate = jobSheet.CompleteDate
	restSearchJobSheet.Support = jobSheet.Support
	restSearchJobSheet.ResponseTime = jobSheet.ResponseTime
	restSearchJobSheet.FileList = []RestAttachment{}
	for _, attachment := range attachmentList {
		restSearchJobSheet.FileList = append(restSearchJobSheet.FileList, RestAttachment{SeqNo: attachment.SeqNo, FileName: attachment.FileName})
	}
	return *restSearchJobSheet
}
