package paging

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

var ErrInvalidToken = errors.New("invalid page token, must be valid base64 or empty string")

type UnsortableToken string

// NewUnsortableToken encodes non-sortable id and time into base64 string separated with ",".
func NewUnsortableToken(id string, t time.Time) UnsortableToken {
	s := fmt.Sprintf("%s,%s", id, t.Format(time.RFC3339Nano))

	return UnsortableToken(base64.StdEncoding.EncodeToString([]byte(s)))
}

// Decode decodes page token into string id and time,
// token must be a string splitted with "," and encoded into base64/
// Encoded token example: `ef6e33ec-5b1d-4ade-8dcc-b508262ee859,2006-01-02T15:04:05.999999999Z07:00`.
func (t UnsortableToken) Decode() (string, time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return "", time.Time{}, ErrInvalidToken
	}

	parts := strings.Split(string(b), ",")
	if len(parts) != 2 {
		return "", time.Time{}, ErrInvalidToken
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, parts[1])
	if err != nil {
		return "", time.Time{}, ErrInvalidToken
	}

	return parts[0], parsedTime, nil
}

type OffsetToken string

// NewOffsetToken encodes limit and offset into base64 string.
func NewOffsetToken(limit, offset uint64) OffsetToken {
	limitBytes := uint64ToBytes(limit)
	offsetBytes := uint64ToBytes(offset)
	return OffsetToken(base64.StdEncoding.EncodeToString(append(limitBytes, offsetBytes...)))
}

// Decode decodes token into limit and offset.
// If token is empty returns default values and no error.
// If token is invalid, returns error.
func (t OffsetToken) Decode() (limit, offset uint64, err error) {
	if t == "" {
		return 0, 0, nil
	}

	b, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return 0, 0, ErrInvalidToken
	}

	if len(b) != 16 {
		return 0, 0, ErrInvalidToken
	}

	return bytesToUint64(b[:8]), bytesToUint64(b[8:]), nil
}

func uint64ToBytes(i uint64) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Len = 8
	sh.Cap = 8
	sh.Data = uintptr(unsafe.Pointer(&i))

	return b[:]
}

func bytesToUint64(b []byte) uint64 {
	return *(*uint64)(unsafe.Pointer(&b[0]))
}
