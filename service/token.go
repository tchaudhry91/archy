package service

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func (s *Server) buildToken(user string) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"iat":  time.Now().AddDate(1, 0, 0).Unix(),
	})
	return token.SignedString(s.signingSecret)
}
