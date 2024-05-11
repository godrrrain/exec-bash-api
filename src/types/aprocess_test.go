package types

import (
	"os/exec"
	"testing"
)

func TestActiveProcesses_Add(t *testing.T) {
	ap := NewActiveProcesses()

	cmd := exec.Command("echo", "test")
	id := "328a8031-e9f5-4de4-9345-0e0a73fc8381"

	ap.Add(id, cmd)

	if _, ok := ap.Load(id); !ok {
		t.Errorf("Expected to find command with id %s, but it was not found", id)
	}
}

func TestActiveProcesses_Delete(t *testing.T) {
	ap := NewActiveProcesses()

	cmd := exec.Command("echo", "test")
	id := "328a8031-e9f5-4de4-9345-0e0a73fc8381"

	ap.Add(id, cmd)
	ap.Delete(id)

	if _, ok := ap.Load(id); ok {
		t.Errorf("Expected to not find command with id %s after deletion, but it was found", id)
	}
}
