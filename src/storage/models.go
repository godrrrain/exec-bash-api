package storage

type Command struct {
	Id           int    `json:"id"`
	Command_uuid string `json:"command_uuid"`
	Description  string `json:"description"`
	Script       string `json:"script"`
	Status       string `json:"status"`
	Output       string `json:"output"`
}
