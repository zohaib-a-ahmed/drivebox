package upload

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zohaib-a-ahmed/drivebox/pkg/auth"
	"google.golang.org/api/drive/v3"
)

var UploadParentCmd = &cobra.Command{
	Use:   "parent <path_to_file>",
	Short: "Upload a file to Google Drive under a parent directory",
	Long:  `Upload a file to your Google Drive. Specify the local path after selecting an existing parent directory`,
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

		choice := PromptUserForAction()
		switch choice {
		case "1":
			parentID, err := SearchParentDirectory(driveService)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			UploadFileToDrive(filePath, driveService, parentID)
		case "2":
			ForegoParentDirectory(filePath, driveService)
		case "3":
			parentID, err := CreateParentDirectory(driveService)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			UploadFileToDrive(filePath, driveService, parentID)
		case "4":
			fmt.Println("Quitting...")
			return
		default:
			fmt.Println("Invalid choice. Quitting...")
		}
	},
}

func PromptUserForAction() string {
	fmt.Println("Select an action:")
	fmt.Println("1: Search for a parent directory")
	fmt.Println("2: Forego specifying a parent directory")
	fmt.Println("3: Create a new parent directory")
	fmt.Println("4: Quit")
	fmt.Print("Selection: ")
	var choice string
	fmt.Scanln(&choice)
	return choice
}

func SearchParentDirectory(driveService *drive.Service) (string, error) {
	// Implement search functionality here
	fmt.Println("searching rents")
	return "searching rents", nil
}

func ForegoParentDirectory(filePath string, driveService *drive.Service) {
	// Implement logic to upload to Google Drive's root or a default location
	fmt.Println("foregoing rents")
}

func CreateParentDirectory(driveService *drive.Service) (string, error) {
	// Implement directory creation logic here
	fmt.Println("creating rents")
	return "creating rent", nil
}
