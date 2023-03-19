package event_default

import "github.com/infrago/event"

func Driver() event.Driver {
	return &defaultDriver{}
}

func init() {
	event.Register("default", Driver())
}
