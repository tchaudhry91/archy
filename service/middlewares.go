package service

import (
	"net/http"
	"time"
)

func (s *Server) log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func(begin time.Time) {
			s.logger.Log("Path", req.URL.Path, "Method", req.Method, "Took", time.Since(begin))
		}(time.Now())
		h(w, req)
	}
}
