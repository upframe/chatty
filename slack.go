package chatty

// Message ...
type Message struct {
	OK              bool   `json:"ok,omitempty"`
	Channel         string `json:"channel,omitempty"`
	Timestamp       string `json:"ts,omitempty"`
	ThreadTimestamp string `json:"thread_ts,omitempty"`
	Text            string `json:"text,omitempty"`
	User            string `json:"user,omitempty"`
}
