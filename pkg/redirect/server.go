package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flo80/redirect/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type responseStatus struct {
	Status  bool
	Message string
	Content []storage.Redirect
}

// Server settings for redirect server
type Server struct {
	listenAddress      string         // ip:port to listen on, for all interfaces empty, e.g. ":8080"
	adminHost          string         // hostname for administration of redirects (REST API at /redirects)
	storage.Redirector                // storage of all redirects: hostname, URL, target
	mux                *http.ServeMux // mux for handlers
	logger             *log.Logger    //logger to be used BUG: not yet implemented
}

// NewServer creates new server, sets handle functions but does not start listening.
// Per default it uses http.DefaultServeMux, a default ruslog logger and an empty MapRedirect for storage
func NewServer(listenAddress string, opts ...Option) *Server {
	redirector := storage.MapRedirect{}

	s := &Server{
		listenAddress: listenAddress,
		adminHost:     "",
		Redirector:    &redirector,
		mux:           http.DefaultServeMux,
		logger:        &log.Logger{},
	}

	for _, opt := range opts {
		opt(s)
	}

	s.mux.HandleFunc("/", s.Handler)

	if s.adminHost != "" {
		s.mux.HandleFunc(s.adminHost+"/redirects/ping", s.AdminAPI)
		s.mux.HandleFunc(s.adminHost+"/redirects/list", s.AdminAPI)
		s.mux.HandleFunc(s.adminHost+"/redirects/add", s.AdminAPI)
		s.mux.HandleFunc(s.adminHost+"/redirects/delete", s.AdminAPI)
		s.mux.HandleFunc(s.adminHost+"/redirects/deleteHost", s.AdminAPI)
	}
	return s
}

//StartServer calls http.ListenAndServe
func (s *Server) StartServer() error {
	log.Printf("Starting redirect server on address %v", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, s.mux)
}

// Option defines options to set for server
type Option func(*Server)

// WithAdmin allows to enable the REST API
func WithAdmin(adminHost string) Option {
	return func(s *Server) { s.adminHost = adminHost }
}

// WithLogger allows to pass a custom logger
func WithLogger(logger *log.Logger) Option {
	return func(s *Server) { s.logger = logger }
}

//WithMux allows to pass a custom mux
func WithMux(mux *http.ServeMux) Option {
	return func(s *Server) { s.mux = mux }
}

//WithRedirector allows to pass a custom redirector
func WithRedirector(redirector storage.Redirector) Option {
	return func(s *Server) { s.Redirector = redirector }
}

// Handler for http.HandleFunc for redirects
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	target, err := s.Redirector.GetTarget(r.Host, r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("no redirect found: %v", err)
		return
	}
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	log.Printf("request received for host %v and url %v, redirected to %v", r.Host, r.URL, target)

}

// AdminAPI is the http.Handler for API
// API supports following GET functions
//
//   /redirects/ping - only receive status ok
//   /redirects/list - list all redirects
//   /redirects/list?host=x - list all redirects for host x
//   /redirects/list?host=x&url=y - show redirect for host x with url y
//   /redirects/add?host=x&url=y&target=z - add or change redirect for host x with url y to target z
//   /redirects/delete?host=x&url=y - delete redirect for host x and url y
//   /redirects/deleteHost?host=x - delete all redirects for host x
//
// add, delete and deleteHost reply with a status
//   Status: true iftrue
//   Message: additional information
//   Content: []Redirect
func (s *Server) AdminAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	log.Debugf("received request %v", r)

	red := s.Redirector

	urlSplit := strings.Split(r.URL.Path, "/")
	if len(urlSplit) != 3 {
		http.NotFound(w, r)
		return
	}
	function := urlSplit[2]

	params := r.URL.Query()

	host := ""
	if hosts := params["host"]; len(hosts) > 0 {
		host = hosts[0]
	}

	url := ""
	if urls := params["url"]; len(urls) > 0 {
		url = urls[0]
	}

	target := ""
	if targets := params["target"]; len(targets) > 0 {
		target = targets[0]
	}

	var response responseStatus

	log.Debugf("parsed request %v %v %v %v", function, host, url, target)

	switch function {
	case "ping":
		response = responseStatus{true, "pong", nil}
	case "list":
		if host == "" {
			response = responseStatus{true, "all redirects", red.GetAllRedirects()}
		} else if url == "" {
			response = responseStatus{true, "redirects for host", red.GetRedirectsForHost(host)}
		} else {
			response = responseStatus{true, "redirects for host and url", red.GetRedirect(host, url)}
		}
	case "add":
		if host == "" || url == "" || target == "" {
			response = responseStatus{false, "request malformed", nil}
		} else {
			err := red.AddRedirect(storage.Redirect{host, url, target})
			if err != nil {
				response = responseStatus{false, err.Error(), nil}
			} else {
				response = responseStatus{true, "redirect added", red.GetRedirect(host, url)}
			}
		}
	case "delete":
		if host == "" || url == "" {
			response = responseStatus{false, "request malformed", nil}
		} else {
			red.RemoveRedirect(storage.Redirect{host, url, ""})
			response = responseStatus{true, "redirect deleted", nil}
		}
	case "deleteHost":
		if host == "" {
			response = responseStatus{false, "request malformed", nil}
		} else {
			red.RemoveAllRedirectsForHost(storage.Redirect{host, "", ""})
			response = responseStatus{true, "host deleted", nil}
		}
	default:
		http.NotFound(w, r)
		return
	}

	log.Debugf("sending response %v", response)
	json.NewEncoder(w).Encode(response)

}
