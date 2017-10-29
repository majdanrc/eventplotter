package providers

import (
	"github.com/majdanrc/eventplotter/events"
	"github.com/majdanrc/eventplotter/streamer"
)

type EventProvider interface {
	ProvideEvents(streamer.Stream) ([]events.Event, error)
}
