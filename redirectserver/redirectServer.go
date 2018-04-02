package redirectserver

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Server settings for redirect server
type Server struct {
	// ip:port to listen on, for all interfaces empty, e.g. ":8080"
	listenAddress string
	// hostname for administration of redirects
	adminHost string
	// map of all redirects: hostname, URL -> target
	Redirects hostRedirects

	// mux for handlers
	mux *http.ServeMux

	//logger
	logger *log.Logger
}

//NewServer creates new server
func NewServer(listenAddress string, adminHost string, opts ...Option) *Server {
	s := &Server{
		listenAddress: listenAddress,
		adminHost:     adminHost,
		Redirects:     hostRedirects{},
		mux:           http.DefaultServeMux,
		logger:        log.New(ioutil.Discard, "redirectServer", log.LstdFlags),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.mux.HandleFunc("/", s.Redirects.Handler)

	return s
}

//StartServer calls http.ListenAndServe
func (s *Server) StartServer() error {
	log.Printf("Starting server on address %v", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, s.mux)
}

// Option defines options to set for server
type Option func(*Server)

// WithLogger allows to pass a custom logger
func WithLogger(logger *log.Logger) Option {
	return func(s *Server) { s.logger = logger }
}

//WithMux allows to pass a custom mux
func WithMux(mux *http.ServeMux) Option {
	return func(s *Server) { s.mux = mux }
}
