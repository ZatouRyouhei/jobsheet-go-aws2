package service

import (
	"jobsheet-go-aws2/database/client"
	"jobsheet-go-aws2/database/jobsheet"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/dto"
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetClientList(c echo.Context) error {
	clientList, err := client.Scan()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	var restClientList []dto.RestClient
	for _, client := range clientList {
		restClientList = append(restClientList, dto.NewRestClient(client))
	}
	// IDの昇順にソートする。
	sort.Slice(restClientList, func(i, j int) bool {
		return restClientList[i].ID < restClientList[j].ID
	})
	return c.JSON(http.StatusCreated, restClientList)
}

func RegistClient(c echo.Context) error {
	restClient := new(dto.RestClient)
	err := c.Bind(restClient)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetClient := new(model.Client)
	targetClient.ID = restClient.ID
	targetClient.Name = restClient.Name
	err = client.PutItem(*targetClient)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "regist client")
}

func DeleteClient(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	resultFlg := "0"
	// 業務日誌で使用されているか確認
	clientList, err := jobsheet.SearchForClient(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	if len(clientList) > 0 {
		// 使用されている場合は削除しない
		resultFlg = "1"
	} else {
		// 使用されていない場合は削除する。
		err = client.DeleteItem(id)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	return c.String(http.StatusOK, resultFlg)
}
