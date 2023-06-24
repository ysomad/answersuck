package paging

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

var errInvalidToken = errors.New("paging: invalid token, must store id and RFC3339 datetime separated with comma")

type Token string

// NewToken encodes primary key and time into base64 string separated with ",".
func NewToken(pk string, t time.Time) Token {
	return Token(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%s", pk, t.Format(time.RFC3339Nano)))))
}

// decodeToken decodes page token into string uuid and time,
// token must be a string splitted with "," and encoded into base64,
// not decoded token example: `ef6e33ec-5b1d-4ade-8dcc-b508262ee859,2006-01-02T15:04:05.999999999Z07:00`.
func (t Token) Decode() (string, time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return "", time.Time{}, err
	}

	parts := strings.Split(string(b), ",")
	if len(parts) != 2 {
		return "", time.Time{}, errInvalidToken
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, parts[1])
	if err != nil {
		return "", time.Time{}, err
	}

	return parts[0], parsedTime, nil
}
