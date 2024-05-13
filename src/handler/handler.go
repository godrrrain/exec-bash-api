package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"sync"

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
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid body"))
		return
	}

	if reqCommand.Script == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("script must be given"))
		return
	}

	uid := uuid.New()

	commandUuid := uid.String()
	status := "EXECUTING"

	cmd := exec.Command("/bin/sh", "-c", reqCommand.Script)
	h.aprocesses.Add(commandUuid, cmd)

	defaultOutput := ""

	err = h.storage.CreateCommand(context.Background(), commandUuid, reqCommand.Description, reqCommand.Script, status, defaultOutput)
	if err != nil {
		log.Printf("failed to create command %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	go func() {
		output, err := cmd.CombinedOutput()
		status = "EXECUTED"
		if err != nil {
			log.Printf("script execution error %s\n", err.Error())
			status = "FAILED"

			if err.Error() == "signal: killed" {
				status = "STOPPED"
			}
		}

		h.aprocesses.Delete(commandUuid)

		err = h.storage.UpdateCommandOutput(context.Background(), commandUuid, string(output))
		if err != nil {
			log.Println(err.Error())
		}

		err = h.storage.UpdateCommandStatus(context.Background(), commandUuid, status)
		if err != nil {
			log.Println(err.Error())
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
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
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
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid limit query"))
		return
	}

	offsetStr := params.Get("offset")
	if offsetStr == "" {
		offsetStr = "-1"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Printf("invalid offset query %s\n", err.Error())
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid offset query"))
		return
	}

	status := params.Get("status")

	commands, err := h.storage.GetCommands(context.Background(), status, limit, offset)
	if err != nil {
		log.Printf("failed to get commands %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, CommandsToResponse(commands))
}

func (h *Handler) StopCommand(c *gin.Context) {
	commandUuid := c.Param("uuid")

	aprocess, ok := h.aprocesses.Load(commandUuid)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	if aprocess.Cmd.Process != nil {
		err := aprocess.Cmd.Process.Kill()
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse("error while stopping command"))
			return
		}
	}

	h.aprocesses.Delete(commandUuid)

	c.JSON(http.StatusOK, NewMessageResponse("command stopped successfully"))
}

func (h *Handler) DeleteCommand(c *gin.Context) {
	command_uuid := c.Param("uuid")

	if !IsValidUUID(command_uuid) {
		c.Status(http.StatusNotFound)
		return
	}

	err := h.storage.DeleteCommand(context.Background(), command_uuid)
	if err != nil {
		if err.Error() == "command not found" {
			log.Println(err.Error())
			c.Status(http.StatusNotFound)
			return
		}
		log.Printf("failed to delete command: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, NewMessageResponse("command deleted successfully"))
}

func (h *Handler) CreateDurableCommand(c *gin.Context) {

	var reqCommand CreateCommandRequest

	err := json.NewDecoder(c.Request.Body).Decode(&reqCommand)
	if err != nil {
		log.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid body"))
		return
	}

	if reqCommand.Script == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("script must be given"))
		return
	}

	uid := uuid.New()

	commandUuid := uid.String()
	status := "EXECUTING"

	cmd := exec.Command("/bin/sh", "-c", reqCommand.Script)
	h.aprocesses.Add(commandUuid, cmd)

	defaultOutput := ""

	err = h.storage.CreateCommand(context.Background(), commandUuid, reqCommand.Description, reqCommand.Script, status, defaultOutput)
	if err != nil {
		log.Printf("failed to create banner %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("error while creating pipe for stdout"))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("error while creating pipe for stderr"))
		return
	}

	if err := cmd.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("error while starting command"))
		return
	}

	var buffer bytes.Buffer
	var mutex sync.Mutex
	var tempOutput string

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			mutex.Lock()
			buffer.WriteString(scanner.Text() + "\n")
			tempOutput = buffer.String()
			mutex.Unlock()

			err = h.storage.UpdateCommandOutput(context.Background(), commandUuid, tempOutput)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			mutex.Lock()
			buffer.WriteString("ERROR: " + scanner.Text() + "\n")
			tempOutput = buffer.String()
			mutex.Unlock()

			err = h.storage.UpdateCommandOutput(context.Background(), commandUuid, tempOutput)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	go func() {
		err := cmd.Wait()
		status = "EXECUTED"
		if err != nil {
			log.Printf("script execution error: %s\n", err)
			status = "FAILED"

			if err.Error() == "signal: killed" {
				status = "STOPPED"
			}
		}

		h.aprocesses.Delete(commandUuid)

		err = h.storage.UpdateCommandStatus(context.Background(), commandUuid, status)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	c.JSON(http.StatusCreated, CreateCommandResponse{
		Uuid: commandUuid,
	})
}
