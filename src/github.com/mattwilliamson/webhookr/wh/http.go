package wh

import (
	"net/http"
	"regexp"
)

var webhookrRe = regexp.MustCompile(`/(\w+)(/.+)`)

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request, t string, c map[string]string) {
    err := s.templates[t].ExecuteTemplate(w, "base", c)

    if err != nil {
    	s.Log.Error("Error rendering template %v", err)
    }
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	s.Log.Debug("RequestURI: %v", r.RequestURI)

	if len(r.RequestURI) > 1 {
		s.webhookrHandler(w, r)
	} else {
		s.indexHandler(w, r)
	}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	c := map[string]string {
		"STATIC_URL": s.StaticUrl.String(),
		"nav_home_class": "active",
	}
	s.renderTemplate(w, r, "index.html", c)
}

func (s *Server) webhookrHandler(w http.ResponseWriter, r *http.Request) {
	ids := webhookrRe.FindStringSubmatch(r.RequestURI)
	id := ids[0]
	remainder := ids[1]
	shouldPost := len(remainder) > 1 || r.Method != "GET" || r.ContentLength > 0

	s.Log.Debug("webhookr path:%v id: %v", r.RequestURI, id)
	s.Log.Debug("webhookr remainder:%v id: %v", r.RequestURI, remainder)

	if shouldPost {
		s.Log.Debug("This is a POST")
		
	} else {
		c := map[string]string {
			"STATIC_URL": s.StaticUrl.String(),
			"Path": r.URL.String(),
			"ID": id,
		}
		s.renderTemplate(w, r, "webhook.html", c)
	}
}

func (s *Server) helpHandler(w http.ResponseWriter, r *http.Request) {
	c := map[string]string {
		"STATIC_URL": s.StaticUrl.String(),
		"nav_help_class": "active",
	}
    s.renderTemplate(w, r, "help.html", c)
}

func (s *Server) newHookrHandler(w http.ResponseWriter, r *http.Request) {
	id := RandomId()
	http.Redirect(w, r, "/"+id, 302)
}



