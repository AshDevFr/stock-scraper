package state

import (
	"stock_scraper/types"
	"sync"
	"time"
)

var (
	once          sync.Once
	contentsMutex = &sync.Mutex{}
	alertsMutex   = &sync.Mutex{}
	alertInterval = 5
)

type State struct {
	alerts   map[string]*time.Time
	contents map[string]*types.Result
}

func NewState() *State {
	return &State{
		alerts:   make(map[string]*time.Time),
		contents: make(map[string]*types.Result),
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
		nextAlert := lastAlert.Add(time.Minute * time.Duration(alertInterval))

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

func GetContent(uuid string) *types.Result {
	state := GetState()
	contentsMutex.Lock()
	content := state.contents[uuid]
	contentsMutex.Unlock()
	return content
}

func SetContent(uuid string, result types.Result) {
	state := GetState()
	contentsMutex.Lock()
	state.contents[uuid] = &result
	contentsMutex.Unlock()
}
