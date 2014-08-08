package wh

import (
    "fmt"
    "net/http"
	"net/url"
	"html/template"
	"github.com/mattwilliamson/log"
)


type Server struct {
	StaticUrl *url.URL
	StaticDir string
	Host string
	Port uint
	*log.Log
	templates map[string]*template.Template
    Registry *ClientRegistry
}

func (s *Server) ListenAndServe() error {
	log := s.Log

	// Setup templates
	s.processTemplates()

    // Setup file server
    if !s.StaticUrl.IsAbs() {
    	staticUrl := s.StaticUrl.String()
    	log.Info("Hosting static files from %v at %v", s.StaticDir, staticUrl)
    	http.Handle(staticUrl, http.StripPrefix(staticUrl, http.FileServer(http.Dir(s.StaticDir))))
    }

    // Bind handlers to urls
    http.HandleFunc("/", s.rootHandler)
    http.HandleFunc("/new", s.newHookrHandler)
    http.HandleFunc("/help", s.helpHandler)

    // Start Listening
    bindAddress := fmt.Sprintf("%v:%d", s.Host, s.Port)
    log.Info("Listening on %v...", bindAddress)
    err := http.ListenAndServe(bindAddress, nil)

    return err
}

func New() *Server {
	server := &Server{}
    server.Registry = &NewReigstry()

    return server
}