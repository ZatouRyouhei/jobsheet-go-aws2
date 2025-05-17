package service

import (
	"crypto/sha512"
	"fmt"
	"jobsheet-go-aws2/database/jobsheet"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/database/user"
	"jobsheet-go-aws2/dto"
	jobsheetjwt "jobsheet-go-aws2/jwt"
	"jobsheet-go-aws2/random"

	"encoding/hex"
	"jobsheet-go-aws2/mail"
	"log/slog"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	paramUser := new(dto.RestUser)
	err := c.Bind(paramUser)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetUser, err := user.PassCheck(paramUser.Id, paramUser.Password)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	// JWT生成
	t, err := jobsheetjwt.CreateToken(targetUser)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, dto.RestLoginUser{
		User: dto.RestUser{
			Id:       targetUser.Id,
			Password: "",
			Name:     targetUser.Name,
			SeqNo:    targetUser.SeqNo,
		},
		Token: t,
	})
}

func GetUserList(c echo.Context) error {
	// JWTからユーザ情報を取得する。
	userInfo := jobsheetjwt.GetUserInfo(c)
	slog.Info("Info", slog.Any("Name", "ユーザID"+userInfo.Id))

	userList, err := user.Scan()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var rUserList []dto.RestUser
	for _, user := range userList {
		rUserList = append(rUserList, dto.NewRestUser(user))
	}
	// IDの昇順にソートする。
	sort.Slice(rUserList, func(i, j int) bool {
		return rUserList[i].SeqNo < rUserList[j].SeqNo
	})
	return c.JSON(http.StatusCreated, rUserList)
}

func RegistUser(c echo.Context) error {
	mode := c.Param("mode") // initial:新規登録、update:更新、password:パスワード初期化
	paramUser := new(dto.RestUser)
	err := c.Bind(paramUser)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	dbUser, err := user.GetItem(paramUser.Id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var targetUser model.User
	if dbUser.Id == "" {
		// 新規登録
		targetUser.Id = paramUser.Id
		targetUser.Name = paramUser.Name
		// パスワード生成 5文字のランダムの文字列
		password, err := random.MakeRandomStr(5)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		// 生成したパスワードをハッシュ値に変換する
		passwordSHA512 := sha512.Sum512([]byte(password))
		targetUser.Password = hex.EncodeToString(passwordSHA512[:])
		// 連番取得
		nextSeqNo := 1
		userList, err := user.Scan()
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		if len(userList) > 0 {
			// IDの降順にソートする。
			sort.Slice(userList, func(i, j int) bool {
				return userList[i].SeqNo > userList[j].SeqNo
			})
			nextSeqNo = userList[0].SeqNo + 1
		}
		targetUser.SeqNo = nextSeqNo
		err = user.PutItem(targetUser)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}

		// メール送信
		sendParam := mail.SendParam{
			Title:  "【業務日誌】ユーザ新規登録",
			Body:   fmt.Sprintf("ユーザを新規登録しました。\n パスワードは %s です。", password),
			MailTo: targetUser.Id,
		}
		err = mail.SendMessage(sendParam)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	} else {
		// 更新もしくはパスワード初期化
		targetUser.Id = paramUser.Id
		targetUser.Name = paramUser.Name
		targetUser.Password = ""
		password := ""
		if mode == "password" {
			// パスワード初期化の時はパスワード生成
			// パスワード生成 5文字のランダムの文字列
			password, err = random.MakeRandomStr(5)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
			// 生成したパスワードをハッシュ値に変換する
			passwordSHA512 := sha512.Sum512([]byte(password))
			targetUser.Password = hex.EncodeToString(passwordSHA512[:])
		}
		// ユーザ情報更新
		err = user.UpdateItem(targetUser)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		// パスワードが初期化された場合はメール通知する
		if mode == "password" {
			// メール送信
			sendParam := mail.SendParam{
				Title:  "【業務日誌】パスワード初期化",
				Body:   fmt.Sprintf("パスワードが初期化されました。\n パスワードは %s です。", password),
				MailTo: targetUser.Id,
			}
			err = mail.SendMessage(sendParam)
			if err != nil {
				slog.Error("Error", slog.Any("error", err))
				return c.String(http.StatusBadRequest, err.Error())
			}
		}
	}
	return c.String(http.StatusOK, "user updated")
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	jobSheetList, err := jobsheet.SearchForUser(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	if len(jobSheetList) > 0 {
		// 業務日誌で使用中のユーザは削除しない。
		return c.String(http.StatusOK, "1")
	} else {
		// 使用されてない場合は削除する。
		err = user.DeleteItem(id)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, "0")
	}
}

func ChangePassword(c echo.Context) error {
	var restUser dto.RestUser
	err := c.Bind(&restUser)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetUser, err := user.GetItem(restUser.Id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	if targetUser.Id != "" {
		targetUser.Password = restUser.Password
		err = user.UpdateItem(*targetUser)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, "change password")
	} else {
		return c.String(http.StatusBadRequest, "bad request")
	}
}

func ChangeSeq(c echo.Context) error {
	var restUserList dto.RestUserList
	err := c.Bind(&restUserList)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	userList := restUserList.UserList
	for seqNo := 1; seqNo <= len(userList); seqNo++ {
		targetId := userList[seqNo-1]
		targetUser, err := user.GetItem(targetId)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		targetUser.SeqNo = seqNo
		err = user.UpdateSeqNo(*targetUser)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	return c.String(http.StatusOK, "change seq")
}
