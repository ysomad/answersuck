package domain

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrDeviceInvalidUserAgent = errors.New("invalid user agent")
	ErrDeviceInvalidIP        = errors.New("invalid ip address")
)

type Device struct {
	UserAgent string `json:"ua"`
	IP        string `json:"ip"`
}

func (d Device) validateIP() error {
	if net.ParseIP(d.IP) == nil {
		return ErrDeviceInvalidIP
	}

	return nil
}

func (d Device) Validate() error {
	// TODO: add validation for User-Agent
	if d.UserAgent == "" {
		return ErrDeviceInvalidUserAgent
	}

	if err := d.validateIP(); err != nil {
		return fmt.Errorf("d.validateIP: %w", err)
	}

	return nil
}
