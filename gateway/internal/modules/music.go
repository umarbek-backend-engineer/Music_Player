package modules

type MusicChunk struct {
	UploadID  string `json:"Uploaded ID"`
	FileName  string `json:"filename"`
	ChunkData []byte `json:"chunk,omitempty"`
	EOF       bool   `json:"eof"`
	// Progres float32 // optional, shows what percentage of the file is uploaded    -    Calculate Progress = bytes_sent / total_bytes
}
