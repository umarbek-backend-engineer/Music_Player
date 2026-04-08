package utils

import (
	"fmt"
	"mime/multipart"
)

func FileValidator(fileHeader *multipart.FileHeader) error {
	const maxFileSize = 10 << 20

	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("File too big")
	}

	contentType := fileHeader.Header.Get("Content-Type")

	if contentType != "audio/mpeg" && contentType != "audio/wav" {
		return fmt.Errorf("Only MP3 or WAV allowed")
	}
	return nil
}
