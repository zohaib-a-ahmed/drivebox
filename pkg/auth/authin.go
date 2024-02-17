package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// inCmd represents the 'auth in' command
var InCmd = &cobra.Command{
	Use:   "in",
	Short: "Sign in to Google Drive",
	Long:  `Sign in to Google Drive using OAuth authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Signing in to Google Drive...")
		// Implementation of the sign-in logic
	},
}
