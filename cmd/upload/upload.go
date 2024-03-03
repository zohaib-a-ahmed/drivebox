package upload

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zohaib-a-ahmed/drivebox/pkg/auth"
	"google.golang.org/api/drive/v3"
)

func init() {
	UploadCmd.AddCommand(UploadParentCmd)
}

var UploadCmd = &cobra.Command{
	Use:   "upload <path_to_file>",
	Short: "Upload a file to Google Drive",
	Long:  `Upload a file to your Google Drive. Specify the local path.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Ensure valid command skeleton
		if len(args) < 1 {
			fmt.Println("Path to the local file must be provided. Usage is 'drivebox upload <path_to_file>'")
			return
		}

		var filePath = args[0]
		// Ensure unload destination is valid
		if CheckValidPath(filePath) != 200 {
			return
		}

		// Create drive service
		driveService, err := auth.CreateDriveService()
		if err != nil {
			log.Fatalf("Failed to create Google Drive service: %v", err)
		}

		UploadFileToDrive(filePath, driveService, "")
	},
}

func CheckValidPath(path string) (code int) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("The file does not exist at the specified path: %s\n", path)
		return 400
	}
	return 200
}

func UploadFileToDrive(filePath string, svc *drive.Service, parentID string) (int, error) {

	// Gather file information
	file, err := os.Open(filePath)
	if err != nil {
		return 400, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return 400, err
	}

	// Initialize parents slice based on parentID
	var parents []string
	if parentID != "" {
		parents = []string{parentID}
	}

	// Create File metadata
	f := &drive.File{
		Name:    fileInfo.Name(),
		Parents: parents,
	}

	// Create and upload the file
	res, err := svc.Files.Create(f).
		Media(file).
		ProgressUpdater(func(now, size int64) { fmt.Printf("%d, %d\r", now, size) }).
		Do()

	if err != nil {
		return 400, err
	}
	if res.Id == "" {
		return 200, nil
	}

	fmt.Println("Successful Upload.")
	return 200, nil
}
