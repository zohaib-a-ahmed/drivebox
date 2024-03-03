package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var InCmd = &cobra.Command{
	Use:   "in",
	Short: "Sign in to Google Drive",
	Long:  `Sign in to Google Drive using OAuth2.0 authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		signIn()
		log.Println("Authentication successful!")
	},
}

func signIn() {
	server := &http.Server{Addr: ":9999"}

	// Generate the OAuth URL and open it in the browser.
	url := NewConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	//fmt.Printf("Opening browser to visit: %s\n", url)

	// Open URL in the user's browser.
	if err := browser.OpenURL(url); err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}

	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Extract the code from the query string.
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found in the request", http.StatusBadRequest)
			return
		}

		// Exchange the authorization code for an access token.
		token, err := NewConfig().Exchange(ctx, code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := saveToken("token.json", token); err != nil {
			http.Error(w, "Failed to save token", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Authentication successful! You may now close this window.")
		time.Sleep(3 * time.Second)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	})

	// Start the HTTP server.
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

func loadTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func CreateDriveService() (*drive.Service, error) {
	// Create a new OAuth2 HTTP client using the token
	ctx := context.Background()
	config := NewConfig()

	token, err := loadTokenFromFile("token.json")
	if err != nil {
		log.Fatalf("Unable to load token; Check authenticatio with 'drivebox auth check' \n %v", err)
	}

	client := config.Client(ctx, token)

	// Create a new Google Drive service client
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("cannot create drive service: %v", err)
	}
	return srv, nil
}
