package storage

import (
	"context"
	"errors"
)

type mockStorage struct {
	commands []Command
}

func NewMockStorage() *mockStorage {
	testData := []Command{
		{
			Id:           0,
			Command_uuid: "aba5be85-f261-451f-a279-a3b0bbb6d1c8",
			Description:  "Test script 1",
			Script:       "echo Hello",
			Status:       "EXECUTED",
			Output:       "Hello",
		},
		{
			Id:           1,
			Command_uuid: "d3f135e8-7585-450f-bce3-41d40baec2ee",
			Description:  "Test script 2",
			Script:       "sleep -3",
			Status:       "FAILED",
			Output:       "",
		},
		{
			Id:           2,
			Command_uuid: "39b0fc6c-1a6a-43a7-85a3-6b6614190e77",
			Description:  "Test script 3",
			Script:       "sleep 3; echo Hi; sleep 5",
			Status:       "EXECUTING",
			Output:       "",
		},
		{
			Id:           3,
			Command_uuid: "c6235826-fe6f-47cf-93ee-74ea9d17910e",
			Description:  "Test script 4",
			Script:       "for i in $(seq 1 5); do echo Hyo; sleep 5; done",
			Status:       "STOPPED",
			Output:       "Hyo\nHyo\nHyo\n",
		},
	}

	return &mockStorage{
		commands: testData,
	}
}

func (ms *mockStorage) CreateCommand(ctx context.Context, command_uuid, description, script, status, output string) error {
	id := len(ms.commands)
	ms.commands = append(ms.commands, Command{
		Id:           id,
		Command_uuid: command_uuid,
		Description:  description,
		Script:       script,
		Status:       status,
		Output:       output,
	})

	return nil
}

func (ms *mockStorage) GetCommand(ctx context.Context, command_uuid string) (Command, error) {
	var command Command
	isFinded := false

	for i := 0; i < len(ms.commands); i++ {
		if ms.commands[i].Command_uuid == command_uuid {
			command = ms.commands[i]
			isFinded = true
			break
		}
	}

	if !isFinded {
		return command, errors.New("command not found")
	}

	return command, nil
}

func (ms *mockStorage) GetCommands(ctx context.Context, status string, limit int, offset int) ([]Command, error) {
	var commands []Command
	if status != "" {
		for _, cmd := range ms.commands {
			if cmd.Status == status {
				commands = append(commands, cmd)
			}
		}
	} else {
		commands = ms.commands
	}

	if limit > 0 && len(commands) > limit {
		commands = commands[:limit]
	}

	if offset > 0 {
		commands = commands[offset:]
	}

	return commands, nil
}

func (ms *mockStorage) UpdateCommandStatus(ctx context.Context, command_uuid string, status string) error {
	for i := 0; i < len(ms.commands); i++ {
		if ms.commands[i].Command_uuid == command_uuid {
			ms.commands[i].Status = status
			break
		}
	}

	return nil
}

func (ms *mockStorage) UpdateCommandOutput(ctx context.Context, command_uuid string, output string) error {
	for i := 0; i < len(ms.commands); i++ {
		if ms.commands[i].Command_uuid == command_uuid {
			ms.commands[i].Output = output
			break
		}
	}

	return nil
}
