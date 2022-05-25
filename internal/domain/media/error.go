package media

import "errors"

var (
	ErrInvalidMimeType = errors.New("unsupported media format")
	ErrAlreadyExist    = errors.New("media with given id already exist")
	ErrAccountNotFound = errors.New("account with given account_id not found")
)
