package modelos

// RespuestaMetadataAudioDTO is the data transfer object for search responses
type RespuestaMetadataAudioDTO struct {
	ObjAudio MetadataAudio // Found audio information
	Codigo   int           // HTTP-like status code (e.g., 200, 404)
	Mensaje  string        // Result description
}
