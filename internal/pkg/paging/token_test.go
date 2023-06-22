package paging

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testPK    string
	testTime  time.Time
	testToken Token
)

func generateTestToken(pk string, t time.Time) Token {
	return Token(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%s", pk, t.Format(time.RFC3339Nano)))))
}

func TestMain(m *testing.M) {
	testPK = "80862cb4-947a-4d64-8dbe-858fea7d84f2"
	testTime = time.Now()
	testToken = generateTestToken(testPK, testTime)

	os.Exit(m.Run())
}

func TestNewToken(t *testing.T) {
	type args struct {
		uuid string
		t    time.Time
	}

	tests := []struct {
		name string
		args args
		want Token
	}{
		{
			name: "success",
			args: args{
				uuid: testPK,
				t:    testTime,
			},
			want: testToken,
		},
		{
			name: "empty pk",
			args: args{
				uuid: "",
				t:    testTime,
			},
			want: generateTestToken("", testTime),
		},
		{
			name: "empty time",
			args: args{
				uuid: testPK,
				t:    time.Time{},
			},
			want: generateTestToken(testPK, time.Time{}),
		},
		{
			name: "empty pk and time",
			args: args{
				uuid: "",
				t:    time.Time{},
			},
			want: generateTestToken("", time.Time{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewToken(tt.args.uuid, tt.args.t)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToken_Decode(t *testing.T) {
	tests := []struct {
		name         string
		encodedToken Token
		wantPK       string
		wantTime     time.Time
		wantErr      bool
	}{
		{
			name:         "invalid token",
			encodedToken: "asdffg",
			wantPK:       "",
			wantTime:     time.Time{},
			wantErr:      true,
		},
		{
			name:         "success",
			encodedToken: testToken,
			wantPK:       testPK,
			wantTime:     testTime,
			wantErr:      false,
		},
		{
			name:         "invalid token 2",
			encodedToken: "MTIzNDUsYXNkZmcsMTIzNCx3ZXJ0",
			wantPK:       "",
			wantTime:     time.Time{},
			wantErr:      true,
		},
		{
			name:         "invalid token time",
			encodedToken: "cGssMjAyMy0wNC0wOVQwNjMzOjU0KzAwOjAw",
			wantPK:       "",
			wantTime:     time.Time{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPK, gotTime, err := tt.encodedToken.Decode()

			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.wantPK, gotPK)
			assert.ObjectsAreEqual(tt.wantTime, gotTime)
		})
	}
}
