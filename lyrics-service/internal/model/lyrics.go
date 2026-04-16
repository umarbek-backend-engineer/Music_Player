package model

type segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type Respond struct {
	Lyrics   []segment `json:"lyrics"`
	Language string    `json:"language"`
}
