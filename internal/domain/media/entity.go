package media

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidMimeType = errors.New("unsupported media format")
	ErrAlreadyExist    = errors.New("media with given id already exist")
	ErrAccountNotFound = errors.New("account with given account_id not found")
	ErrNotFound        = errors.New("media not found")
)

const MaxUploadSize = 5 << 20 // 5 MB

type MimeType string

const (
	TypeImageJPEG = MimeType("image/jpeg")
	TypeImagePNG  = MimeType("image/png")
	TypeAudioMP4  = MimeType("audio/mp4")
	TypeAudioAAC  = MimeType("audio/aac")
	TypeAudioMPEG = MimeType("audio/mpeg")
)

func (t MimeType) valid() bool {
	switch t {
	case TypeImageJPEG, TypeImagePNG, TypeAudioAAC, TypeAudioMP4, TypeAudioMPEG:
		return true
	}

	return false
}

type Media struct {
	Id        string    `json:"id"`
	URL       string    `json:"url"`
	Type      MimeType  `json:"type"`
	AccountId string    `json:"-"`
	CreatedAt time.Time `json:"-"`
}

func (m *Media) generateId() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	m.Id = id.String()

	return nil
}

func (m *Media) filenameFromId(fn string) string {
	return fmt.Sprintf("%s-%s", m.Id, fn)
}

func (m *Media) newTempFile(filename string, buf []byte) (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	if _, err = io.Copy(f, bytes.NewReader(buf)); err != nil {
		return nil, fmt.Errorf("io.Copy: %w", err)
	}

	return f, nil
}

func (m *Media) deleteTempFile(filename string) error {
	return os.Remove(filename)
}
