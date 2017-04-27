package http

import (
	"net/http"

	"github.com/nlopes/slack"

	"github.com/upframe/chatty"
)

func interactive(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error) {
	action := &slack.AttachmentActionCallback{}

	err := json.NewDecoder(r.Body).Decode(action)
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

	// TODO: finish proccessing

	team := c.Teams[action.Team.ID]
	team.Websites = append(team.Websites, action.Actions[0].Value)

	return 0, nil
}
