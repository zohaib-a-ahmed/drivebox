package unload

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var UnloadCmd = &cobra.Command{
	Use:   "unload <file_name> <optional_path_destination>",
	Short: "Download a file from Google Drive",
	Long: `Download a file from your Google Drive to a local path. 
If no destination is provided, the file will be downloaded to the current directory.`,
	Args: cobra.MinimumNArgs(1), // Ensures at least one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		destination := "./" // Default to current directory if no destination is provided
		if len(args) > 1 {
			destination = args[1]
			// Validate or adjust the destination path if necessary
			if _, err := os.Stat(destination); os.IsNotExist(err) {
				log.Fatalf("The specified destination does not exist: %s", destination)
			}
		}

		fmt.Printf("Downloading '%s' to '%s'...\n", fileName, destination)
		// Placeholder for search, selection, and download logic
	},
}
