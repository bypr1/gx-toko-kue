package runner

import (
	"service/internal/app/database"
	"service/internal/pkg/config"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "xtreme:migration",
		Long: "Running Migration",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitTZ()

			DBClose := config.InitDB()
			defer DBClose()

			xtremedb.Migrate(config.PgSQL, database.Migrations(config.PgSQL))
		},
	})
}
