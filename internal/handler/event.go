package handler

import (
	"goo/internal/events"
	"goo/internal/types"
	"goo/pkg/event"
	e "github.com/docker/go-events"

)

type Event struct {
}

func (*Event) Handler() {
	event.Listens = map[interface{}][]e.Sink{
		&types.UserSaveEvent{}: {
			&events.UserDel{},
			&events.UserSave{},
		},
	}
}