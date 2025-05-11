package service

import (
	"jobsheet-go-aws2/database/business_system"
	"jobsheet-go-aws2/database/jobsheet"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/dto"
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetBusinessSystemList(c echo.Context) error {
	id := c.Param("id")
	var businessSystemList []model.BusinessSystem
	var err error
	if id == "" {
		businessSystemList, err = business_system.Scan()
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	} else {
		paramId, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		businessSystemList, err = business_system.Search(paramId)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	}

	var restBusinessSystemList []dto.RestBusinessSystem
	for _, businessSystem := range businessSystemList {
		restBusinessSystemList = append(restBusinessSystemList, dto.NewRestBusinessSystem(businessSystem))
	}
	// IDの昇順にソートする。
	sort.Slice(restBusinessSystemList, func(i, j int) bool {
		return restBusinessSystemList[i].ID < restBusinessSystemList[j].ID
	})
	return c.JSON(http.StatusCreated, restBusinessSystemList)
}

func RegistSystem(c echo.Context) error {
	restSystem := new(dto.RestBusinessSystem)
	err := c.Bind(restSystem)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetBusinessSystem := new(model.BusinessSystem)
	targetBusinessSystem.ID = restSystem.ID
	targetBusinessSystem.Name = restSystem.Name
	targetBusinessSystem.BusinessID = restSystem.Business.ID
	err = business_system.PutItem(*targetBusinessSystem)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "regist business")
}

func DeleteSystem(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	resultFlg := "0"
	// 業務日誌で使用されているか確認
	jobSheetList, err := jobsheet.SearchForBusinessSystem(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	if len(jobSheetList) > 0 {
		// 使用されている場合は削除しない
		resultFlg = "1"
	} else {
		// 使用されていない場合は削除する。
		err = business_system.DeleteItem(id)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	return c.String(http.StatusOK, resultFlg)
}
