package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zohaib-a-ahmed/drivebox/cmd/unload"
	"github.com/zohaib-a-ahmed/drivebox/cmd/upload"
	"github.com/zohaib-a-ahmed/drivebox/pkg/auth"
)

var rootCmd = &cobra.Command{
	Use:   "drivebox",
	Short: "Drivebox is a CLI tool for managing Google Drive files",
	Long:  `Drivebox allows you to easily upload, download, and manage your Google Drive files from the command line.`,
}

func main() {
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(upload.UploadCmd)
	rootCmd.AddCommand(unload.UnloadCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
