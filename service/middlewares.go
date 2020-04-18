package service

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type key string

const userK key = "user"

func (s *Server) log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func(begin time.Time) {
			s.logger.Log("Path", req.URL.Path, "Method", req.Method, "Took", time.Since(begin))
		}(time.Now())
		h(w, req)
	}
}

func (s *Server) loggedIn(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract Token
		tokenStr := req.Header.Get("token")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return s.signingSecret, nil
		})
		if err != nil {
			s.respond(w, req, nil, http.StatusUnauthorized, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user, ok := claims["user"]
			if !ok {
				s.respond(w, req, nil, http.StatusUnauthorized, err)
				return
			}
			ctx := context.WithValue(req.Context(), userK, user.(string))
			req = req.WithContext(ctx)
			h(w, req)
		} else {
			s.respond(w, req, nil, http.StatusUnauthorized, errors.New("User not found in token claims"))
			return
		}
	}
}
