package upframy

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"time"

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

	link := "https://api.icndb.com/jokes/random?escape=javascript&firstName=" + firstName + "&lastName=" + lastName

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

type wdtt struct {
	Message string `json:"message"`
}

func whatDoesTrumpThink(message *slack.MessageEvent) {
	link := "https://api.whatdoestrumpthink.com/api/v1/quotes/random"

	r, err := http.Get(link)
	if err != nil {
		logger.Printf("Error on whatDoesTrumpThink: %v", err)
		return
	}
	defer r.Body.Close()

	resp := &wdtt{}
	if err = json.NewDecoder(r.Body).Decode(resp); err != nil {
		logger.Printf("Error on whatDoesTrumpThink: %v", err)
		return
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(resp.Message, message.Channel))
}

var foxSayings = []string{
	`"Ring-ding-ding-ding-dingeringeding!
Gering-ding-ding-ding-dingeringeding!
Gering-ding-ding-ding-dingeringeding!"`,
	`"Wa-pa-pa-pa-pa-pa-pow!
Wa-pa-pa-pa-pa-pa-pow!
Wa-pa-pa-pa-pa-pa-pow!"`,
	`"Hatee-hatee-hatee-ho!
Hatee-hatee-hatee-ho!
Hatee-hatee-hatee-ho!"`,
	`"Joff-tchoff-tchoffo-tchoffo-tchoff!
Tchoff-tchoff-tchoffo-tchoffo-tchoff!
Joff-tchoff-tchoffo-tchoffo-tchoff!"`,
	`"Jacha-chacha-chacha-chow!
Chacha-chacha-chacha-chow!
Chacha-chacha-chacha-chow!"`,
	`"Fraka-kaka-kaka-kaka-kow!
Fraka-kaka-kaka-kaka-kow!
Fraka-kaka-kaka-kaka-kow!"`,
	`"A-hee-ahee ha-hee!
A-hee-ahee ha-hee!
A-hee-ahee ha-hee!"`,
	`"A-oo-oo-oo-ooo!
Woo-oo-oo-ooo!"`,
}

func whatDoesTheFoxSay(message *slack.MessageEvent) {
	rand.Seed(time.Now().Unix())
	id := rand.Intn(len(foxSayings) - 1)

	msg := rtm.NewOutgoingMessage(foxSayings[id], message.Channel)
	if message.ThreadTimestamp != "" {
		msg.ThreadTimestamp = message.ThreadTimestamp
	}

	rtm.SendMessage(msg)
}
