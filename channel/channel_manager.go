package manager

import "sync"
import "time"
import "errors"

type ChManager struct {
	chMap   map[string]chan string
	rwmutex sync.RWMutex
}

func NewManager() *ChManager {
	return &ChManager{chMap: make(map[string]chan string)}
}

var UnkownIdError = errors.New("Unkown Id")
var ExistsIdError = errors.New("Exists Id")

func (cm *ChManager) Add(id string, ch chan string) error {

	cm.rwmutex.RLock()

	_, ok := cm.chMap[id]

	cm.rwmutex.RUnlock()

	if ok {
		return ExistsIdError
	}

	cm.rwmutex.Lock()

	cm.chMap[id] = ch

	cm.rwmutex.Unlock()

	return nil
}

func (cm *ChManager) Close(id string) {

	cm.rwmutex.Lock()
	ch, ok := cm.chMap[id]

	if ok {

		close(ch)

		delete(cm.chMap, id)
	}

	cm.rwmutex.Unlock()
}

func (cm *ChManager) Send(id string, msg string) error {

	var (
		ch chan string
		ok bool
	)
	cm.rwmutex.RLock()

	ch, ok = cm.chMap[id]

	cm.rwmutex.RUnlock()

	if !ok {
		return UnkownIdError
	}

	select {
	case ch <- msg:
	case <-time.After(time.Second):
		cm.Close(id)
	}

	return nil
}
