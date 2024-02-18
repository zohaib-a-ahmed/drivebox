package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var InCmd = &cobra.Command{
	Use:   "in",
	Short: "Sign in to Google Drive",
	Long:  `Sign in to Google Drive using OAuth2.0 authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting authentication process...")
		signIn()
	},
}

func signIn() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server := &http.Server{Addr: ":9999"}

	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:9999/oauth/callback",
	}

	// Generate the OAuth URL and open it in the browser.
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	//fmt.Printf("Opening browser to visit: %s\n", url)

	// Open URL in the user's browser.
	if err := browser.OpenURL(url); err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}

	// This function is defined within your signIn() function or globally if you prefer.
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Extract the code from the query string.
		code := r.URL.Query().Get("code")
		if code == "" {
			// Handle the error case.
			http.Error(w, "Code not found in the request", http.StatusBadRequest)
			log.Println("Code not found in the request")
			return
		}

		// Exchange the authorization code for an access token.
		token, err := config.Exchange(ctx, code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to exchange token: %v\n", err)
			return
		}

		if err := saveToken("token.json", token); err != nil {
			http.Error(w, "Failed to save token", http.StatusInternalServerError)
			log.Printf("Failed to save token: %v\n", err)
			return
		}

		fmt.Fprintf(w, "Authentication successful! You may now close this window.")
		log.Println("Authentication successful! Shutting down server..")
		time.Sleep(3 * time.Second)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	})

	// Start the HTTP server.
	fmt.Println("Starting local server to receive the authorization code...")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func saveToken(path string, token *oauth2.Token) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file %s for token: %v", path, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(token)
	if err != nil {
		return fmt.Errorf("failed to encode token to file %s: %v", path, err)
	}

	return nil
}
