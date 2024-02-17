package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// outCmd represents the 'auth out' command
var OutCmd = &cobra.Command{
	Use:   "out",
	Short: "Sign out from Google Drive",
	Long:  `Sign out from Google Drive and clear the authentication token.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Signing out from Google Drive...")
		// Implementation of the sign-out logic
	},
}
