package utils

import (
	"net/http"
	"os"
	"path/filepath"
)

func SaveFile(uploadDir, filename string, content []byte) (string, error) {

	// checking the size of the file type
	size := len(content)
	if size > 10*1024*1024 {
		return "", ErrFileTooLarg
	}

	// checking the type fo the file
	fileType := http.DetectContentType(content)
	if fileType != "audio/mpeg" && fileType != "audio/wav" {
		return "", ErrInvalidFile
	}

	// checking if the files for storage exists. if not create one
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	//build path
	fullpath := filepath.Join(uploadDir, filename)

	// write file
	if err := os.WriteFile(fullpath, content, 0644); err != nil {
		return "", nil
	}

	return fullpath, nil
}
