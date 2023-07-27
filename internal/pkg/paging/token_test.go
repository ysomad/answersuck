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
	testToken UnsortableToken
)

func genUnsortableTestToken(pk string, t time.Time) UnsortableToken {
	return UnsortableToken(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%s", pk, t.Format(time.RFC3339Nano)))))
}

func TestMain(m *testing.M) {
	testPK = "80862cb4-947a-4d64-8dbe-858fea7d84f2"
	testTime = time.Now()
	testToken = genUnsortableTestToken(testPK, testTime)

	os.Exit(m.Run())
}

func TestNewUnsortableToken(t *testing.T) {
	type args struct {
		uuid string
		t    time.Time
	}

	tests := []struct {
		name string
		args args
		want UnsortableToken
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
			want: genUnsortableTestToken("", testTime),
		},
		{
			name: "empty time",
			args: args{
				uuid: testPK,
				t:    time.Time{},
			},
			want: genUnsortableTestToken(testPK, time.Time{}),
		},
		{
			name: "empty pk and time",
			args: args{
				uuid: "",
				t:    time.Time{},
			},
			want: genUnsortableTestToken("", time.Time{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUnsortableToken(tt.args.uuid, tt.args.t)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUnsortableToken_Decode(t *testing.T) {
	tests := []struct {
		name         string
		encodedToken UnsortableToken
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

func TestNewOffsetToken(t *testing.T) {
	type args struct {
		limit  uint64
		offset uint64
	}
	tests := []struct {
		name string
		args args
		want OffsetToken
	}{
		{
			name: "success",
			args: args{
				limit:  100,
				offset: 500,
			},
			want: "ZAAAAAAAAAD0AQAAAAAAAA==",
		},
		{
			name: "0 limit and offset",
			args: args{
				limit:  0,
				offset: 0,
			},
			want: "AAAAAAAAAAAAAAAAAAAAAA==",
		},
		{
			name: "max limit and offset",
			args: args{
				limit:  18446744073709551615,
				offset: 18446744073709551615,
			},
			want: "/////////////////////w==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOffsetToken(tt.args.limit, tt.args.offset)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOffsetToken_Decode(t *testing.T) {
	tests := []struct {
		name       string
		tr         OffsetToken
		wantLimit  uint64
		wantOffset uint64
		wantErr    bool
	}{
		{
			name:       "success",
			tr:         "ZAAAAAAAAAD0AQAAAAAAAA==",
			wantLimit:  100,
			wantOffset: 500,
			wantErr:    false,
		},
		{
			name:       "0 limit, 0 offset",
			tr:         "AAAAAAAAAAAAAAAAAAAAAA==",
			wantLimit:  0,
			wantOffset: 0,
			wantErr:    false,
		},
		{
			name:       "max limit and offset",
			tr:         "/////////////////////w==",
			wantLimit:  18446744073709551615,
			wantOffset: 18446744073709551615,
			wantErr:    false,
		},
		{
			name:       "empty token",
			tr:         "",
			wantLimit:  0,
			wantOffset: 0,
			wantErr:    true,
		},
		{
			name:       "invalid base64",
			tr:         "asdasdasdasdasdasdasdasdasdaasdasdasdasdada==",
			wantLimit:  0,
			wantOffset: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset, err := tt.tr.Decode()
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.wantLimit, gotLimit)
			assert.Equal(t, tt.wantOffset, gotOffset)
		})
	}
}
