package service

func (s *Server) routes() {
	// User
	s.router.HandleFunc("/register", s.log(s.handleRegister())).Methods("POST")
	s.router.HandleFunc("/token", s.log(s.handleGetToken())).Methods("POST")

	// History
	s.router.HandleFunc("/entries", s.log(s.loggedIn(s.handleGetEntries()))).Methods("GET")
	s.router.HandleFunc("/entries", s.log(s.loggedIn(s.handlePutEntries()))).Methods("POST")
}
