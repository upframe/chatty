package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nlopes/slack"

	"github.com/upframe/chatty"
)

// NOTE: if in the future we start having problems with non-received messages
// we must start using action.ResponseURL to send instead of writing directly.

func interactive(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error) {
	action := &slack.AttachmentActionCallback{}

	err := json.Unmarshal([]byte(r.PostFormValue("payload")), action)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if action.Token != c.VerificationToken {
		return 0, nil
	}

	if _, ok := c.Teams[action.Team.ID]; !ok {
		return 0, nil
	}

	if action.CallbackID != "accept_channel" {
		return 0, nil
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	if len(action.Actions) == 0 {
		return http.StatusBadRequest, errors.New("missing action")
	}

	if action.Actions[0].Value == "" {
		w.Write([]byte("You've rejected this website approval request."))
		return 0, nil
	}

	team := c.Teams[action.Team.ID]
	team.Websites = append(team.Websites, action.Actions[0].Value)

	w.Write([]byte("You've approved this website approval request."))
	return 0, nil
}
