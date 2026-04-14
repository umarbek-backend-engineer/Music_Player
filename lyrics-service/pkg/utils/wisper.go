package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func SendToWisper(data []byte, filename string) (string, error) {
	body := &bytes.Buffer{}
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

	req, err := http.NewRequest(
		"POST",
		"http://transcription-service:5001/transcribe",
		body,
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	return string(resBody), nil
}
