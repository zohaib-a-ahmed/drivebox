package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	AuthCmd.AddCommand(InCmd)
	AuthCmd.AddCommand(OutCmd)
	AuthCmd.AddCommand(CheckCmd)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// NewConfig creates and returns a new oauth2.Config instance using environment variables.
func NewConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:9999/oauth/callback",
	}
}

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Google Drive",
	Long:  `Begin OAuth authentication or refresh the token with Google Drive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("subcommand required: in, out, or check")
	},
}
