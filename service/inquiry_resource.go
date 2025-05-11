package service

import (
	"jobsheet-go-aws2/database/inquiry"
	"jobsheet-go-aws2/dto"
	"log/slog"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

func GetInquiryList(c echo.Context) error {
	inquiryList, err := inquiry.Scan()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var restInquiryList []dto.RestInquiry
	for _, inquiry := range inquiryList {
		restInquiryList = append(restInquiryList, dto.NewRestInquiry(inquiry))
	}
	// IDの昇順にソートする。
	sort.Slice(restInquiryList, func(i, j int) bool {
		return restInquiryList[i].ID < restInquiryList[j].ID
	})
	return c.JSON(http.StatusCreated, restInquiryList)
}
