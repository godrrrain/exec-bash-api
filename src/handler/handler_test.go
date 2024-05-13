package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/godrrrain/exec-bash-api/src/storage"
	"github.com/godrrrain/exec-bash-api/src/types"
)

func TestHandler_CreateCommand(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	aprocesses := types.NewActiveProcesses()

	handler := NewHandler(mockStorage, aprocesses)
	r := gin.New()
	r.POST("/api/v1/commands", handler.CreateCommand)

	// ok
	testBody := `{"description" : "Простой скрипт с циклом 8", "script": "echo Hyo; sleep 2"}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/commands", bytes.NewBufferString(testBody))

	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Test ok. Unexpected status code, expected %d, got %d instead", 201, w.Code)
	}
}

func TestHandler_GetCommand(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	aprocesses := types.NewActiveProcesses()

	handler := NewHandler(mockStorage, aprocesses)
	r := gin.New()
	r.GET("/api/v1/commands/:uuid", handler.GetCommand)

	// ok
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/commands/aba5be85-f261-451f-a279-a3b0bbb6d1c8", nil)

	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Test ok. Unexpected status code, expected %d, got %d instead", 200, w.Code)
	}

	// command not found
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/commands/210c0ac1-a2ae-4600-bd04-7cd8ada90e7d", nil)

	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Test command not found. Unexpected status code, expected %d, got %d instead", 404, w.Code)
	}

	// incorrect uuid
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/commands/10", nil)

	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Test incorrect uuid. Unexpected status code, expected %d, got %d instead", 404, w.Code)
	}
}

func TestHandler_GetCommands(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	aprocesses := types.NewActiveProcesses()

	handler := NewHandler(mockStorage, aprocesses)
	r := gin.New()
	r.GET("/api/v1/commands", handler.GetCommands)

	// ok
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/commands", nil)

	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Test ok. Unexpected status code, expected %d, got %d instead", 200, w.Code)
	}

	// ok with parametrs
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/commands?status=EXECUTING&limit=5&offset=1", nil)

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Test ok with parametrs. Unexpected status code, expected %d, got %d instead", 200, w.Code)
	}

	// incorrect limit
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/commands?limit=seven", nil)

	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Test incorrect limit. Unexpected status code, expected %d, got %d instead", 400, w.Code)
	}

	// incorrect offset
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/commands?offset=five", nil)

	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Test incorrect offset. Unexpected status code, expected %d, got %d instead", 400, w.Code)
	}
}
