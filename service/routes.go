package service

func (s *Server) routes() {
	s.router.HandleFunc("/entries", s.log(s.loggedIn(s.handleGetEntries()))).Methods("GET")
	s.router.HandleFunc("/entries", s.log(s.loggedIn(s.handlePutEntries()))).Methods("POST")
}
