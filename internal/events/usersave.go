package events

import (
	"fmt"
	"github.com/docker/go-events"
)

type UserSave struct {}

func (e *UserSave)Write(event events.Event) error {
	fmt.Println(event)
	return nil
}

func (e *UserSave)Close() error {
	return nil
}
