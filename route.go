package main

import (
	"jobsheet-go-aws2/constant"
	jobsheetjwt "jobsheet-go-aws2/jwt"
	"jobsheet-go-aws2/service"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	// JWT設定
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
		// ログイン処理だけJWT認証から除外する。
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == constant.BASE_URL+"/user/login/" {
				return true
			} else {
				return false
			}
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jobsheetjwt.JwtCustomClaims)
		},
	}))

	e.POST(constant.BASE_URL+"/user/login/", service.Login)
	e.GET(constant.BASE_URL+"/user/getList/", service.GetUserList)
	e.POST(constant.BASE_URL+"/user/regist/:mode", service.RegistUser)
	e.DELETE(constant.BASE_URL+"/user/delete/:id", service.DeleteUser)
	e.POST(constant.BASE_URL+"/user/changePassword/", service.ChangePassword)
	e.POST(constant.BASE_URL+"/user/changeSeq/", service.ChangeSeq)

	e.GET(constant.BASE_URL+"/jobsheet/get/:id", service.GetJobSheet)
	e.POST(constant.BASE_URL+"/jobsheet/search/", service.SearchJobSheet)
	e.POST(constant.BASE_URL+"/jobsheet/regist/", service.RegistJobSheet)
	e.DELETE(constant.BASE_URL+"/jobsheet/delete/:id", service.DeleteJobSheet)
	e.POST(constant.BASE_URL+"/jobsheet/download/", service.DownloadJobSheet)
	e.GET(constant.BASE_URL+"/jobsheet/pdf/:id", service.PdfJobSheet)
	e.GET(constant.BASE_URL+"/jobsheet/stats/:year", service.GetStatsJobSheet)
	e.GET(constant.BASE_URL+"/jobsheet/stats/:year/:systemId", service.GetChart)

	e.POST(constant.BASE_URL+"/attachment/regist/:id", service.RegistAttachment)
	e.GET(constant.BASE_URL+"/attachment/download/:id/:seqNo", service.DownloadAttachment)
	e.DELETE(constant.BASE_URL+"/attachment/delete/:id/:seqNo", service.DeleteAttachment)

	e.GET(constant.BASE_URL+"/holiday/getList/", service.GetHolidayList)
	e.POST(constant.BASE_URL+"/holiday/regist/", service.RegistHoliday)

	e.GET(constant.BASE_URL+"/business/getList/", service.GetBusinessList)
	e.POST(constant.BASE_URL+"/business/regist/", service.RegistBusiness)
	e.DELETE(constant.BASE_URL+"/business/delete/:id", service.DeleteBusiness)

	e.GET(constant.BASE_URL+"/system/getList/:id", service.GetBusinessSystemList)
	e.GET(constant.BASE_URL+"/system/getList/", service.GetBusinessSystemList)
	e.POST(constant.BASE_URL+"/system/regist/", service.RegistSystem)
	e.DELETE(constant.BASE_URL+"/system/delete/:id", service.DeleteSystem)

	e.GET(constant.BASE_URL+"/client/getList/", service.GetClientList)
	e.POST(constant.BASE_URL+"/client/regist/", service.RegistClient)
	e.DELETE(constant.BASE_URL+"/client/delete/:id", service.DeleteClient)

	e.GET(constant.BASE_URL+"/inquiry/getList/", service.GetInquiryList)
}
