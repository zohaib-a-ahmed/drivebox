package upload

import (
	"fmt"

	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to Google Drive",
	Long:  `Upload a file to your Google Drive. Specify the local path and optional Google Drive path.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uploading file to Google Drive...")
	},
}
