package service

func (s *Server) routes() {
	s.router.HandleFunc("/entries", s.log(s.handleGetEntries())).Methods("GET")
	s.router.HandleFunc("/entries", s.log(s.handlePutEntries())).Methods("POST")
}
