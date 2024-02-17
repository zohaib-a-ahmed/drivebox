package unload

import (
	"fmt"

	"github.com/spf13/cobra"
)

// unloadCmd represents the unload command
var UnloadCmd = &cobra.Command{
	Use:   "unload",
	Short: "Download a file from Google Drive",
	Long:  `Download a file from your Google Drive to a local path.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Downloading file from Google Drive...")
	},
}
