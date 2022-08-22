package media

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ysomad/answersuck-backend/internal/pkg/mime"
)

var (
	ErrAlreadyExist     = errors.New("media with given id already exist")
	ErrAccountNotFound  = errors.New("account with given account_id not found")
	ErrNotFound         = errors.New("media not found")
	ErrEmptyFilename    = errors.New("empty filename")
	ErrInvalidAccountId = errors.New("invalid account id")
)

const MaxUploadSize = 5 << 20 // 5 MB

type Media struct {
	Id        string
	Filename  string
	Type      mime.Type
	AccountId string
	CreatedAt time.Time
}

func New(filename, accountId string, t mime.Type) (Media, error) {
	if filename == "" {
		return Media{}, ErrEmptyFilename
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
