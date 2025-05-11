package service

import (
	"fmt"
	"jobsheet-go-aws2/database/attachment"
	"jobsheet-go-aws2/database/business"
	"jobsheet-go-aws2/database/business_system"
	"jobsheet-go-aws2/database/client"
	"jobsheet-go-aws2/database/inquiry"
	"jobsheet-go-aws2/database/jobsheet"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/database/user"
	"jobsheet-go-aws2/dto"
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
)

func RegistJobSheet(c echo.Context) error {
	var restJobSheet = new(dto.RestJobSheet)
	err := c.Bind(restJobSheet)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	newflg := false
	// 新規登録の時はIDを自動採番
	// 登録時に同じIDがあった場合は上書き（後勝ち）する。
	if restJobSheet.ID == "" {
		newflg = true
		now := time.Now()
		idHeader := now.Format("2006-01")
		nextSeqNo := "001"
		idHeaderList, err := jobsheet.ScanForIDHeader(idHeader)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		if len(idHeaderList) > 0 {
			// IDの降順にソートする
			sort.Slice(idHeaderList, func(i, j int) bool {
				return idHeaderList[i].ID > idHeaderList[j].ID
			})
			maxSeqNo, err := strconv.Atoi(idHeaderList[0].ID[8:])
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			nextSeqNo = fmt.Sprintf("%03d", maxSeqNo+1)
		}
		restJobSheet.ID = idHeader + "-" + nextSeqNo
	}

	targetJobSheet := new(model.JobSheet)
	targetJobSheet.ID = restJobSheet.ID
	targetJobSheet.ClientID = restJobSheet.ClientID
	targetJobSheet.BusinessSystemID = restJobSheet.BusinessSystemID
	targetJobSheet.InquiryID = restJobSheet.InquiryID
	targetJobSheet.Department = restJobSheet.Department
	targetJobSheet.Person = restJobSheet.Person
	targetJobSheet.OccurDate = restJobSheet.OccurDate
	targetJobSheet.OccurTime = restJobSheet.OccurTime
	targetJobSheet.Title = restJobSheet.Title
	targetJobSheet.Content = restJobSheet.Content
	targetJobSheet.ContactID = restJobSheet.ContactID
	targetJobSheet.LimitDate = restJobSheet.LimitDate
	targetJobSheet.DealID = restJobSheet.DealID
	targetJobSheet.CompleteDate = restJobSheet.CompleteDate
	targetJobSheet.Support = restJobSheet.Support
	targetJobSheet.ResponseTime = restJobSheet.ResponseTime

	if newflg {
		//新規登録
		err := jobsheet.PutItem(*targetJobSheet)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	} else {
		// 更新
		err := jobsheet.UpdateItem(*targetJobSheet)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	return c.String(http.StatusOK, targetJobSheet.ID)
}

func GetJobSheet(c echo.Context) error {
	id := c.Param("id")
	targetJobsheet, err := jobsheet.GetItem(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetAttachmentList, err := attachment.Scan(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	restJobSheet := dto.NewRestSearchJobSheet(*targetJobsheet, targetAttachmentList)
	return c.JSON(http.StatusCreated, restJobSheet)
}

func SearchJobSheet(c echo.Context) error {
	condition := new(dto.RestSearchConditionJobSheet)
	err := c.Bind(condition)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	jobSheetList, err := jobsheet.Search(condition)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var restSearchJobSheetList []dto.RestSearchJobSheet
	for _, targetJobsheet := range jobSheetList {
		targetAttachmentList, err := attachment.Scan(targetJobsheet.ID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		restSearchJobSheetList = append(restSearchJobSheetList, dto.NewRestSearchJobSheet(targetJobsheet, targetAttachmentList))
	}
	// IDの降順にソートする。
	sort.Slice(restSearchJobSheetList, func(i, j int) bool {
		return restSearchJobSheetList[i].ID > restSearchJobSheetList[j].ID
	})
	return c.JSON(http.StatusCreated, restSearchJobSheetList)
}

func DeleteJobSheet(c echo.Context) error {
	id := c.Param("id")
	err := jobsheet.DeleteItem(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	// 添付ファイルも削除する。
	err = attachment.DeleteItems(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "delete jobsheet")
}

func DownloadJobSheet(c echo.Context) error {
	condition := new(dto.RestSearchConditionJobSheet)
	err := c.Bind(condition)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	jobSheetList, err := jobsheet.Search(condition)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	// 自ファイルからの相対バスだとファイルが見つからない。
	// ビルドで生成されるバイナリファイルからの相対パスを指定する。
	f, err := excelize.OpenFile("template/template.xlsx")
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, "fileOpenerror")
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Error", slog.Any("error", err))
		}
	}()
	// シート名
	sheetName := "Sheet1"
	// ロケーション
	loc, _ := time.LoadLocation("Asia/Tokyo")
	// 今日の日付
	today := time.Now()
	// セルスタイル
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
	})
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}

	for i, jobSheet := range jobSheetList {
		// セルスタイル設定
		f.SetCellStyle(sheetName, "A"+strconv.Itoa(i+3), "Q"+strconv.Itoa(i+3), style)
		// 番号
		f.SetCellValue(sheetName, "A"+strconv.Itoa(i+3), jobSheet.ID)
		// ステータス
		status := ""
		if jobSheet.CompleteDate == "" {
			if jobSheet.LimitDate != "" {
				limitDate, err := time.ParseInLocation("2006-01-02 15:04:05", jobSheet.LimitDate+" "+"23:59:59", loc)
				if err != nil {
					slog.Error("Error", slog.Any("error", err))
				}
				if limitDate.Before(today) {
					// 期限超過
					status = "期限超過"
				} else {
					// 期限前
					diffDays := limitDate.Sub(today).Hours() / 24
					if diffDays <= 3 {
						status = "あと" + strconv.Itoa(int(diffDays)+1) + "日"
					}
				}
			}
		} else {
			status = "完了"
		}
		f.SetCellValue(sheetName, "B"+strconv.Itoa(i+3), status)
		// 顧客
		client, err := client.GetItem(jobSheet.ClientID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		f.SetCellValue(sheetName, "C"+strconv.Itoa(i+3), client.Name)
		// 業務
		businessSystem, err := business_system.GetItem(jobSheet.BusinessSystemID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		business, err := business.GetItem(businessSystem.BusinessID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		f.SetCellValue(sheetName, "D"+strconv.Itoa(i+3), business.Name)
		// システム
		f.SetCellValue(sheetName, "E"+strconv.Itoa(i+3), businessSystem.Name)
		// 問合せ区分
		inquiry, err := inquiry.GetItem(jobSheet.InquiryID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		f.SetCellValue(sheetName, "F"+strconv.Itoa(i+3), inquiry.Name)
		// 部署
		f.SetCellValue(sheetName, "G"+strconv.Itoa(i+3), jobSheet.Department)
		// 担当者
		f.SetCellValue(sheetName, "H"+strconv.Itoa(i+3), jobSheet.Person)
		// 発生日時
		occurDateTime := jobSheet.OccurDate + " " + jobSheet.OccurTime
		f.SetCellValue(sheetName, "I"+strconv.Itoa(i+3), occurDateTime)
		// 窓口
		contact, err := user.GetItem(jobSheet.ContactID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		f.SetCellValue(sheetName, "J"+strconv.Itoa(i+3), contact.Name)
		// タイトル
		f.SetCellValue(sheetName, "K"+strconv.Itoa(i+3), jobSheet.Title)
		// 内容
		f.SetCellValue(sheetName, "L"+strconv.Itoa(i+3), jobSheet.Content)
		// 完了期限
		f.SetCellValue(sheetName, "M"+strconv.Itoa(i+3), jobSheet.LimitDate)
		// 対応詳細
		f.SetCellValue(sheetName, "N"+strconv.Itoa(i+3), jobSheet.Support)
		// 対応者
		dealer := ""
		if jobSheet.DealID != "" {
			deal, err := user.GetItem(jobSheet.DealID)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			dealer = deal.Name
		}
		f.SetCellValue(sheetName, "O"+strconv.Itoa(i+3), dealer)
		// 完了日
		f.SetCellValue(sheetName, "P"+strconv.Itoa(i+3), jobSheet.CompleteDate)
		// 対応時間
		f.SetCellValue(sheetName, "Q"+strconv.Itoa(i+3), jobSheet.ResponseTime)
	}

	buf, _ := f.WriteToBuffer()
	response := c.Response()
	response.Writer.Header().Set("Content-Disposition", "attachment; filename=業務日誌.xlsx")
	return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func PdfJobSheet(c echo.Context) error {
	id := c.Param("id")
	targetJobSheet, err := jobsheet.GetItem(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pdf.AddPage()
	err = pdf.AddTTFFont("genju", "font/GenJyuuGothic-Regular.ttf")
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	err = pdf.SetFont("genju", "", 10)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	// 出力日
	rect := gopdf.Rect{W: 80, H: 20}
	pdf.SetX(460)
	pdf.SetY(60)
	op := gopdf.CellOption{
		Align: gopdf.Right | gopdf.Bottom,
	}
	now := time.Now()
	pdf.CellWithOption(&rect, now.Format("2006年01月02日"), op)
	// ID
	drawText(&pdf, 30, 30, targetJobSheet.ID)
	// 顧客
	client, err := client.GetItem(targetJobSheet.ClientID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	drawText(&pdf, 60, 110, client.Name)
	// 業務
	businessSystem, err := business_system.GetItem(targetJobSheet.BusinessSystemID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	business, err := business.GetItem(businessSystem.BusinessID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	drawText(&pdf, 230, 110, business.Name)
	// システム
	drawText(&pdf, 390, 110, businessSystem.Name)
	// 問合せ区分
	inquiry, err := inquiry.GetItem(targetJobSheet.InquiryID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	drawText(&pdf, 60, 160, inquiry.Name)
	// 部署
	drawText(&pdf, 230, 160, targetJobSheet.Department)
	// 担当者
	drawText(&pdf, 390, 160, targetJobSheet.Person)
	// 発生日
	occurDateArr := strings.Split(targetJobSheet.OccurDate, "-")
	occurDateStr := occurDateArr[0] + "年" + occurDateArr[1] + "月" + occurDateArr[2] + "日"
	occurTimeArr := strings.Split(targetJobSheet.OccurTime, ":")
	occurTimeStr := occurTimeArr[0] + "時" + occurTimeArr[1] + "分"
	drawText(&pdf, 60, 220, occurDateStr+" "+occurTimeStr)
	// 窓口
	contact, err := user.GetItem(targetJobSheet.ContactID)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	drawText(&pdf, 230, 220, contact.Name)
	// // タイトル
	drawText(&pdf, 60, 270, targetJobSheet.Title)
	// 内容
	// 改行の可能性がある場合は、pdf.MultiCellもしくはpdf.MultiCellWithOptionを使用する。
	// op := gopdf.CellOption{
	// 	Align: gopdf.Left,
	// 	// セルの幅にテキストがおさまらないときの挙動 pdf.MultiCellWithOptionを使用するときのオプション。
	// 	BreakOption: &gopdf.BreakOption{
	// 		// 単語の途中でも改行するモード
	// 		Mode: gopdf.BreakModeStrict,
	// 		// BreakModeStrictの場合で単語の途中で改行される場合のセパレータ文字列
	// 		Separator: "-",
	// 		// 単語の途中では改行しないモード
	// 		// Mode:           gopdf.BreakModeIndicatorSensitive,
	// 		// BreakModeIndicatorSensitiveの場合に単語の区切りとなる文字を指定
	// 		// BreakIndicator: ' ',
	// 	},
	// }
	// テキストの途中に改行コードが入っている場合の処理
	rect = gopdf.Rect{W: 480, H: 20}
	pdf.SetX(60)
	pdf.SetY(330)
	contents := strings.Split(targetJobSheet.Content, "\n")
	for _, content := range contents {
		pdf.MultiCell(&rect, content)
	}
	// 完了期限
	if targetJobSheet.LimitDate != "" {
		limitDateArr := strings.Split(targetJobSheet.LimitDate, "-")
		limitDateStr := limitDateArr[0] + "年" + limitDateArr[1] + "月" + limitDateArr[2] + "日"
		drawText(&pdf, 60, 500, limitDateStr)
	}
	// 対応詳細
	rect = gopdf.Rect{W: 480, H: 20}
	pdf.SetX(60)
	pdf.SetY(560)
	supports := strings.Split(targetJobSheet.Support, "\n")
	for _, support := range supports {
		pdf.MultiCell(&rect, support)
	}
	// 対応者
	dealer := ""
	if targetJobSheet.DealID != "" {
		deal, err := user.GetItem(targetJobSheet.DealID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		dealer = deal.Name
	}
	drawText(&pdf, 60, 730, dealer)
	// 完了日
	if targetJobSheet.CompleteDate != "" {
		completeDateArr := strings.Split(targetJobSheet.CompleteDate, "-")
		completeDateStr := completeDateArr[0] + "年" + completeDateArr[1] + "月" + completeDateArr[2] + "日"
		drawText(&pdf, 230, 730, completeDateStr)
	}
	// 対応時間
	drawText(&pdf, 390, 730, strconv.FormatFloat(targetJobSheet.ResponseTime, 'f', -1, 64))

	A4 := *gopdf.PageSizeA4
	// 引数の2つめはテンプレートファイルのページ番号
	tp1 := pdf.ImportPage("template/jobSheet.pdf", 1, "/MediaBox")
	pdf.UseImportedTemplate(tp1, 0, 0, A4.W, A4.H)

	// 改ページする場合は新たにテンプレートのページを追加する。
	// pdf.AddPage()
	// tp1 = pdf.ImportPage("template/jobSheet.pdf", 1, "/MediaBox")
	// pdf.UseImportedTemplate(tp1, 0, 0, A4.W, A4.H)

	// 位置合わせ用の罫線表示　完成後にコメントアウトする
	// A4Tate := gopdf.Rect{W: A4.W, H: A4.H}
	// drawGrid(&pdf, &A4Tate)

	response := c.Response()
	response.Writer.Header().Set("Content-Disposition", "attachment; filename=業務日誌.pdf")
	return c.Blob(http.StatusOK, "application/pdf", pdf.GetBytesPdf())
}

func drawText(pdf *gopdf.GoPdf, x float64, y float64, s string) {
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, s)
}

func drawGrid(pdf *gopdf.GoPdf, page *gopdf.Rect) {
	ww := 10.0
	for i := 1; i < int(page.W/ww); i++ {
		if i%10 == 0 {
			pdf.SetLineWidth(0.8)
			pdf.SetStrokeColor(50, 50, 100)
		} else {
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(100, 100, 130)
		}
		x := float64(i) * ww
		pdf.Line(x, 0, x, page.H)
	}
	for i := 1; i < int(page.H/ww); i++ {
		if i%10 == 0 {
			pdf.SetLineWidth(0.8)
			pdf.SetStrokeColor(50, 50, 100)
		} else {
			pdf.SetLineWidth(0.3)
			pdf.SetStrokeColor(100, 100, 130)
		}
		y := float64(i) * ww
		pdf.Line(0, y, page.W, y)
	}
}

func GetStatsJobSheet(c echo.Context) error {
	year := c.Param("year")
	restStatJobSheetList := []dto.RestStatJobSheet{}
	systemList, err := business_system.Scan()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	// IDの昇順にソートする。
	sort.Slice(systemList, func(i, j int) bool {
		return systemList[i].ID < systemList[j].ID
	})

	slog.Info("Info", slog.Any("systemList", systemList)) //　デバッグ用

	for _, system := range systemList {
		restStatJobSheet := new(dto.RestStatJobSheet)
		restStatJobSheet.BusinessSystem = dto.NewRestBusinessSystem(system)

		slog.Info("Info", slog.Any("restStatJobSheet.BusinessSystem", restStatJobSheet.BusinessSystem)) //　デバッグ用

		// 1年の合計時間用
		responseTimeSum := 0.0
		// 4月から3月まで件数をカウントする。
		for monthIdx := 1; monthIdx <= 12; monthIdx++ {
			statsYear, err := strconv.Atoi(year)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			if monthIdx <= 3 {
				statsYear += 1
			}
			jst, _ := time.LoadLocation("Asia/Tokyo")
			dateFrom := time.Date(statsYear, time.Month(monthIdx), 1, 0, 0, 0, 0, jst)
			dateFromStr := dateFrom.Format("2006-01-02")
			nextStatsYear := statsYear
			nextMonthIdx := monthIdx + 1
			if nextMonthIdx == 13 {
				nextStatsYear += 1
				nextMonthIdx = 1
			}
			dateTo := time.Date(nextStatsYear, time.Month(nextMonthIdx), 1, 0, 0, 0, 0, jst)
			dateToStr := dateTo.Format("2006-01-02")
			// 発生件数
			occurJobSheetList, err := jobsheet.ScanForOccur(dateFromStr, dateToStr, system.ID)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			occurCnt := len(occurJobSheetList)

			slog.Info("Info", slog.Any("occurJobSheetList", occurJobSheetList)) //　デバッグ用
			slog.Info("Info", slog.Any("occurCnt", occurCnt))                   //　デバッグ用

			// 完了件数
			completeJobSheetList, err := jobsheet.ScanForComplete(dateFromStr, dateToStr, system.ID)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			completeCnt := len(completeJobSheetList)
			// 対応時間
			responseTime := 0.0
			for _, v := range completeJobSheetList {
				responseTime += v.ResponseTime
			}
			switch monthIdx {
			case 1:
				restStatJobSheet.OccurCnt1 = occurCnt
				restStatJobSheet.CompleteCnt1 = completeCnt
				restStatJobSheet.ResponseTime1 = responseTime
			case 2:
				restStatJobSheet.OccurCnt2 = occurCnt
				restStatJobSheet.CompleteCnt2 = completeCnt
				restStatJobSheet.ResponseTime2 = responseTime
			case 3:
				restStatJobSheet.OccurCnt3 = occurCnt
				restStatJobSheet.CompleteCnt3 = completeCnt
				restStatJobSheet.ResponseTime3 = responseTime
			case 4:
				restStatJobSheet.OccurCnt4 = occurCnt
				restStatJobSheet.CompleteCnt4 = completeCnt
				restStatJobSheet.ResponseTime4 = responseTime
			case 5:
				restStatJobSheet.OccurCnt5 = occurCnt
				restStatJobSheet.CompleteCnt5 = completeCnt
				restStatJobSheet.ResponseTime5 = responseTime
			case 6:
				restStatJobSheet.OccurCnt6 = occurCnt
				restStatJobSheet.CompleteCnt6 = completeCnt
				restStatJobSheet.ResponseTime6 = responseTime
			case 7:
				restStatJobSheet.OccurCnt7 = occurCnt
				restStatJobSheet.CompleteCnt7 = completeCnt
				restStatJobSheet.ResponseTime7 = responseTime
			case 8:
				restStatJobSheet.OccurCnt8 = occurCnt
				restStatJobSheet.CompleteCnt8 = completeCnt
				restStatJobSheet.ResponseTime8 = responseTime
			case 9:
				restStatJobSheet.OccurCnt9 = occurCnt
				restStatJobSheet.CompleteCnt9 = completeCnt
				restStatJobSheet.ResponseTime9 = responseTime
			case 10:
				restStatJobSheet.OccurCnt10 = occurCnt
				restStatJobSheet.CompleteCnt10 = completeCnt
				restStatJobSheet.ResponseTime10 = responseTime
			case 11:
				restStatJobSheet.OccurCnt11 = occurCnt
				restStatJobSheet.CompleteCnt11 = completeCnt
				restStatJobSheet.ResponseTime11 = responseTime
			case 12:
				restStatJobSheet.OccurCnt12 = occurCnt
				restStatJobSheet.CompleteCnt12 = completeCnt
				restStatJobSheet.ResponseTime12 = responseTime
			}
			responseTimeSum += responseTime
		}
		// 未完了の件数を求める。
		leftJobSheetList, err := jobsheet.ScanForUnComplete(system.ID)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		leftCnt := len(leftJobSheetList)
		restStatJobSheet.LeftCnt = leftCnt
		// 1年の合計時間
		restStatJobSheet.ResponseTimeSum = responseTimeSum
		restStatJobSheetList = append(restStatJobSheetList, *restStatJobSheet)
	}
	return c.JSON(http.StatusCreated, restStatJobSheetList)
	//return c.String(http.StatusOK, "ok") // デバッグ用
}

func GetChart(c echo.Context) error {
	year := c.Param("year")
	paramSystemId := c.Param("systemId")
	systemId, err := strconv.Atoi(paramSystemId)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetSystem, err := business_system.GetItem(systemId)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	if targetSystem.ID != 0 {
		restStatJobSheet := new(dto.RestStatJobSheet)
		restStatJobSheet.BusinessSystem = dto.NewRestBusinessSystem(*targetSystem)
		for monthIdx := 1; monthIdx <= 12; monthIdx++ {
			statsYear, err := strconv.Atoi(year)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			if monthIdx <= 3 {
				statsYear += 1
			}
			jst, _ := time.LoadLocation("Asia/Tokyo")
			dateFrom := time.Date(statsYear, time.Month(monthIdx), 1, 0, 0, 0, 0, jst)
			dateFromStr := dateFrom.Format("2006-01-02")
			nextStatsYear := statsYear
			nextMonthIdx := monthIdx + 1
			if nextMonthIdx == 13 {
				nextStatsYear += 1
				nextMonthIdx = 1
			}
			dateTo := time.Date(nextStatsYear, time.Month(nextMonthIdx), 1, 0, 0, 0, 0, jst)
			dateToStr := dateTo.Format("2006-01-02")
			// 発生件数
			occurJobSheetList, err := jobsheet.ScanForOccur(dateFromStr, dateToStr, systemId)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			occurCnt := len(occurJobSheetList)

			slog.Info("Info", slog.Any("occurJobSheetList", occurJobSheetList)) //　デバッグ用
			slog.Info("Info", slog.Any("occurCnt", occurCnt))                   //　デバッグ用

			// 完了件数
			completeJobSheetList, err := jobsheet.ScanForComplete(dateFromStr, dateToStr, systemId)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			completeCnt := len(completeJobSheetList)
			// 対応時間
			responseTime := 0.0
			for _, v := range completeJobSheetList {
				responseTime += v.ResponseTime
			}
			switch monthIdx {
			case 1:
				restStatJobSheet.OccurCnt1 = occurCnt
				restStatJobSheet.CompleteCnt1 = completeCnt
				restStatJobSheet.ResponseTime1 = responseTime
			case 2:
				restStatJobSheet.OccurCnt2 = occurCnt
				restStatJobSheet.CompleteCnt2 = completeCnt
				restStatJobSheet.ResponseTime2 = responseTime
			case 3:
				restStatJobSheet.OccurCnt3 = occurCnt
				restStatJobSheet.CompleteCnt3 = completeCnt
				restStatJobSheet.ResponseTime3 = responseTime
			case 4:
				restStatJobSheet.OccurCnt4 = occurCnt
				restStatJobSheet.CompleteCnt4 = completeCnt
				restStatJobSheet.ResponseTime4 = responseTime
			case 5:
				restStatJobSheet.OccurCnt5 = occurCnt
				restStatJobSheet.CompleteCnt5 = completeCnt
				restStatJobSheet.ResponseTime5 = responseTime
			case 6:
				restStatJobSheet.OccurCnt6 = occurCnt
				restStatJobSheet.CompleteCnt6 = completeCnt
				restStatJobSheet.ResponseTime6 = responseTime
			case 7:
				restStatJobSheet.OccurCnt7 = occurCnt
				restStatJobSheet.CompleteCnt7 = completeCnt
				restStatJobSheet.ResponseTime7 = responseTime
			case 8:
				restStatJobSheet.OccurCnt8 = occurCnt
				restStatJobSheet.CompleteCnt8 = completeCnt
				restStatJobSheet.ResponseTime8 = responseTime
			case 9:
				restStatJobSheet.OccurCnt9 = occurCnt
				restStatJobSheet.CompleteCnt9 = completeCnt
				restStatJobSheet.ResponseTime9 = responseTime
			case 10:
				restStatJobSheet.OccurCnt10 = occurCnt
				restStatJobSheet.CompleteCnt10 = completeCnt
				restStatJobSheet.ResponseTime10 = responseTime
			case 11:
				restStatJobSheet.OccurCnt11 = occurCnt
				restStatJobSheet.CompleteCnt11 = completeCnt
				restStatJobSheet.ResponseTime11 = responseTime
			case 12:
				restStatJobSheet.OccurCnt12 = occurCnt
				restStatJobSheet.CompleteCnt12 = completeCnt
				restStatJobSheet.ResponseTime12 = responseTime
			}
		}
		return c.JSON(http.StatusCreated, restStatJobSheet)
	} else {
		return c.String(http.StatusBadRequest, "bad request")
	}
}
