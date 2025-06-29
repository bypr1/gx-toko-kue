package runner

import (
	"fmt"
	xtremecore "github.com/globalxtreme/go-core/v2"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"net/http"
	"service/internal/app/privateapi"
	"service/internal/pkg/config"
	"time"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "xtreme:private-api",
		Long: "Running Private API",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitPrivateAPI()
			config.InitTZ()
			config.InitCors()
			config.InitMail()
			config.InitValidation()
			config.InitCache(time.Hour, time.Hour)

			DBClose := config.InitDB()
			defer DBClose()

			rabbitMQClose := config.InitRabbitMQ()
			defer rabbitMQClose()

			logCleanup := xtremepkg.InitLogRPC()
			defer logCleanup()

			newRoute := mux.NewRouter()
			xtremecore.RegisterRouter(newRoute, privateapi.Register)

			fmt.Println(fmt.Sprintf("Server started on %s", config.PrivateHost))

			err := http.ListenAndServe(config.PrivateHost, newRoute)
			if err != nil {
				panic(err)
			}
		},
	})
}
