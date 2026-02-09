package email

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

func HandleSendEmail(ctx context.Context, t *asynq.Task) error {
	var payload Payload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// === SIMULASI KIRIM EMAIL ===
	log.Println("[EMAIL WORKER]")
	log.Printf("To      : %s", payload.To)
	log.Printf("Subject : %s", payload.Subject)
	log.Printf("Body    : %s", payload.Body)

	// jika gagal → return error → Asynq auto retry
	return nil
}
