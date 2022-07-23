package v1

import (
	"context"
	"errors"

	"github.com/answersuck/vault/internal/domain/session"
)

type (
	accountIdCtxKey struct{}
	sessionIdCtxKey struct{}
	deviceCtxKey    struct{}
)

var (
	errAccountIdNotFound = errors.New("account id not found in context")
	errSessionIdNotFound = errors.New("session id not found in context")
	errDeviceNotFound    = errors.New("device not found in context")
)

func getAccountId(c context.Context) (string, error) {
	v := c.Value(accountIdCtxKey{})

	aid, ok := v.(string)
	if !ok || aid == "" {
		return "", errAccountIdNotFound
	}

	return aid, nil
}

func getSessionId(c context.Context) (string, error) {
	v := c.Value(sessionIdCtxKey{})

	sid, ok := v.(string)
	if !ok || sid == "" {
		return "", errSessionIdNotFound
	}

	return sid, nil
}

func getDevice(c context.Context) (session.Device, error) {
	v := c.Value(deviceCtxKey{})

	d, ok := v.(session.Device)
	if !ok {
		return session.Device{}, errDeviceNotFound
	}

	return d, nil
}
