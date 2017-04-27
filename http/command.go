package http

import (
	"net/http"
	"strings"

	"github.com/nlopes/slack"
	"github.com/upframe/chatty"
)

// NOTE: if in the future we start having problems with non-received messages
// we must start using ResponseURL to send instead of writing directly.

type commandEvent struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	UserID      string
	UserName    string
	Command     string
	Text        string
	ResponseURL string
}

func command(w http.ResponseWriter, r *http.Request, c *chatty.Config) (int, error) {
	cmd := &commandEvent{
		Token:       r.FormValue("token"),
		TeamID:      r.FormValue("team_id"),
		TeamDomain:  r.FormValue("team_domain"),
		ChannelID:   r.FormValue("channel_id"),
		ChannelName: r.FormValue("channel_name"),
		UserID:      r.FormValue("user_id"),
		UserName:    r.FormValue("user_name"),
		Command:     r.FormValue("command"),
		Text:        r.FormValue("text"),
		ResponseURL: r.FormValue("response_url"),
	}

	if cmd.Token != c.VerificationToken {
		return 0, nil
	}

	if _, ok := c.Teams[cmd.TeamID]; !ok {
		return 0, nil
	}

	if cmd.Command != "/chatty" {
		return 0, nil
	}

	team := c.Teams[cmd.TeamID]

	cmd.Text = strings.TrimSpace(cmd.Text)

	switch {
	case cmd.Text == "help":
		w.Write([]byte("This is your help"))
	case cmd.Text == "websites":
		w.Write(websitesResponse(team))
	case strings.HasPrefix(cmd.Text, "channel"):
		w.Write(channelResponse(team, cmd.Text))
	}

	return 0, nil
}

func websitesResponse(team *chatty.Team) []byte {
	if len(team.Websites) == 0 {
		return []byte("Your team hasn't approved any website yet.")
	}

	response := "Here is the list of approved websites:"

	for i := range team.Websites {
		response += "\n- " + team.Websites[i]
	}

	return []byte(response)
}

func channelResponse(team *chatty.Team, message string) []byte {
	if message == "channel" {
		if team.Channel == "" {
			return []byte("Your team hasn't set the channel I'm going to use. To do so, you can use \"/chatty channel #myChannel\".")
		}

		return []byte("Currently, Chatty is configured to use <#" + team.Channel + ">.")
	}

	args := strings.Split(message, " ")
	if len(args) > 2 {
		return []byte("You sent more stuff than I thought you would! The correct syntax is \"/chatty channel #myChannel\".")
	}

	channel := args[1]
	if !strings.HasPrefix(channel, "<#") || !strings.HasSuffix(channel, ">") {
		return []byte("That is not a channel name!")
	}

	channel = strings.TrimPrefix(channel, "<#")
	channel = strings.TrimSuffix(channel, ">")
	team.Channel = channel

	// Ignore err cause this is not much important now
	team.Slack.PostMessage(channel, "This is my new workplace! So sweet :blush:", slack.PostMessageParameters{})
	return []byte("Channel <#" + channel + "> successfully set!")
}
