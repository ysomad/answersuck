package entity

import (
	"errors"
	"mime"
	"strings"
	"time"
)

type MediaType int8

const (
	MediaTypeImage MediaType = iota + 1
	MediaTypeAudio
	MediaTypeVideo
)

type Media struct {
	URL        string
	Type       MediaType
	Author     string
	CreateTime time.Time
}

var (
	errInvalidMediaURL      = errors.New("invalid media url")
	errUnsupportedMediaType = errors.New("unsupported media type")
)

func NewMedia(url, author string) (Media, error) {
	parts := strings.Split(url, ".")
	if len(parts) == 0 {
		return Media{}, errInvalidMediaURL
	}

	mimeType := mime.TypeByExtension("." + parts[len(parts)-1])
	if mimeType == "" {
		return Media{}, errUnsupportedMediaType
	}

	var mtype MediaType

	switch mimeType {
	case "image/gif", "image/jpeg", "image/png", "image/svg+xml", "image/webp":
		mtype = MediaTypeImage
	case "audio/wav", "audio/mpeg", "audio/aac", "audio/webm":
		mtype = MediaTypeAudio
	case "video/mp4", "video/mpeg", "video/webm":
		mtype = MediaTypeVideo
	default:
		return Media{}, errUnsupportedMediaType
	}

	return Media{
		URL:        url,
		Type:       mtype,
		Author:     author,
		CreateTime: time.Now(),
	}, nil
}
