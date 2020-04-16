package service

func (s *Server) routes() {
	s.router.HandleFunc("/entries", s.handleGetEntries()).Methods("GET")
	s.router.HandleFunc("/entries", s.handlePutEntries()).Methods("POST")
}
