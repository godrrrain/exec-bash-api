package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/godrrrain/exec-bash-api/src/storage"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) GetCommands(c *gin.Context) {

	params := c.Request.URL.Query()

	limitStr := params.Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("invalid limit query %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid limit query",
		})
		return
	}

	offsetStr := params.Get("offset")
	if offsetStr == "" {
		offsetStr = "-1"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Printf("invalid offset query %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid offset query",
		})
		return
	}

	status := params.Get("status")

	commands, err := h.storage.GetCommands(context.Background(), status, limit, offset)
	if err != nil {
		log.Printf("failed to get commands %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CommandsToResponse(commands))
}
