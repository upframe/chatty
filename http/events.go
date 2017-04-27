package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
	"github.com/nlopes/slack"

	"github.com/upframe/chatty"
)

type teste struct {
	Token     string        `json:"token"`
	Type      string        `json:"type"`
	Challenge string 	`json:"challenge"`
	Team      string	`json:"team_id"`
	APIAppID  string	`json:"api_app_id"`
	Event     *slack.Msg	`json:"event"`
	EventID   string        `json:"event_id"`
	EventTime int           `json:"event_time"`
}

func events(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error) {
	info := &teste{}

	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if info.Type == "url_verification" {
		w.Write([]byte(info.Challenge))
		return 0, nil
	}

	if info.Token != c.VerificationToken {
		return 0, nil
	}

	fmt.Println(info.Event)

	if _, ok := c.Teams[info.Team]; !ok {
		return 0, nil
	}

	team := c.Teams[info.Team]

	if info.Event.User == "" {
		return 0, nil
	}

	if _, ok := team.Connections[info.Event.ThreadTimestamp]; !ok {
		return 0, nil
	}

	conn := team.Connections[info.Event.ThreadTimestamp]
	err = websocket.Message.Send(conn, info.Event.Text)
	if err != nil {
		// log here
		fmt.Println("error")
	}

	return 0, nil
}
