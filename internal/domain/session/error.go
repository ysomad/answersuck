package session

import "errors"

var (
	ErrCannotBeTerminated    = errors.New("current session cannot be terminated, use logout instead")
	ErrContextNotFound       = errors.New("session not found in context")
	ErrDeviceMismatch        = errors.New("device doesn't match with device of current session")
	ErrDeviceContextNotFound = errors.New("device not found in context")
	ErrAlreadyExist          = errors.New("session with given id already exist")
	ErrAccountNotFound       = errors.New("session cannot be created, account with given account id is not found")
	ErrNotFound              = errors.New("session not found")
	ErrNotDeleted            = errors.New("session has not been deleted")
)
