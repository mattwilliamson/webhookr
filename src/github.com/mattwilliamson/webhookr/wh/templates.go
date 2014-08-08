package wh

import (
	"os"
	"html/template"
	"path/filepath"
)

// tPath takes a template name and gets the absolute path to it
func (s *Server) tPath(templatePath string) string {
	p, err := filepath.Abs(filepath.Join(s.StaticDir, "/templates/", templatePath))

	if err != nil {
		s.Log.Critical("Could not resolve path to template %v", err)
		os.Exit(1)
	}

	return p
}

// processTemplates loads the templates and caches them
func (s *Server) processTemplates() {
	bt := s.tPath("base.html")
	s.templates = make(map[string]*template.Template)
	s.templates["index.html"] = template.Must(template.ParseFiles(s.tPath("index.html"), bt))
	s.templates["help.html"] = template.Must(template.ParseFiles(s.tPath("help.html"), bt))
	s.templates["webhook.html"] = template.Must(template.ParseFiles(s.tPath("webhook.html"), bt))
	s.templates["posted.html"] = template.Must(template.ParseFiles(s.tPath("posted.html")))
}

