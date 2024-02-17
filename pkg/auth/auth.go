package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	AuthCmd.AddCommand(InCmd)
	AuthCmd.AddCommand(OutCmd)
	AuthCmd.AddCommand(CheckCmd) // Assuming you have a CheckCmd defined somewhere
}

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Google Drive",
	Long:  `Begin OAuth authentication or refresh the token with Google Drive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("subcommand required: in, out, or check")
	},
}
