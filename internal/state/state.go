package state

import (
	"sync"
	"time"
)

var once sync.Once
var alertsMutex = &sync.Mutex{}
var contentsMutex = &sync.Mutex{}

type State struct {
	alerts   map[string]*time.Time
	contents map[string]*string
}

func NewState() *State {
	return &State{
		alerts:   make(map[string]*time.Time),
		contents: make(map[string]*string),
	}
}

var instance *State

func GetState() *State {
	once.Do(func() {
		instance = NewState()
	})

	return instance
}

func ShouldRunAlert(uuid string, callback func()) {
	state := GetState()
	alertsMutex.Lock()
	lastAlert := state.alerts[uuid]
	alertsMutex.Unlock()

	if lastAlert != nil {
		nextAlert := lastAlert.Add(time.Minute * 1)

		if nextAlert.After(time.Now()) {
			return
		}
	}
	now := time.Now()

	alertsMutex.Lock()
	state.alerts[uuid] = &now
	alertsMutex.Unlock()

	go callback()
}

func GetContent(uuid string) *string {
	state := GetState()
	contentsMutex.Lock()
	content := state.contents[uuid]
	contentsMutex.Unlock()
	return content
}

func SetContent(uuid string, content string) {
	state := GetState()
	contentsMutex.Lock()
	state.contents[uuid] = &content
	contentsMutex.Unlock()
}
