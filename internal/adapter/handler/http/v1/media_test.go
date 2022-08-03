package v1

import (
	"net/http"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func Test_newMediaHandler(t *testing.T) {
	type args struct {
		d *Deps
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMediaHandler(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMediaHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mediaHandler_upload(t *testing.T) {
	type fields struct {
		log      *zap.Logger
		validate validate
		media    mediaService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &mediaHandler{
				log:      tt.fields.log,
				validate: tt.fields.validate,
				media:    tt.fields.media,
			}
			h.upload(tt.args.w, tt.args.r)
		})
	}
}
