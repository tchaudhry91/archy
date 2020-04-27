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

	"github.com/tchaudhry91/archy/history"

	"github.com/tchaudhry91/archy/service/store"
)

// Server is the base struct for the HTTP history service
type Server struct {
	db            *store.MongoStore
	logger        log.Logger
	router        *mux.Router
	signingSecret []byte
	httpS         http.Server
}

// NewServer returns a new ZSH History handling service
func NewServer(db *store.MongoStore, logger log.Logger, router *mux.Router, bindAddr, signingSecret string) *Server {
	s := &Server{
		db:     db,
		logger: logger,
		router: router,
		httpS: http.Server{
			Addr:    bindAddr,
			Handler: router,
		},
		signingSecret: []byte(signingSecret),
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

func (s *Server) respond(w http.ResponseWriter, req *http.Request, data interface{}, status int, err error) {
	defer s.logger.Log("status", status, "err", err)
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
			Start   uint64 `json:"start,omitempty"`
			End     uint64 `json:"end,omitempty"`
			Command string `json:"command,omitempty"`
			Machine string `json:"machine,omitempty"`
			Limit   int64  `json:"limit,omitempty"`
		}
		type Response struct {
			Entries []history.Entry `json:"entries,omitempty"`
			Err     string          `json:"err,omitempty"`
		}

		userV := req.Context().Value(userK)
		if userV == nil {
			s.respond(w, req, nil, http.StatusUnauthorized, errors.New("No user found"))
			return
		}

		user := userV.(string)

		r := Request{
			Start:   uint64(time.Now().AddDate(0, 0, -7).Unix()),
			End:     uint64(time.Now().Unix()),
			Machine: "",
			Command: "",
			Limit:   100,
		}
		// Decode request
		qparams := req.URL.Query()
		if startS, ok := qparams["start"]; ok {
			start, err := strconv.Atoi(startS[0])
			if err != nil {
				s.respond(w, req, nil, http.StatusBadRequest, errors.New("Invalid parameter"))
				return
			}
			r.Start = uint64(start)
		}
		if endS, ok := qparams["end"]; ok {
			end, err := strconv.Atoi(endS[0])
			if err != nil {
				s.respond(w, req, nil, http.StatusBadRequest, errors.New("Invalid parameter"))
				return
			}
			r.End = uint64(end)
		}
		if limitS, ok := qparams["limit"]; ok {
			limit, err := strconv.Atoi(limitS[0])
			if err != nil {
				s.respond(w, req, nil, http.StatusBadRequest, errors.New("Invalid parameter"))
				return
			}
			r.Limit = int64(limit)
		}

		if machineQ, ok := qparams["machine"]; ok {
			r.Machine = machineQ[0]
		}

		if commandQ, ok := qparams["command"]; ok {
			r.Command = commandQ[0]
		}

		filter := store.SelectTimerangeFilter(r.Start, r.End)
		if r.Command != "" {
			filter = store.AndMergeFilters(filter, store.SearchCommandFilter(r.Command))
		}
		if r.Machine != "" {
			filter = store.AndMergeFilters(filter, store.SelectMachineFilter(r.Machine))
		}

		// Fetch the entries
		entries, err := s.db.GetEntries(req.Context(), user, filter, r.Limit)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, req, Response{Entries: entries}, http.StatusOK, nil)
		return
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

		userV := req.Context().Value(userK)
		if userV == nil {
			s.respond(w, req, nil, http.StatusUnauthorized, errors.New("No user supplied"))
			return
		}

		user := userV.(string)

		// Decode Request
		err := s.jsonDecode(w, req, &r)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusBadRequest, err)
			return
		}

		// PutEntries
		changed, err := s.db.StoreEntries(req.Context(), user, r.Entries)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, req, Response{Updated: changed}, http.StatusOK, nil)
		return
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			User     string `json:"user,omitempty"`
			Password string `json:"password,omitempty"`
		}
		type Response struct {
			Err string `json:"err,omitempty"`
		}
		r := Request{}

		err := s.jsonDecode(w, req, &r)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusBadRequest, err)
			return
		}

		user, err := store.NewUser(r.User, r.Password)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusInternalServerError, err)
			return
		}

		// Register the user
		err = s.db.PutUser(req.Context(), user)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, req, Response{}, http.StatusOK, err)
	}
}

func (s *Server) handleGetToken() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			User     string `json:"user,omitempty"`
			Password string `json:"password,omitempty"`
		}
		type Response struct {
			Token string `json:"token,omitempty"`
			Err   string `json:"err,omitempty"`
		}

		r := Request{}

		err := s.jsonDecode(w, req, &r)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusBadRequest, err)
			return
		}

		user, err := s.db.GetUser(req.Context(), r.User)
		if err != nil {
			s.respond(w, req, Response{Err: err.Error()}, http.StatusInternalServerError, err)
			return
		}
		if !user.CheckPassword(r.Password) {
			s.respond(w, req, Response{Err: "Incorrect password"}, http.StatusUnauthorized, errors.New("Incorrect password"))
		}
		token, err := s.buildToken(user.User)
		if err != nil {
			s.respond(w, req, Response{Err: "Failed to build token"}, http.StatusInternalServerError, err)
		}
		s.respond(w, req, Response{Token: token}, http.StatusOK, nil)
	}
}
