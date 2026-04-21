package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"mime/multipart"
	"net/http"
	"time"

	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/model"
)

func SendToWisper(data []byte, filename string) (model.Respond, error) {
	body := &bytes.Buffer{}
	// create multipart writer for body
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return model.Respond{}, err
	}

	_, err = part.Write(data)
	if err != nil {
		return model.Respond{}, err
	}
	writer.Close()

	// creating request to transribe service to get the lyrics with timestamps
	req, err := http.NewRequest(
		"POST",
		"http://transcription-service:5001/transcribe",
		body,
	)

	if err != nil {
		return model.Respond{}, err
	}

	// setting headers for content type of formdata
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 500 * time.Second,
	}

	// sending the created request
	res, err := client.Do(req)
	if err != nil {
		return model.Respond{}, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return model.Respond{}, fmt.Errorf("transcription service error: %s\n%v", res.Status, err)
	}

	// return the josn slice
	var lyricsBody model.Respond

	err = json.NewDecoder(res.Body).Decode(&lyricsBody)
	if err != nil {
		return model.Respond{}, err
	}

	// // get the real response transcriptc
	// bodyBytes, _ := io.ReadAll(res.Body)
	// fmt.Println(string(bodyBytes))

	return lyricsBody, nil
}
