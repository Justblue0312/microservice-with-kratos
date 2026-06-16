package asynq

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TypeGoodbyeSaid = "goodbye:said"

type GoodbyeSaidPayload struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewGoodbyeSaidTask(name, message, timestamp string) (*asynq.Task, error) {
	payload := GoodbyeSaidPayload{Name: name, Message: message, Timestamp: timestamp}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeGoodbyeSaid, data), nil
}
