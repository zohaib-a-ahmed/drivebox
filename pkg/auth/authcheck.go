package auth

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check Google Drive authentication",
	Long:  `Check whether the current session is authenticated with Google OAuth.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Checking authentication...")

		token, err := loadToken("token.json")
		if err != nil {
			log.Println("Failed to load token. Use 'drivebox auth in' to authenticate.")
			return
		}

		client := NewConfig().Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/drive/v3/files/root?fields=id")
		if err != nil {
			log.Fatalf("Failed to make outgoing API request; config (client) credentials may not have been set up:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			log.Println("Current session authorized!")
		} else {
			log.Println("Current session is not authorized. Use 'drivebox auth in' to authenticate.")
		}
	},
}

func loadToken(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}
