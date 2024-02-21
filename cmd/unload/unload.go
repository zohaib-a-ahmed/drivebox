package unload

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zohaib-a-ahmed/drivebox/pkg/auth"
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

		driveService, err := auth.CreateDriveService()
		if err != nil {
			log.Fatalf("Failed to create Google Drive service: %v", err)
		}

		// Search for the file on Google Drive
		query := fmt.Sprintf("name contains '%s'", fileName)
		call := driveService.Files.List().Q(query).PageSize(6).Fields("files(id, name)")
		files, err := call.Do()
		if err != nil {
			log.Fatalf("Failed to retrieve files: %v", err)
		}
		if len(files.Files) == 0 {
			fmt.Println("No files found.")
			return
		}
		fmt.Println("Files found:")
		for i, file := range files.Files {
			fmt.Printf("%d: %s \n", i+1, file.Name)
		}
	},
}
