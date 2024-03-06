package auth

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var SetUpCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up access to Google Drive",
	Long:  `Set up access to Google Drive by providing client credentials under your own project.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load(".env")
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("Error loading .env file: %v", err)
		}

		clientID := os.Getenv("GOOGLE_CLIENT_ID")
		clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

		if clientID != "" && clientSecret != "" {
			log.Println("Credentials already exist.")
			if !userConsent("Do you want to change these credentials? (yes/no): ") {
				log.Println("Exiting setup.")
				return
			}
		}

		clientID = getUserInput("Enter GOOGLE_CLIENT_ID: ")
		clientSecret = getUserInput("Enter GOOGLE_CLIENT_SECRET: ")

		updateEnvFile(clientID, clientSecret)
		log.Println("Setup complete. Check credential validity with 'drivebox auth in'.")
	},
}

func userConsent(prompt string) bool {
	answer := getUserInput(prompt)
	return strings.ToLower(answer) == "yes"
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func updateEnvFile(clientID, clientSecret string) {
	file, err := os.OpenFile(".env", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open .env file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("GOOGLE_CLIENT_ID=%s\nGOOGLE_CLIENT_SECRET=%s\n", clientID, clientSecret))
	if err != nil {
		log.Fatalf("Failed to write to .env file: %v", err)
	}
}
