package app

import (
	"net"
	"net/http"
)

func (s *resolver) StartHTTP(listenAddr string) error {
	server := http.Server{
		Addr: listenAddr,
	}
	http.HandleFunc("/", s.locationHandler)
	go func() {
		<-s.ctx.Done()
		s.logger("shutting down http server: %s", server.Shutdown(s.ctx))
	}()
	s.logger("Starting http server at %s", listenAddr)
	return server.ListenAndServe()
}

func (s *resolver) locationHandler(w http.ResponseWriter, r *http.Request) {
	if e := r.ParseForm(); e != nil {
		s.logger("ParseForm: %s", e)
		http.Error(w, "parse error", http.StatusBadRequest)
		return
	}
	ip := net.ParseIP(r.Form.Get("ip"))
	if ip == nil {
		http.Error(w, "invalid ip", http.StatusBadRequest)
		return
	}

	resp, err := s.GetLocationJSON(ip, r.Form.Get("format"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)

}
