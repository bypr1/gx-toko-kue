package thirdparty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/gocraft/work"
	"net/http"
	"service/internal/pkg/constant"
)

type Telegram struct {
	BootToken  string
	RoomChatId float64
}

func (tel *Telegram) SetFromDeveloper() {
	tel.BootToken = constant.TELEGRAM_BOOT_TEKON_DEV
	tel.RoomChatId = constant.TELEGRAM_ROOM_CHAT_DEV
}

func (tel *Telegram) Queue(message string) {
	if len(tel.BootToken) == 0 {
		tel.SetFromDeveloper()
	}

	args := map[string]interface{}{
		"booToken":   tel.BootToken,
		"roomChatId": tel.RoomChatId,
		"message":    message,
	}

	ctx := work.NewEnqueuer(constant.QUEUE_HIGH, xtremepkg.RedisPool)
	_, err := ctx.Enqueue(constant.JOB_TELEGRAM_MESSAGE, args)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
}

func (tel *Telegram) Send(message string) error {
	if len(tel.BootToken) == 0 {
		tel.SetFromDeveloper()
	}

	content, err := json.Marshal(map[string]interface{}{"chat_id": tel.RoomChatId, "text": message, "parse_mode": "Markdown"})
	if err != nil {
		fmt.Println(fmt.Sprintf("TELEGRAM: Unable to marshal. Err: %v", err))
		return err
	}

	res, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tel.BootToken),
		"application/json",
		bytes.NewBuffer(content),
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("TELEGRAM: Unable to send message. Err: %v", err))
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Unexpected status %v", res.Status))
	}

	return nil
}
