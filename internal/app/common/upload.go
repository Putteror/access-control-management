package common

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// UploadFile saves a multipart file to the specified directory.
// It generates a unique filename to prevent conflicts.
func UploadFile(fileHeader *multipart.FileHeader, destinationDir string) (string, error) {
	// 1. Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 2. Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// 3. Generate a unique filename
	fileExtension := filepath.Ext(fileHeader.Filename)
	fileName := uuid.New().String() + fileExtension
	filePath := filepath.Join(destinationDir, fileName)

	// 4. Create a new file on the server to save the uploaded content
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// 5. Copy the uploaded file content to the new file
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	// 6. Return the path of the saved file
	return filePath, nil
}

// GetImageURL constructs the full URL for a saved image.
// This function is for demonstration and needs to be adapted to your web server setup.
func GetImageURL(baseURI, filePath string) string {
	if filePath == "" {
		return ""
	}
	// Replace backslashes with forward slashes for URL compatibility
	safePath := strings.ReplaceAll(filePath, "\\", "/")
	return fmt.Sprintf("%s/%s", baseURI, safePath)
}
