package http

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/upframe/chatty"
)

// setup creates a new team in our database and sets up its connection.
func setup(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error) {
	resp, err := slack.GetOAuthResponse(c.ClientID, c.ClientSecret, r.URL.Query().Get("code"), "", false)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !resp.Ok {
		return http.StatusInternalServerError, errors.New("Setup is not OK")
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	hash := md5.Sum([]byte(resp.TeamID))

	team := &chatty.Team{
		ID:           resp.TeamID,
		Name:         resp.TeamName,
		Token:        resp.AccessToken,
		Channel:      "C50SNCVSR",
		WebsocketURI: "/websocket/" + hex.EncodeToString(hash[:]),
	}

	c.Teams[team.ID] = team
	if c.TeamService.Save(team) != nil {
		return http.StatusInternalServerError, err
	}

	c.Logger.Infof("Successfully authorized team %s", resp.TeamID)
	c.Teams[team.ID].Start(c)
	return 0, nil
}
