package http

import (
	"net/http"

	"github.com/upframe/chatty"
)

type action struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type interactiveMessage struct {
	Actions    []action `json:"actions"`
	CallbackID string   `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTS    string `json:"action_ts"`
	MessageTS   string `json:"message_ts"`
	Token       string `json:"token"`
	ResponseURL string `json:"response_url"`
}

func interactive(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error) {

	return 0, nil
}
