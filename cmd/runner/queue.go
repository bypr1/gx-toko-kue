package runner

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	xtremequeue "github.com/globalxtreme/go-core/v2/queue"
	"github.com/spf13/cobra"
	"service/internal/app/queue"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "xtreme:queue",
		Long: "Running Queue",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitTZ()

			logCleanup := xtremepkg.InitLogRPC()
			defer logCleanup()

			queueNames := cmd.Flags().String("q", constant.QUEUE_HIGH, "Queue name")
			configurations := queue.Register()

			worker := xtremequeue.Queue{Names: *queueNames}
			worker.Work(configurations)
		},
	})
}
