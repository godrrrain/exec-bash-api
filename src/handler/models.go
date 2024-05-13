package handler

import (
	"github.com/godrrrain/exec-bash-api/src/storage"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type CreateCommandRequest struct {
	Description string `json:"description"`
	Script      string `json:"script" binding:"required"`
}

type CreateCommandResponse struct {
	Uuid string `json:"uuid"`
}

type CommandResponse struct {
	Command_uuid string `json:"command_uuid"`
	Description  string `json:"description"`
	Script       string `json:"script"`
	Status       string `json:"status"`
	Output       string `json:"output"`
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func NewErrorResponse(err string) ErrorResponse {
	return ErrorResponse{
		Error: err,
	}
}

func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{
		Message: message,
	}
}

func CommandToResponse(command storage.Command) CommandResponse {
	return CommandResponse{
		Command_uuid: command.Command_uuid,
		Description:  command.Description,
		Script:       command.Script,
		Status:       command.Status,
		Output:       command.Output,
	}
}

func CommandsToResponse(commands []storage.Command) []CommandResponse {
	if commands == nil {
		return nil
	}

	res := make([]CommandResponse, len(commands))

	for index, value := range commands {
		res[index] = CommandToResponse(value)
	}

	return res
}
