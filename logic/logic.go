package logic

import (
	"github.com/RGITHackathonFall2024/auth/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	return utils.Must(token.SignedString(nil))
}
