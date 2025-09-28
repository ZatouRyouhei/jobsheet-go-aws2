package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"jobsheet-go-aws2/database/holiday"

	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/dto"
	"log/slog"
	"net/http"
	"sort"

	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/labstack/echo/v4"
)

func GetHolidayList(c echo.Context) error {
	holidayList, err := holiday.Scan()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var restHolidayList []dto.RestHoliday
	for _, holiday := range holidayList {
		restHolidayList = append(restHolidayList, dto.RestHoliday{
			Holiday: holiday.Holiday,
			Name:    holiday.Name,
		})
	}
	// 日付の昇順にソートする。
	sort.Slice(restHolidayList, func(i, j int) bool {
		return restHolidayList[i].Holiday < restHolidayList[j].Holiday
	})
	return c.JSON(http.StatusCreated, restHolidayList)
}

func RegistHoliday(c echo.Context) error {
	// テーブルをクリアする。
	err := holiday.ClearTable()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	// ファイルは内閣府のホームページに公開されているCSVファイルの想定
	// 文字コードはUTF8（元ファイルはSJISだがUTF8に変換したうえで取り込む）
	// ヘッダーあり
	// ダブルクォーテーションなし
	form, err := c.FormFile("file")
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	src, err := form.Open()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "file cant open")
	}
	defer src.Close()

	// SJISの文字コードとして読み込む（SJISの変換がうまくいかない）
	// reader := transform.NewReader(src, japanese.ShiftJIS.NewDecoder())
	// s := bufio.NewScanner(reader)
	// UTF8として読み込む
	s := bufio.NewScanner(src)

	// エラーリスト
	errorList := []dto.RestErrorMessage{}

	// データ読み込み
	var holidays []model.Holiday
	rowNum := 1 // 行番号
	for s.Scan() {
		// ヘッダー行は無視する
		if rowNum == 1 {
			rowNum += 1
			continue
		}
		result := strings.Split(s.Text(), ",")
		if len(result) != 2 {
			errorList = append(errorList, dto.RestErrorMessage{LineNo: rowNum, ErrorMsg: "フォーマットエラー（日付と祝日名称を入力してください。）"})
			rowNum += 1
			continue
		}
		// スラッシュ区切りをハイフン区切りに変換する
		holidayArr := strings.Split(result[0], "/")
		holiday := holidayArr[0] + "-" + fmt.Sprintf("%0*s", 2, holidayArr[1]) + "-" + fmt.Sprintf("%0*s", 2, holidayArr[2])
		name := result[1]
		holidays = append(holidays, model.Holiday{
			Holiday: holiday,
			Name:    name,
		})
		rowNum += 1
	}

	// 入力データ0件の時はエラー
	if rowNum-1 < 2 {
		errorList = append(errorList, dto.RestErrorMessage{LineNo: 0, ErrorMsg: "ヘッダーを含め2行以上入力してください。"})
	}

	if len(errorList) == 0 {
		written, err := holiday.BatchWriteItem(holidays)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		slog.Info("Info", slog.Any("result info", strconv.Itoa(written)+"件取り込みました。"))
	}

	return c.JSON(http.StatusCreated, errorList)
}

// SJISからUTF-8に変換する関数
func convertSJISToUTF8(input string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(input)), japanese.ShiftJIS.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
