package modules

type AddLyricsPayload struct {
	MusicID string `json:"music_id" binding:"required"`
	Text    string `json:"text" binding:"required"`
}

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type Respond struct {
	Lyrics   []Segment `json:"lyrics"`
	Language string    `json:"language"`
}
