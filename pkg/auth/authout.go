package auth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var OutCmd = &cobra.Command{
	Use:   "out",
	Short: "Sign out from Google Drive",
	Long:  `Sign out from Google Drive and clear the authentication token.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Signing out from Google Drive...")
		err := os.Remove("token.json")
		if err != nil {
			if os.IsNotExist(err) {
				// The file does not exist, which means the user is not currently authenticated
				fmt.Println("No current session is authenticated. Run 'drivebox auth in' to authenticate.")
			} else {
				// An error other than the file not existing occurred
				fmt.Printf("Error removing token file: %v\n", err)
			}
		} else {
			fmt.Println("Successfully signed out.")
		}
	},
}
