package main

// ImmichUp CLI tool for uploading media files to Immich server
// build:
// GOOS=linux GOARCH=amd64 go build -o iu-linux iu.go

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBaseURL = "http://127.0.0.1:2283/api" // replace as needed
)

type Response struct {
	ID        string `json:"id"`
	Duplicate bool   `json:"duplicate"`
}

func upload(filePath string, apiKey string, baseURL string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	modTime := fileInfo.ModTime()

	// Create a buffer to write our multipart form to
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file
	fileField, err := writer.CreateFormFile("assetData", filepath.Base(filePath))
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(fileField, file)
	if err != nil {
		return err
	}

	// Add other form fields
	_ = writer.WriteField("deviceAssetId", fmt.Sprintf("%s-%d", filePath, modTime.Unix()))
	_ = writer.WriteField("deviceId", "golang")
	_ = writer.WriteField("fileCreatedAt", modTime.Format(time.RFC3339))
	_ = writer.WriteField("fileModifiedAt", modTime.Format(time.RFC3339))
	_ = writer.WriteField("isFavorite", "false")

	// Close the writer
	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", baseURL+"/assets", &requestBody)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", response)
	return nil
}

// Common image and video file extensions
var mediaExtensions = map[string]bool{
	// Images
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".tiff": true,
	".webp": true,
	// Videos
	".mp4":  true,
	".mov":  true,
	".avi":  true,
	".mkv":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
}

// findMediaFiles scans a directory recursively for media files
func findMediaFiles(directory string) ([]string, error) {
	var mediaFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file extension is in our map
		ext := strings.ToLower(filepath.Ext(path))
		if mediaExtensions[ext] {
			mediaFiles = append(mediaFiles, path)
		}

		return nil
	})

	return mediaFiles, err
}

func main() {
	var serverURL string
	flag.StringVar(&serverURL, "url", defaultBaseURL, "Immich server API URL")
	flag.Parse()

	// Check if directory is provided as command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run iu.go /path/to/your/media/directory")
		fmt.Println("Make sure to set the IMAGEUP_API_KEY environment variable")
		os.Exit(1)
	}

	// Check if API key is set
	if os.Getenv("IMAGEUP_API_KEY") == "" {
		fmt.Println("Error: IMAGEUP_API_KEY environment variable not set")
		fmt.Println("Set it with: export IMAGEUP_API_KEY=your_api_key")
		os.Exit(1)
	}

	// Get directory from command line argument
	targetDirectory := os.Args[1]

	// Find all media files
	mediaFiles, err := findMediaFiles(targetDirectory)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d media files\n", len(mediaFiles))

	// Upload each file
	for _, filePath := range mediaFiles {
		fmt.Printf("Uploading %s...\n", filePath)
		err := upload(filePath, os.Getenv("IMAGEUP_API_KEY"), serverURL)
		if err != nil {
			fmt.Printf("Error uploading file %s: %v\n", filePath, err)
		}
	}
}
