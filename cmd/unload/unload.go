package unload

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zohaib-a-ahmed/drivebox/pkg/auth"
	"google.golang.org/api/drive/v3"
)

var UnloadCmd = &cobra.Command{
	Use:   "unload <file_name> <optional_path_destination>",
	Short: "Download a file from Google Drive",
	Long: `Download a file from your Google Drive to a local path.
If no destination is provided, the file will be downloaded to the current directory.
Enter the number of the file you wish to download, use 'refine <query>' to narrow down your search, or type 'quit' to exit the command.`,
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

		handleUserSelection(driveService, files.Files, destination)
	},
}

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func searchFiles(driveService *drive.Service, query string) ([]*drive.File, error) {
	call := driveService.Files.List().Q(query).PageSize(6).Fields("files(id, name)")
	files, err := call.Do()
	if err != nil {
		return nil, err
	}
	return files.Files, nil
}

func handleUserSelection(driveService *drive.Service, files []*drive.File, destination string) {
	for {
		input := getUserInput("Enter the number of the file to download, 'refine <query>' to search again, or 'quit' to exit: ")

		if input == "quit" {
			fmt.Println("Exiting command.")
			return // Exit the function, effectively ending the command
		}

		if strings.HasPrefix(input, "refine ") {
			query := strings.TrimSpace(strings.TrimPrefix(input, "refine"))
			fmt.Println("Refining search with:", query)
			var err error
			files, err = searchFiles(driveService, fmt.Sprintf("name contains '%s'", query))
			if err != nil {
				log.Printf("Failed to retrieve files: %v\n", err)
				continue
			}
			if len(files) == 0 {
				fmt.Println("No files found. Try refining your search.")
				continue
			}
			for i, file := range files {
				fmt.Printf("%d: %s\n", i+1, file.Name)
			}
		} else {
			selection, err := strconv.Atoi(input)
			if err != nil || selection < 1 || selection > len(files) {
				fmt.Println("Invalid selection. Please enter a valid number or 'quit' to exit.")
				continue
			}
			downloadFile(driveService, files[selection-1].Id, destination)
			break
		}
	}
}

func downloadFile(driveService *drive.Service, fileId, destinationPath string) error {
	file, err := driveService.Files.Get(fileId).Fields("name", "mimeType").Do()
	if err != nil {
		return fmt.Errorf("failed to get file: %v", err)
	}

	var resp *http.Response
	var fileName string // Variable to hold the final file name including extension

	if strings.Contains(file.MimeType, "google-apps") {
		// Determine the correct export MIME type and corresponding file extension
		exportMimeType, fileExtension := determineExportFormat(file.MimeType)
		resp, err = driveService.Files.Export(fileId, exportMimeType).Download()
		if err != nil {
			return fmt.Errorf("failed to export and download file: %v", err)
		}
		// Append the appropriate file extension to the original file name
		fileName = sanitizeFileName(file.Name) + fileExtension
	} else {
		// For binary files, directly download and use the original file name
		resp, err = driveService.Files.Get(fileId).Download()
		if err != nil {
			return fmt.Errorf("failed to download file: %v", err)
		}
		fileName = sanitizeFileName(file.Name)
	}
	defer resp.Body.Close()

	// Ensure the destination path ends with a separator
	if !strings.HasSuffix(destinationPath, "/") && !strings.HasSuffix(destinationPath, "\\") {
		destinationPath += "/"
	}

	// Create a file at the specified destination with the correct file extension
	outFile, err := os.Create(destinationPath + fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	// Write the content to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("Download complete: %s\n", destinationPath+fileName)
	return nil
}

// determineExportFormat returns the MIME type and file extension for exporting Google Docs formats.
func determineExportFormat(mimeType string) (exportMimeType, fileExtension string) {
	switch mimeType {
	case "application/vnd.google-apps.document":
		return "application/pdf", ".pdf"
	case "application/vnd.google-apps.spreadsheet":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ".xlsx"
	case "application/vnd.google-apps.presentation":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation", ".pptx"
	// Add more cases as needed
	default:
		return "application/pdf", ".pdf" //
	}
}

// sanitizeFileName cleans up the file name to prevent path traversal vulnerabilities or issues with illegal characters.
func sanitizeFileName(name string) string {
	// Simple sanitization logic; consider expanding based on your needs.
	return strings.ReplaceAll(name, "/", "_")
}
