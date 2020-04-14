package service

import (
	"net/http"

	"github.com/tchaudhry91/zsh-archaeologist/service/store"
)

// Server is the base struct for the HTTP history service
type Server struct {
	db       *store.MongoStore
	bindAddr string
}

func (s *Server) handleGetEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			user  string
			start int64
			end   int64
		}
		type Response struct {
		}
	}
}
