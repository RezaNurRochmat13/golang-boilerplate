package email

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TypeSendEmail = "email:send"

type Payload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewSendEmailTask(payload Payload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendEmail, data), nil
}
