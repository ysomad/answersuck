package media

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidType      = errors.New("unsupported media format")
	ErrAlreadyExist     = errors.New("media with given id already exist")
	ErrAccountNotFound  = errors.New("account with given account_id not found")
	ErrNotFound         = errors.New("media not found")
	ErrEmptyFilename    = errors.New("empty filename")
	ErrInvalidAccountId = errors.New("invalid account id")
)

const MaxUploadSize = 5 << 20 // 5 MB

type Type string

const (
	TypeImageJPEG = Type("image/jpeg")
	TypeImagePNG  = Type("image/png")
	TypeImageWEBP = Type("image/webp")
	TypeAudioMP4  = Type("audio/mp4")
	TypeAudioAAC  = Type("audio/aac")
	TypeAudioMPEG = Type("audio/mpeg")
)

func (t Type) valid() bool {
	switch t {
	case TypeImageJPEG, TypeImagePNG, TypeImageWEBP, TypeAudioMP4, TypeAudioMPEG, TypeAudioAAC:
		return true
	}
	return false
}

type Media struct {
	Id        string
	Filename  string
	Type      Type
	AccountId string
	CreatedAt time.Time
}

func New(filename, accountId string, t Type) (Media, error) {
	if filename == "" {
		return Media{}, ErrEmptyFilename
	}

	if !t.valid() {
		return Media{}, ErrInvalidType
	}

	_, err := uuid.Parse(accountId)
	if err != nil {
		return Media{}, ErrInvalidAccountId
	}

	uuid4, err := uuid.NewRandom()
	if err != nil {
		return Media{}, err
	}
	id := uuid4.String()

	return Media{
		Id:        id,
		Filename:  fmt.Sprintf("%s-%s", id, filename),
		Type:      t,
		AccountId: accountId,
		CreatedAt: time.Now(),
	}, nil
}

func (m Media) removeTmpFile() { os.Remove(m.Filename) }
