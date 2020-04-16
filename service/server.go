package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/tchaudhry91/zsh-archaeologist/history"

	"github.com/tchaudhry91/zsh-archaeologist/service/store"
)

// Server is the base struct for the HTTP history service
type Server struct {
	db     *store.MongoStore
	logger log.Logger
	router *mux.Router
	httpS  http.Server
}

// NewServer returns a new ZSH History handling service
func NewServer(db *store.MongoStore, logger log.Logger, router *mux.Router, bindAddr string) *Server {
	s := &Server{
		db:     db,
		logger: logger,
		router: router,
		httpS: http.Server{
			Addr:    bindAddr,
			Handler: router,
		},
	}
	s.routes()
	return s
}

// Start begins listening for requests on the bindAddr. Blocks.
func (s *Server) Start() error {
	return s.httpS.ListenAndServe()
}

// Shutdown gracefully terminates the server
func (s *Server) Shutdown() error {
	return s.httpS.Shutdown(context.TODO())
}

func (s *Server) respond(w http.ResponseWriter, req *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Log("err", errors.Errorf("Failed to Encode result to JSON:%v", err))
		}
	}
}

func (s *Server) jsonDecode(w http.ResponseWriter, req *http.Request, v interface{}) error {
	return json.NewDecoder(req.Body).Decode(v)
}

func (s *Server) handleGetEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Start uint64 `json:"start,omitempty"`
			End   uint64 `json:"end,omitempty"`
			Limit int64  `json:"limit,omitempty"`
		}
		type Response struct {
			Entries []history.Entry `json:"entries,omitempty"`
			Err     string          `json:"err,omitempty"`
		}

		user := "tchaudhry"

		r := Request{
			Start: uint64(time.Now().AddDate(0, 0, 7).Unix()),
			End:   uint64(time.Now().Unix()),
			Limit: 100,
		}
		// Decode request
		qparams := req.URL.Query()
		if startS, ok := qparams["start"]; ok {
			start, err := strconv.Atoi(startS[0])
			if err != nil {
				s.respond(w, req, nil, http.StatusBadGateway)
			}
			r.Start = uint64(start)
		}
		if endS, ok := qparams["end"]; ok {
			end, err := strconv.Atoi(endS[0])
			if err != nil {
				s.respond(w, req, nil, http.StatusBadGateway)
			}
			r.End = uint64(end)
		}
		if limitS, ok := qparams["limit"]; ok {
			limit, err := strconv.Atoi(limitS[0])
			if err != nil {
				s.respond(w, req, nil, http.StatusBadGateway)
			}
			r.Limit = int64(limit)
		}

		// Fetch the entries
		entries, err := s.db.GetEntries(req.Context(), user, store.SelectTimerangeFilter(r.Start, r.End), r.Limit)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusInternalServerError)
		}
		s.respond(w, req, Response{Entries: entries}, http.StatusOK)
	}
}

func (s *Server) handlePutEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Entries []history.Entry `json:"entries,omitempty"`
		}
		type Response struct {
			Updated int64  `json:"updated,omitempty"`
			Err     string `json:"err,omitempty"`
		}
		r := Request{}
		user := "tchaudhry"

		// Decode Request
		s.jsonDecode(w, req, &r)

		// PutEntries
		changed, err := s.db.StoreEntries(req.Context(), user, r.Entries)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusInternalServerError)
		}

		s.respond(w, req, Response{Updated: changed, Err: err.Error()}, http.StatusOK)
	}
}
