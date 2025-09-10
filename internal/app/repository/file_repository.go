package repository

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// FileRepository defines the contract for a repository that manages file storage.
// This interface can be implemented by different storage types (local, S3, etc.).
type FileRepository interface {
	Save(fileHeader *multipart.FileHeader, folderPath string) (string, error)
	Delete(filePath string) error
}

// fileSystemRepo is the concrete implementation for local file storage.
type fileSystemRepo struct {
	basePath string
}

// NewFileSystemRepo creates a new instance of FileSystemRepository.
func NewFileSystemRepo(basePath string) FileRepository {
	return &fileSystemRepo{basePath: basePath}
}

// Save saves a file to the local file system and returns its path.
func (r *fileSystemRepo) Save(fileHeader *multipart.FileHeader, folderPath string) (string, error) {

	save_file_path := filepath.Join(r.basePath, folderPath)

	// 1. Ensure the upload directory exists.
	if err := os.MkdirAll(save_file_path, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Create a unique filename.
	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	internalFilePath := filepath.Join(save_file_path, filename)
	returnFilePath := filepath.Join(folderPath, filename)

	// 3. Open the destination file for writing.
	dst, err := os.Create(internalFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// 4. Open the source file from the multipart form data.
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file header: %w", err)
	}
	defer src.Close()

	// 5. Use io.Copy for efficient and safe file copying.
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	if err := dst.Sync(); err != nil {
		return "", fmt.Errorf("failed to sync file: %w", err)
	}

	return returnFilePath, nil
}

// Delete removes a file from the local file system.
func (r *fileSystemRepo) Delete(filePath string) error {
	internalFilePath := filepath.Join(r.basePath, filePath)
	return os.Remove(internalFilePath)
}
