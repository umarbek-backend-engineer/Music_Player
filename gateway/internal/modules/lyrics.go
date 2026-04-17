package modules

type AddLyricsPayload struct {
	MusicID string `json:"music_id" binding:"required"`
	Text    string `json:"text" binding:"required"`
}
