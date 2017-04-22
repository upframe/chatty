package upframy

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/nlopes/slack"
)

// pingPong sends a ping/pong message using the same case of the sent one.
// e.g.: "ping" -> "pong"; "PiNg" -> "PoNg"; and so on.
func pingPong(message *slack.MessageEvent) {
	answer := message.Text
	var replacement string

	switch answer[1] {
	case 'I':
		replacement = "O"
	case 'i':
		replacement = "o"
	case 'O':
		replacement = "I"
	case 'o':
		replacement = "i"
	default:
		// Wut, what? This must not happen!
		replacement = "(You found a bug)"
	}

	answer = answer[:1] + replacement + answer[2:]

	if message.Channel[0] != 'D' {
		answer = "<@" + message.User + "> " + answer
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(answer, message.Channel))
}

// icndb is the Internet Chuck Norris Database API response type
type icndb struct {
	Type  string `json:"type"`
	Value struct {
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
		Categories []string `json:"categories"`
	} `json:"value"`
}

// makeFunOfUser creates a joke using Internet Chuck Norris Database about
// the user!
func makeFunOfUser(message *slack.MessageEvent) {
	firstName := url.QueryEscape(users[message.User].FirstName)
	lastName := url.QueryEscape(users[message.User].LastName)

	link := "https://api.icndb.com/jokes/random?escape=javascript&firstName=" + firstName + "lastName=" + lastName

	r, err := http.Get(link)
	if err != nil {
		logger.Printf("Error on makeFunOfUser: %v", err)
		return
	}
	defer r.Body.Close()

	resp := &icndb{}
	if err = json.NewDecoder(r.Body).Decode(resp); err != nil {
		logger.Printf("Error on makeFunOfUser: %v", err)
		return
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(resp.Value.Joke, message.Channel))
}
