package chatty

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"

	"golang.org/x/net/websocket"
)

const exitMessage = "*My fan just left.*"

// Config contains the configuration of the app.
type Config struct {
	Mu                *sync.Mutex
	Logger            *logrus.Logger
	Router            *mux.Router
	Teams             map[string]*Team
	TeamService       TeamService
	ClientID          string
	ClientSecret      string
	VerificationToken string
}

// Team contains the fields related to the teams connected
// to this application.
type Team struct {
	ID           string                     `json:"id"`
	Name         string                     `json:"name"`
	Scope        string                     `json:"scope"`
	Channel      string                     `json:"channel"`
	Websites     []string                   `json:"websites"`
	WebsocketURI string                     `json:"websocket_uri"`
	Token        string                     `json:"token"`
	Connections  map[string]*websocket.Conn `json:"-"`
	Slack        *slack.Client              `json:"-"`
	Logger       *logrus.Logger             `json:"-"`
}

// TeamService allows the management of the team object.
type TeamService interface {
	Loads() ([]*Team, error)
	Load(id string) (*Team, error)
	Save(t *Team) error
	Delete(t *Team) error
}

// Start starts the program for this team.
func (t *Team) Start(c *Config) {
	t.Slack = slack.New(t.Token)
	t.Logger = c.Logger
	t.Connections = make(map[string]*websocket.Conn)
	c.Router.Handle(t.WebsocketURI, websocket.Handler(t.websocket))
	c.Logger.Infof("Successfully started team %s on %s", t.ID, t.WebsocketURI)
}

// websocket is the handler for the socket that each team has.
func (t *Team) websocket(ws *websocket.Conn) {
	var err error
	var thread string

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			// This means the connection was closed.
			break
		}

		params := slack.PostMessageParameters{ThreadTimestamp: thread}
		_, ts, err := t.Slack.PostMessage(t.Channel, reply, params)
		if err != nil {
			t.Logger.Errorf("Can't post '%s' message for '%s' team", reply, t.ID)
			break
		}

		if thread == "" {
			thread = ts
			t.Connections[thread] = ws
		}
	}

	params := slack.PostMessageParameters{ThreadTimestamp: thread}
	_, _, err = t.Slack.PostMessage(t.Channel, exitMessage, params)
	if err != nil {
		t.Logger.Errorf("Can't post '%s' message for '%s' team", exitMessage, t.ID)
	}
}
