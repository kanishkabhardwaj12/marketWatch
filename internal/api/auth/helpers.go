package auth

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func savePhoto(file multipart.File, filename, storagePath string) (string, error) {
	defer file.Close()
	uniqueFilename := fmt.Sprintf("%s_%s", uuid.NewString(), filename)

	os.MkdirAll(storagePath, os.ModePerm)
	fullPath := filepath.Join(storagePath, uniqueFilename)

	dest, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		return "", err
	}

	return "/photo_uploads/" + uniqueFilename, nil
}

func saveFile(file multipart.File, filename, storagePath string) (string, error) {
	defer file.Close()
	uniqueFilename := fmt.Sprintf("%s_%s", uuid.NewString(), filename)

	os.MkdirAll(storagePath, os.ModePerm)
	fullPath := filepath.Join(storagePath, uniqueFilename)

	dest, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		return "", err
	}

	return "/tradebook_uploads/" + uniqueFilename, nil
}
