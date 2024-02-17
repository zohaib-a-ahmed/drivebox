package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// inCmd represents the 'auth in' command
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check Google Drive authentication",
	Long:  `Check whether current session is authenticated with Google OAuth.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking authentication...")
		// Implementation of the sign-in logic
	},
}
