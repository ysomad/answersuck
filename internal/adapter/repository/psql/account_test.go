package psql

import (
	"context"
	"testing"

	"github.com/answersuck/host/internal/pkg/postgres"
	"go.uber.org/zap"
)

func TestAccountRepo_UpdatePassword(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Client *postgres.Client
	}
	type args struct {
		ctx       context.Context
		accountId string
		password  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AccountRepo{
				Logger: tt.fields.Logger,
				Client: tt.fields.Client,
			}
			if err := r.UpdatePassword(tt.args.ctx, tt.args.accountId, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("AccountRepo.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
