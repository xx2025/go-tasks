package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"go-tasks/boot/config"
	"time"
)

type JwtCustomClaims struct {
	UserId    int
	RoleId    int
	LoginTime int64 //登录时间戳
	jwt.StandardClaims
}

func NewToken(userId, roleId int) (string, int64, error) {
	loginTime := time.Now().Unix()
	claims := JwtCustomClaims{
		UserId:    userId,
		RoleId:    roleId,
		LoginTime: loginTime,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "go-tasks",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  loginTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := config.MasterConfigure.JwtSecretKey
	if secretKey == "" {
		return "", 0, nil
	}

	signedToken, err := token.SignedString([]byte(secretKey))

	return signedToken, loginTime, err
}

func ParseToken(tokenString string) (*JwtCustomClaims, error) {
	// 解析 Token
	token, _ := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.MasterConfigure.JwtSecretKey), nil
	})

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {

		return claims, nil
	} else {
		return nil, fmt.Errorf("token已失效")
	}

}
