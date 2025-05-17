package jwt

import (
	"jobsheet-go-aws2/database/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

// JWT生成
func CreateToken(user *model.User) (string, error) {
	claims := &JwtCustomClaims{
		user.Id,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, nil
}

// JWTからユーザ情報を取得する
func GetUserInfo(c echo.Context) *JwtCustomClaims {
	// JWTの認証を通過した場合、デフォルトでユーザー情報がc.Get("user")に格納される
	usr := c.Get("user").(*jwt.Token)
	claims := usr.Claims.(*JwtCustomClaims)
	return claims
}
