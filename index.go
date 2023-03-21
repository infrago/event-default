package event_default

import (
	"github.com/infrago/event"
	"github.com/infrago/infra"
)

func Driver() event.Driver {
	return &defaultDriver{}
}

func init() {
	infra.Register("default", Driver())
}
