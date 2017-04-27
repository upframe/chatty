package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/upframe/chatty"
)

// Serve ...
func Serve(c *chatty.Config) {
	c.Router = mux.NewRouter()

	c.Router.NotFoundHandler = &notFoundHandler{Config: c}
	// TODO: remove "/api" when Caddy fixes the proxy issue
	c.Router.HandleFunc("/api/setup", i(setup, c))
	c.Router.HandleFunc("/api/events", i(events, c))
	c.Router.HandleFunc("/api/interactive", i(interactive, c))
	c.Router.HandleFunc("/api/command", i(command, c))

	for _, team := range c.Teams {
		team.Start(c)
	}

	c.Logger.Infof("Listening on port %s.", "1562")
	if err := http.ListenAndServe(":1562", c.Router); err != nil {
		c.Logger.Fatal(err)
	}
}

type response struct {
	ID      string `json:"ID,omitempty"`
	Code    int
	Content interface{}
	Error   error `json:"-"`
}

type handler func(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error)

func i(h handler, c *chatty.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			code int
			err  error
		)

		defer func() {
			if code == 0 && err == nil {
				return
			}

			msg := &response{Code: code}

			if err != nil {
				msg.Content = err.Error()
			} else {
				msg.Content = http.StatusText(code)
			}

			if code >= 400 {
				t := time.Now()
				msg.ID = t.Format("20060102150405")
			}

			if code >= 400 && err != nil {
				c.Logger.Error(err)
			}

			if code != 0 {
				w.WriteHeader(code)
			}

			data, e := json.MarshalIndent(msg, "", "\t")
			if e != nil {
				c.Logger.Print(e)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}()

		code, err = h(w, r, c)
	}
}

type notFoundHandler struct {
	Config *chatty.Config
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i(func(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error) {
		w.WriteHeader(404)
		w.Write([]byte("no page for " + r.URL.String()))

		return 0, nil
	}, h.Config)(w, r)
}
