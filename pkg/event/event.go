package event

import (
    "errors"
    "github.com/docker/go-events"
    "reflect"
    "sync"
)

var eventMap sync.Map

var Listens map[interface{}][]events.Sink

func Event(fn interface{}) error {
    eName := reflect.TypeOf(fn).String()
    v, ok := eventMap.Load(eName)
    if !ok {
        return errors.New("map error")
    }
    b := v.(*events.Broadcaster)
    err := b.Write(fn)
    return err
}

func Register() error {
    for key, value := range Listens {
        eName := reflect.TypeOf(key).String()
        e := events.NewBroadcaster(value...)
        eventMap.Store(eName, e)
    }
    return nil
}
