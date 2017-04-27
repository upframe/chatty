package chatty

import (
	"net/http"
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

	// TODO: if channel == "", send a message to choose the channel

	c.Router.Handle(t.WebsocketURI, websocket.Handler(t.websocket))
	c.Logger.Infof("Successfully started team %s on %s", t.ID, t.WebsocketURI)
}

func (t *Team) allowed(url string) bool {
	for _, website := range t.Websites {
		if website == url {
			return true
		}
	}

	// If it's not allowed, let's notify our dear team that
	// someone from the outside tried to access the websocket
	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			{
				Title:      "Website Approval",
				Text:       "Do you wish to add <" + url + "> to the verified websites' list?",
				Fallback:   "You are unable to either reject or accept.",
				CallbackID: "accept_channel",
				Color:      "#8BC34A",
				Actions: []slack.AttachmentAction{
					{
						Name:  "approve",
						Text:  "Approve",
						Type:  "button",
						Style: "primary",
						Value: "true",
						Confirm: &slack.ConfirmationField{
							Title:       "Are you sure?",
							Text:        "This will let *" + url + "* access your chat from now on.",
							OkText:      "Yes",
							DismissText: "No",
						},
					},
					{
						Name:  "approve",
						Text:  "Reject",
						Type:  "button",
						Value: "false",
						Confirm: &slack.ConfirmationField{
							Title:       "Are you sure?",
							Text:        "This will not let *" + url + "* access your chat.",
							OkText:      "Yes",
							DismissText: "No",
						},
					},
				},
			},
		},
	}

	text := "Someone tried to access your API websocket from a new URL. If it was you, click the button bellow to accept it and add this to the list of verified websites."
	_, _, err := t.Slack.PostMessage(t.Channel, text, params)
	if err != nil {
		t.Logger.Errorf("can't post '%s' message for '%s' team", text, t.ID)
	}

	return false
}

// websocket is the handler for the socket that each team has.
func (t *Team) websocket(ws *websocket.Conn) {
	if !t.allowed(ws.RemoteAddr().String()) {
		err := ws.WriteClose(http.StatusForbidden)
		if err != nil {
			t.Logger.Errorf("error closing the websocket %v", err)
		}
	}

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
			t.Logger.Errorf("can't post '%s' message for '%s' team", reply, t.ID)
			break
		}

		if thread == "" {
			thread = ts
			t.Connections[thread] = ws
		}
	}

	if thread == "" {
		return
	}

	params := slack.PostMessageParameters{ThreadTimestamp: thread}
	_, _, err = t.Slack.PostMessage(t.Channel, exitMessage, params)
	if err != nil {
		t.Logger.Errorf("can't post '%s' message for '%s' team", exitMessage, t.ID)
	}
}
