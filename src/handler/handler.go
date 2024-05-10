package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"

	"github.com/godrrrain/exec-bash-api/src/storage"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type CreateCommandRequest struct {
	Description string `json:"description"`
	Script      string `json:"script"`
}

type CreateCommandResponse struct {
	Uuid string `json:"uuid"`
}

type CommandResponse struct {
	Description string `json:"description"`
	Script      string `json:"script"`
	Status      string `json:"status"`
	Output      string `json:"output"`
}

type Handler struct {
	storage    storage.Storage
	aprocesses *storage.ActiveProcesses
}

func NewHandler(storage storage.Storage, aprocesses *storage.ActiveProcesses) *Handler {
	return &Handler{
		storage:    storage,
		aprocesses: aprocesses,
	}
}

func (h *Handler) CreateCommand(c *gin.Context) {

	var reqCommand CreateCommandRequest

	err := json.NewDecoder(c.Request.Body).Decode(&reqCommand)
	if err != nil {
		log.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid body",
		})
		return
	}

	uid := uuid.New()

	commandUuid := uid.String()
	status := "EXECUTING"

	cmd := exec.Command("/bin/sh", "-c", reqCommand.Script)
	// err = h.storage.CreateCommand(context.Background(), reqCommand.Description, reqCommand.Script, status, output)
	// if err != nil {
	// 	log.Printf("failed to create banner %s\n", err.Error())
	// 	c.JSON(http.StatusInternalServerError, ErrorResponse{
	// 		Error: err.Error(),
	// 	})
	// 	return
	// }

	go func() {
		output, err := cmd.CombinedOutput()
		status = "EXECUTED"
		if err != nil {
			log.Printf("Ошибка выполнения скрипта %s\n", err.Error())
			status = "FAILED"
		}

		log.Println(string(output))

		err = h.storage.CreateCommand(context.Background(), commandUuid, reqCommand.Description, reqCommand.Script, status, string(output))
		if err != nil {
			log.Printf("failed to create banner %s\n", err.Error())
		}
	}()

	c.JSON(http.StatusCreated, CreateCommandResponse{
		Uuid: commandUuid,
	})
}

func (h *Handler) GetCommand(c *gin.Context) {
	command_uuid := c.Param("uuid")

	if !IsValidUUID(command_uuid) {
		c.Status(http.StatusNotFound)
		return
	}

	command, err := h.storage.GetCommand(context.Background(), command_uuid)
	if err != nil {
		if err.Error() == "banner not found" {
			log.Println(err.Error())
			c.Status(http.StatusNotFound)
			return
		}
		log.Printf("failed to get command: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CommandToResponse(command))
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func CommandToResponse(command storage.Command) CommandResponse {
	return CommandResponse{
		Description: command.Description,
		Script:      command.Script,
		Status:      command.Status,
		Output:      command.Output,
	}
}
