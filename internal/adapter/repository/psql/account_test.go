package psql

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/pkg/postgres"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountRepo(t *testing.T) {
	type args struct {
		l *zap.Logger
		c *postgres.Client
	}

	tests := []struct {
		name string
		args args
		want *accountRepo
	}{
		{
			name: "must set args",
			args: args{l: _testLogger, c: _testClient},
			want: &accountRepo{l: _testLogger, c: _testClient},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAccountRepo(tt.args.l, tt.args.c)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_accountRepo_Save(t *testing.T) {
	type args struct {
		ctx  context.Context
		a    *account.Account
		code string
	}
	r := accountRepo{
		l: _testLogger,
		c: _testClient,
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "account saved",
			args: args{
				ctx: context.Background(),
				a: &account.Account{
					Email:    "test@test.ru",
					Password: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Save(tt.args.ctx, tt.args.a, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("accountRepo.Save() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountRepo_FindById(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx       context.Context
		accountId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *account.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			got, err := r.FindById(tt.args.ctx, tt.args.accountId)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountRepo.FindById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountRepo_FindByEmail(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *account.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			got, err := r.FindByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountRepo.FindByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountRepo_FindByNickname(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx      context.Context
		nickname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *account.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			got, err := r.FindByNickname(tt.args.ctx, tt.args.nickname)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.FindByNickname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountRepo.FindByNickname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountRepo_Archive(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx       context.Context
		accountId string
		updatedAt time.Time
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
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			if err := r.Archive(tt.args.ctx, tt.args.accountId, tt.args.updatedAt); (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.Archive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountRepo_SavePasswordToken(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx   context.Context
		email string
		token string
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
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			if err := r.SavePasswordToken(tt.args.ctx, tt.args.email, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.SavePasswordToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountRepo_FindPasswordToken(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    account.PasswordToken
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			got, err := r.FindPasswordToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.FindPasswordToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountRepo.FindPasswordToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountRepo_FindEmailByNickname(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx      context.Context
		nickname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			got, err := r.FindEmailByNickname(tt.args.ctx, tt.args.nickname)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.FindEmailByNickname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("accountRepo.FindEmailByNickname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountRepo_SetPassword(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx context.Context
		dto account.SetPasswordDTO
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
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			if err := r.SetPassword(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.SetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountRepo_Verify(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx       context.Context
		code      string
		updatedAt time.Time
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
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			if err := r.Verify(tt.args.ctx, tt.args.code, tt.args.updatedAt); (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountRepo_FindVerification(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c *postgres.Client
	}
	type args struct {
		ctx       context.Context
		accountId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    account.Verification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &accountRepo{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			got, err := r.FindVerification(tt.args.ctx, tt.args.accountId)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountRepo.FindVerification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountRepo.FindVerification() = %v, want %v", got, tt.want)
			}
		})
	}
}
