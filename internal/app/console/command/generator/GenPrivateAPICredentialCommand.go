package generator

import (
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gorm.io/gorm/utils"
	"log"
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/encryption"
	"time"
)

type GenPrivateAPICredentialCommand struct {
	args    []string
	replace bool
}

func (c *GenPrivateAPICredentialCommand) Command(cobraCmd *cobra.Command) {
	addCommand := cobra.Command{
		Use:  "gen:private-api-credential",
		Long: "Command to generate a handler",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitPrivateAPICredential()

			c.args = args

			c.Handle()
		},
	}

	cobraCmd.AddCommand(&addCommand)
	addCommand.PersistentFlags().BoolVarP(&c.replace, "replace", "r", false, "Replace old credential")
}

func (c *GenPrivateAPICredentialCommand) Handle() {
	if len(c.args) < 1 {
		log.Printf("Please provide a credential name!")
		return
	}

	name := c.args[0]

	credential, ok := config.PrivateAPICredential[name]
	if ok {
		if c.replace {
			c.generateKey(credential)
		}
	} else {
		credential = map[string]interface{}{
			"id": utils.ToString(time.Now().UnixMicro()),
		}

		c.generateKey(credential)
	}

	credential["name"] = name

	c.generatePublicKey(credential)

	c.showCredential(credential)
}

func (c *GenPrivateAPICredentialCommand) showCredential(credential map[string]interface{}) {
	printCredential := func(key, value string) {
		fmt.Printf("%s %s\n", key, value)
	}

	fmt.Println("Please Save ID & Key to config (PrivateAPICredential) & .env file")
	fmt.Println("Please Give ID, Name, & Public Key to another service")
	printCredential(fmt.Sprintf("%-11s", "ID:"), credential["id"].(string))
	printCredential(fmt.Sprintf("%-11s", "Name:"), credential["name"].(string))
	printCredential(fmt.Sprintf("%-11s", "Key:"), credential["key"].(string))
	printCredential(fmt.Sprintf("%-3s", "Public Key:"), credential["secret"].(string))
}

func (c *GenPrivateAPICredentialCommand) generateKey(credential map[string]interface{}) {
	uuid7, _ := uuid.NewV7()

	credential["key"] = uuid7.String()
}

func (c *GenPrivateAPICredentialCommand) generatePublicKey(credential map[string]interface{}) {
	ec := encryption.NewPrivateAPIEncryption(credential["id"].(string))
	publicKey, err := ec.Encrypt(credential["key"].(string))
	if err != nil {
		log.Printf("Please provide a credential name!")
		os.Exit(1)
	}

	credential["secret"] = publicKey
}
