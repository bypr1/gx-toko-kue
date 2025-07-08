package main

import (
	"fmt"
	"net/http"
	"os"
	"service/internal/app/api"
	"service/internal/app/database"
	"service/internal/pkg/config"

	xtremecore "github.com/globalxtreme/go-core/v2"
	xtremedb "github.com/globalxtreme/go-core/v2/database"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "root",
	Long: "Running service api",
	Run: func(cmd *cobra.Command, args []string) {
		xtremepkg.InitDevMode()
		xtremepkg.InitHost()

		config.InitTZ()
		config.InitCors()
		config.InitRPC()
		config.InitValidation()

		DBClose := config.InitCakeDB()
		defer DBClose()

		xtremedb.Migrate(config.CakeSQL, database.Migrations(config.CakeSQL))

		newCors := cors.New(config.CorsOptions)

		newRoute := mux.NewRouter()
		xtremecore.RegisterRouter(newRoute, api.Register)

		fmt.Println(fmt.Sprintf("Server started on %s", xtremepkg.HostFull))

		err := http.ListenAndServe(xtremepkg.Host, newCors.Handler(newRoute))
		if err != nil {
			panic(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// TODO: awalnya dev mode = false
	rootCmd.PersistentFlags().BoolVar(&xtremepkg.DevMode, "dev", false, "Set for development mode")
}

func main() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "xtreme:seeder",
		Long: "Running Seeder",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitTZ()

			DBClose := config.InitCakeDB()
			defer DBClose()

			database.Seeder(config.CakeSQL)
		},
	})
	Execute()
}
