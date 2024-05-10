package storage

import (
	"os"
	"sync"
)

type ActiveProcess struct {
	Process *os.Process
}

type ActiveProcesses struct {
	sync.RWMutex
	m map[string]*ActiveProcess
}

func NewActiveProcesses() *ActiveProcesses {
	return &ActiveProcesses{
		m: make(map[string]*ActiveProcess),
	}
}

func (ap *ActiveProcesses) Add(id string, process *os.Process) {
	command := &ActiveProcess{
		Process: process,
	}
	ap.Lock()
	ap.m[id] = command
	ap.Unlock()
}

func (ap *ActiveProcesses) Delete(id string) {
	ap.Lock()
	delete(ap.m, id)
	ap.Unlock()
}

func (ap *ActiveProcesses) Load(id string) (*ActiveProcess, bool) {
	ap.RLock()
	command, ok := ap.m[id]
	ap.RUnlock()
	return command, ok
}
