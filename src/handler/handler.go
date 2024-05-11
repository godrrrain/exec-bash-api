package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/godrrrain/exec-bash-api/src/storage"
	"github.com/godrrrain/exec-bash-api/src/types"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage    storage.Storage
	aprocesses *types.ActiveProcesses
}

func NewHandler(storage storage.Storage, aprocesses *types.ActiveProcesses) *Handler {
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
	h.aprocesses.Add(commandUuid, cmd)
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
			log.Printf("script execution error %s\n", err.Error())
			status = "FAILED"
		}

		h.aprocesses.Delete(commandUuid)

		log.Println(string(output))

		err = h.storage.CreateCommand(context.Background(), commandUuid, reqCommand.Description, reqCommand.Script, status, string(output))
		if err != nil {
			log.Printf("failed to create command %s\n", err.Error())
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
		if err.Error() == "command not found" {
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

func (h *Handler) StopCommand(c *gin.Context) {
	command_uuid := c.Param("uuid")

	aprocess, ok := h.aprocesses.Load(command_uuid)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	if aprocess.Cmd.Process != nil {
		err := aprocess.Cmd.Process.Kill()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "error while stopping command",
			})
			return
		}
	}

	h.aprocesses.Delete(command_uuid)

	// err := h.storage.UpdateCommandStatus(context.Background(), command_uuid, "STOPPED")
	// if err != nil {
	// 	if err.Error() == "command not found" {
	// 		log.Println(err.Error())
	// 		c.Status(http.StatusNotFound)
	// 		return
	// 	}
	// 	log.Printf("failed to get command: %s\n", err.Error())
	// 	c.JSON(http.StatusInternalServerError, ErrorResponse{
	// 		Error: err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, MessageResponse{
		Message: "command stopped successfully",
	})
}
