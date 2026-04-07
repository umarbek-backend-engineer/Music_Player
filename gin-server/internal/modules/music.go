package modules

type MusicChunk struct {
	UploadID string
	FileName string
	ChunkData []byte
	EOF bool
	Progres float32 // optional, shows what percentage of the file is uploaded    -    Calculate Progress = bytes_sent / total_bytes
}