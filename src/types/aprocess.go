package types

import (
	"os/exec"
	"sync"
)

type ActiveProcess struct {
	Cmd *exec.Cmd
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

func (ap *ActiveProcesses) Add(id string, cmd *exec.Cmd) {
	command := &ActiveProcess{
		Cmd: cmd,
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
