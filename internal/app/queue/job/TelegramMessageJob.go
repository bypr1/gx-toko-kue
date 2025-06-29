package Telegram

import (
	"fmt"
	"github.com/gocraft/work"
	"github.com/mitchellh/mapstructure"
	"log"
	"service/internal/pkg/core"
	"service/internal/pkg/thirdparty"
)

type TelegramMessageJob struct {
	BootToken  string  `json:"bootToken"`
	RoomChatId float64 `json:"roomChatId"`
	Message    string  `json:"message"`
}

func (j *TelegramMessageJob) Consume(job *work.Job) error {
	err := core.ErrorHandler(func() error {
		log.Println(fmt.Sprintf("SendTelegramMessageJob::PROCESSING"))

		mapstructure.Decode(job.Args, &j)

		tel := thirdparty.Telegram{
			BootToken:  j.BootToken,
			RoomChatId: j.RoomChatId,
		}

		err := tel.Send(j.Message)
		if err != nil {
			return err
		}

		log.Println(fmt.Sprintf("SendTelegramMessageJob::SUCCESS"))

		return nil
	})

	return err
}
