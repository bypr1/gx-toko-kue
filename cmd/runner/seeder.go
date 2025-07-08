package runner

import (
	"service/internal/app/database"
	"service/internal/pkg/config"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "xtreme:seeder",
		Long: "Running Seeder",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitTZ()

			DBClose := config.InitDB()
			defer DBClose()

			database.Seeder(config.PgSQL)
		},
	})
}
