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

	c.Router.HandleFunc("/setup", i(setup, c))
	c.Router.HandleFunc("/events", i(events, c))
	c.Router.HandleFunc("/interactive", i(interactive, c))
	c.Router.HandleFunc("/command", i(command, c))

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
