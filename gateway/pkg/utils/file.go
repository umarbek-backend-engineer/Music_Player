package utils

import (
	"fmt"
	"io"
	"mime/multipart"
)

func FileValidator(fileHeader *multipart.FileHeader, file multipart.File) error {
	const maxFileSize = 10 << 20

	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("File too big")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "" && contentType != "audio/mpeg" && contentType != "audio/wav" && contentType != "audio/x-wav" && contentType != "audio/wave" {
		return fmt.Errorf("Only MP3 or WAV allowed")
	}

	seeker, ok := file.(io.Seeker)
	if !ok {
		return fmt.Errorf("Unable to reset uploaded file")
	}

	_, err := seeker.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("Unable to reset uploaded file")
	}

	return nil
}
