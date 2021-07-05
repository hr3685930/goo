package events

import (
	"fmt"
	"github.com/docker/go-events"
)

type UserDel struct {}

func (e *UserDel)Write(event events.Event) error {
	fmt.Println(1111)
	return nil
}

func (e *UserDel)Close() error {
	return nil
}
