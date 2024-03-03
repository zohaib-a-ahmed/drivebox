package upload

import (
	"fmt"
	"log"
	"strconv"

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
			fmt.Println(parentID)
			UploadFileToDrive(filePath, driveService, parentID)
		case "2":
			parentID, err := CreateParentDirectory(driveService)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			UploadFileToDrive(filePath, driveService, parentID)
		case "3":
			fmt.Println("Exiting... Use command 'drivebox upload <path_to_file>' to upload under no directory.")
			return
		default:
			fmt.Println("Invalid choice. Exiting...")
		}
	},
}

func PromptUserForAction() string {
	fmt.Println("Select an action:")
	fmt.Println("1: Search for a parent directory")
	fmt.Println("2: Create a new parent directory")
	fmt.Println("3: Quit")
	fmt.Print("Selection: ")
	var choice string
	fmt.Scanln(&choice)
	return choice
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func searchFiles(svc *drive.Service, query string) ([]*drive.File, error) {
	searchQuery := fmt.Sprintf("name contains '%s' and mimeType = 'application/vnd.google-apps.folder'", query)
	call := svc.Files.List().Q(searchQuery).PageSize(5).Fields("files(id, name)")
	files, err := call.Do()
	if err != nil {
		return nil, err
	}
	return files.Files, nil
}

func SearchParentDirectory(svc *drive.Service) (string, error) {
	var files []*drive.File
	var err error

	for {
		input := getUserInput("Enter Drive parent directory name or 'quit' to exit: ")

		if input == "quit" {
			fmt.Println("Exiting...")
			return "", nil
		}

		if len(files) == 0 { // Initial search or no previous search results
			query := input
			files, err = searchFiles(svc, query)
			if err != nil {
				log.Printf("Failed to retrieve files: %v\n", err)
				continue
			}
			if len(files) == 0 {
				fmt.Println("No files found. Try refining your search query.")
				continue
			}
		}

		for i, file := range files {
			fmt.Printf("%d: %s (ID: %s)\n", i+1, file.Name, file.Id)
		}

		// Selection loop
		for {
			input = getUserInput("Selection: ")

			if input == "quit" {
				fmt.Println("Exiting...")
				return "", nil
			} else {
				selection, err := strconv.Atoi(input)
				if err != nil || selection < 1 || selection > len(files) {
					fmt.Println("Invalid selection - try again...")
					continue
				}
				return files[selection-1].Id, nil
			}
		}
	}
}

func CreateParentDirectory(svc *drive.Service) (string, error) {

	dirName := getUserInput("Name the new directory: ")

	// Search for an existing directory with the same name
	query := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and trashed=false", dirName)
	call := svc.Files.List().Q(query).Fields("files(id, name)")
	files, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("unable to search for directories: %v", err)
	}

	if len(files.Files) > 0 {
		return "", fmt.Errorf("a directory with the name '%s' already exists", dirName)
	}

	dirMetadata := &drive.File{
		Name:     dirName,
		MimeType: "application/vnd.google-apps.folder",
	}
	newDir, err := svc.Files.Create(dirMetadata).Fields("id").Do()
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	fmt.Println("Directory created successfully.")
	return newDir.Id, nil
}
