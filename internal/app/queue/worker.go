package queue

import (
	xtremequeue "github.com/globalxtreme/go-core/v2/queue"
	Telegram "service/internal/app/queue/job"
	"service/internal/pkg/constant"
)

func Register() []xtremequeue.JobConf {
	return []xtremequeue.JobConf{
		{
			Context:     Telegram.TelegramMessageJob{},
			JobFunc:     (*Telegram.TelegramMessageJob).Consume,
			Concurrency: 1,
			QueueName:   constant.QUEUE_HIGH,
			JobName:     constant.JOB_TELEGRAM_MESSAGE,
			Priority:    10,
		},
	}
}
