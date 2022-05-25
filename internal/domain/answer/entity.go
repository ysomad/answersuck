package answer

var (
	// see media.MimeType
	allowedMimeTypes = [2]string{"image/jpeg", "image/png"}
)

type Answer struct {
	Id      int    `json:"id"`
	Text    string `json:"text"`
	MediaId string `json:"mediaId"`
}

// isMimeTypeAllowed checks if media for answer is in allowed mime type array
func (a Answer) isMimeTypeAllowed(mt string) bool {
	var allowed bool

	for _, t := range allowedMimeTypes {
		if mt == t {
			allowed = true
			break
		}
	}

	return allowed
}
