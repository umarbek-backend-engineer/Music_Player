package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func SendToWisper(data []byte, filename string) (string, error) {
	body := &bytes.Buffer{}
	// create multipart writer for body
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	_, err = part.Write(data)
	if err != nil {
		return "", err
	}
	writer.Close()

	// creating request to transribe service to get the lyrics with timestamps
	req, err := http.NewRequest(
		"POST",
		"http://transcription-service:5001/transcribe",
		body,
	)

	if err != nil {
		return "", err
	}

	// setting headers for content type of formdata
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	// sending the created request
	res, err := client.Do(req)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("transcription service error: %s\n%v", res.Status, err)
	}
	defer res.Body.Close()

	// Read the entire HTTP response body into memory as a byte slice.
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	// return the byte slice
	return string(resBody), nil
}
