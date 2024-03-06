package auth

import (
	"log"

	"github.com/spf13/cobra"
)

var SetUpCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up access to Google Drive",
	Long:  `Set up access to Google Drive by providing client credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Setting up!")
	},
}
