package upframy

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

var (
	rtm    *slack.RTM
	bot    string
	logger *log.Logger

	users = make(map[string]slack.UserProfile)
)

// Start ...
func Start(token string, l *log.Logger) {
	logger = l

	api := slack.New(token)
	slack.SetLogger(l)

	rtm = api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			logger.Printf("Logged in as %s\n", ev.Info.User.Name)

			for _, u := range ev.Info.Users {
				users[u.ID] = u.Profile
			}

			bot = ev.Info.User.ID
			go startChat()
		case *slack.MessageEvent:
			handleMessageEvent(ev)
		case *slack.RTMError:
			logger.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			logger.Printf("Invalid credentials")
			return
		case *slack.AckMessage:
			handleAcknowledgements(ev)
		}
	}
}

func handleMessageEvent(ev *slack.MessageEvent) {
	// Ignore messages from the bot itself or messages
	// with subtypes.
	if ev.User == bot || ev.SubType != "" {
		return
	}

	if ev.Channel == chatChannel {
		// If this handler runs, the others musn't run.
		if chatSlackHandler(ev) {
			return
		}
	}

	text := strings.ToLower(ev.Text)
	if strings.Contains(text, "miguel gonçalves") {
		reply(ev, `I'm not slackbot, but I know that if "Miguel Gonçalves" is mentioned again, this channel will be terminated!`)
		return
	}

	isDirectMessage := (ev.Channel[0] == 'D')
	isBotMentioned := strings.Contains(ev.Text, "<@"+bot+">")

	// Ignore messages that aren't direct messages and the bot
	// isn't mentioned.
	if !isDirectMessage && !isBotMentioned {
		return
	}

	// If it's mentioned, remove the mention so we can parse
	// the text later.
	if isBotMentioned {
		ev.Text = strings.Replace(ev.Text, "<@"+bot+">", "", -1)
		ev.Text = strings.TrimSpace(ev.Text)
	}

	switch {
	case text == "ping", text == "pong":
		pingPong(ev)
		return
	case strings.Contains(text, "tell me a joke"):
		makeFunOfUser(ev)
		return
	case strings.Contains(text, "what does trump think"):
		whatDoesTrumpThink(ev)
		return
	case strings.Contains(text, "what does the fox say"):
		whatDoesTheFoxSay(ev)
		return
	}

	var answer string

	switch text {
	case "fuck you":
		answer = "Don't be evil :scream:"
	case "hey", "hi", "hello", "hullo":
		answer = "Hi there!"
	case "bye", "cya", "goodbye":
		answer = "Bye! Gonna miss you :kissing:"
	default:
		if isDirectMessage {
			answer = fmt.Sprintf("Sorry %s, I didn't quite understand what you just said :disappointed:", users[ev.User].FirstName)
		} else {
			answer = fmt.Sprintf("If you mention me again, I will break your nose <@%s>!", ev.User)
		}
	}

	reply(ev, answer)
}

var (
	pendingTimestamps = []*pendingTimestamp{}
)

type pendingTimestamp struct {
	ID   int
	Chan chan string
}

func handleAcknowledgements(ack *slack.AckMessage) {
	for i := range pendingTimestamps {
		if ack.ReplyTo != pendingTimestamps[i].ID {
			continue
		}

		pendingTimestamps[i].Chan <- ack.Timestamp

		// Remove from pendingTimestamps
		pendingTimestamps = append(pendingTimestamps[:i], pendingTimestamps[i+1:]...)
	}
}
