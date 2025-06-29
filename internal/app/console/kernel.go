package console

import (
	"service/internal/app/console/command"
	"service/internal/app/console/command/generator"
	"service/internal/app/console/command/rabbitmq"

	xtremeconsole "github.com/globalxtreme/go-core/v2/console"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

func RegisterCommand(cobraCmd *cobra.Command) {
	xtremeconsole.Commands(cobraCmd, []xtremeconsole.BaseCommand{
		// File Generator
		&generator.GenMigrationCommand{},
		&generator.GenParserCommand{},
		&generator.GenHandlerCommand{},
		&generator.GenModelCommand{},
		&generator.GenPrivateAPICredentialCommand{},

		// RabbitMQ Consumer
		&rabbitmq.RabbitMQConsumerGlobalCommand{},
		&rabbitmq.RabbitMQConsumerLocalCommand{},

		&command.TestCommand{},
	})
}

func RegisterSchedule(sch *gocron.Scheduler) {
	//addSchedule(sch.Every(1).Minute(), &Command.TestCommand{})
}

func addSchedule(schedule *gocron.Scheduler, command xtremeconsole.BaseCommand) {
	schedule.Do(command.Handle)
}
