package service

import (
	"jobsheet-go-aws2/database/business"
	"jobsheet-go-aws2/database/business_system"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/dto"
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetBusinessList(c echo.Context) error {
	businessList, err := business.Scan()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var restBusinessList []dto.RestBusiness
	for _, business := range businessList {
		restBusinessList = append(restBusinessList, dto.NewRestBusiness(business))
	}
	// IDの昇順にソートする。
	sort.Slice(restBusinessList, func(i, j int) bool {
		return restBusinessList[i].ID < restBusinessList[j].ID
	})
	return c.JSON(http.StatusCreated, restBusinessList)
}

func RegistBusiness(c echo.Context) error {
	restBusiness := new(dto.RestBusiness)
	err := c.Bind(restBusiness)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetBusiness := new(model.Business)
	targetBusiness.ID = restBusiness.ID
	targetBusiness.Name = restBusiness.Name
	err = business.PutItem(*targetBusiness)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "regist business")
}

func DeleteBusiness(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	// システムに紐づいているか確認
	businessSystemList, err := business_system.Search(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var resultFlg string
	if len(businessSystemList) > 0 {
		// 紐づいているシステムがある場合は削除しない。
		resultFlg = "1"
	} else {
		// 紐づいているシステムがない場合は削除する。
		err = business.DeleteItem(id)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		resultFlg = "0"
	}
	return c.String(http.StatusOK, resultFlg)
}
