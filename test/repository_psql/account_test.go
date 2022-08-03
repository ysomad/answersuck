package repository_psql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/answersuck/host/internal/adapter/repository/psql"
	"github.com/answersuck/host/internal/domain/account"
	"github.com/answersuck/host/internal/pkg/strings"
)

var _accountRepo *psql.AccountRepo

func insertTestAccount(a account.Account) (account.Account, error) {
	u := strings.NewRandom(10)
	a.Email = fmt.Sprintf("test%s@mail.com", u)
	a.Nickname = "test" + u
	if a.Password == "" {
		a.Password = u
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	err := _accountRepo.Pool.QueryRow(
		context.Background(),
		"INSERT INTO account(email, nickname, password, is_verified, is_archived, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6,$7) RETURNING id",
		a.Email, a.Nickname, a.Password, a.Verified, a.Archived, a.CreatedAt, a.UpdatedAt,
	).Scan(&a.Id)
	return a, err
}

func insertTestVerifCode(accountId string) (string, error) {
	code, err := strings.NewUnique(account.VerifCodeLen)
	if err != nil {
		return "", err
	}

	_, err = _accountRepo.Pool.Exec(
		context.Background(),
		"INSERT INTO verification(code, account_id) VALUES ($1, $2)",
		code, accountId,
	)
	if err != nil {
		return "", err
	}

	return code, nil
}

func insertTestPasswordToken(accountId string, createdAt time.Time) (string, error) {
	t, err := strings.NewUnique(account.PasswordTokenLen)
	if err != nil {
		return "", err
	}

	_, err = _accountRepo.Pool.Exec(
		context.Background(),
		"INSERT INTO password_token(token, account_id, created_at) VALUES ($1, $2, $3)",
		t, accountId, createdAt,
	)
	if err != nil {
		return "", err
	}

	return t, nil
}

func TestAccountRepo_Save(t *testing.T) {
	t.Parallel()

	code, err := strings.NewUnique(account.VerifCodeLen)
	assert.NoError(t, err)

	now := time.Now()

	type args struct {
		ctx  context.Context
		a    account.Account
		code string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "account created",
			args: args{
				ctx: context.Background(),
				a: account.Account{
					Email:     "savetest@test.com",
					Nickname:  "savetest",
					Password:  "secret",
					Verified:  true,
					CreatedAt: now,
					UpdatedAt: now,
				},
				code: code,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "account already exist",
			args: args{
				ctx: context.Background(),
				a:   account.Account{Email: "savetest@test.com"},
			},
			wantErr: true,
			err:     account.ErrAlreadyExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.Save(tt.args.ctx, tt.args.a, tt.args.code)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)
			assert.NotEqual(t, "", got.Id)
			assert.Equal(t, tt.args.a.Email, got.Email)
			assert.Equal(t, tt.args.a.Nickname, got.Nickname)
			assert.Equal(t, tt.args.a.Verified, got.Verified)
			assert.Equal(t, now, got.CreatedAt)
			assert.Equal(t, now, got.UpdatedAt)

			var code string
			err = _accountRepo.Pool.QueryRow(
				context.Background(),
				"SELECT code FROM verification WHERE account_id = $1",
				got.Id,
			).Scan(&code)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.code, code)
		})
	}
}

func TestAccountRepo_FindById(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	a1, err := insertTestAccount(account.Account{Archived: true})
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		accountId string
	}
	tests := []struct {
		name    string
		args    args
		want    account.Account
		wantErr bool
		err     error
	}{
		{
			name: "account found",
			args: args{
				ctx:       context.Background(),
				accountId: a.Id,
			},
			want:    a,
			wantErr: false,
			err:     nil,
		},
		{
			name: "account doesn't exist",
			args: args{
				ctx:       context.Background(),
				accountId: "0eafd279-aa51-46e0-a1e2-3d63c1f6063e",
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
		{
			name: "archived account not found",
			args: args{
				ctx:       context.Background(),
				accountId: a1.Id,
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.FindById(tt.args.ctx, tt.args.accountId)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			// hardcoded because assert.Equal is not working well when comparing structs with time.Time
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Id, got.Id)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.Equal(t, tt.want.Nickname, got.Nickname)
			assert.Equal(t, tt.want.Password, got.Password)
			assert.Equal(t, tt.want.Verified, got.Verified)
			assert.Equal(t, tt.want.Archived, got.Archived)
			assert.Equal(t, tt.want.CreatedAt.Unix(), got.CreatedAt.Unix())
			assert.Equal(t, tt.want.UpdatedAt.Unix(), got.UpdatedAt.Unix())
		})
	}
}

func TestAccountRepo_FindByEmail(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	a1, err := insertTestAccount(account.Account{Archived: true})
	assert.NoError(t, err)

	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    account.Account
		wantErr bool
		err     error
	}{
		{
			name: "account found",
			args: args{
				ctx:   context.Background(),
				email: a.Email,
			},
			want:    a,
			wantErr: false,
			err:     nil,
		},
		{
			name: "account doesn't exist",
			args: args{
				ctx:   context.Background(),
				email: "yeet@test.com",
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
		{
			name: "archived account not found",
			args: args{
				ctx:   context.Background(),
				email: a1.Email,
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.FindByEmail(tt.args.ctx, tt.args.email)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			// hardcoded because assert.Equal is not working well when comparing structs with time.Time
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Id, got.Id)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.Equal(t, tt.want.Nickname, got.Nickname)
			assert.Equal(t, tt.want.Verified, got.Verified)
			assert.Equal(t, tt.want.Archived, got.Archived)
			assert.Equal(t, tt.want.CreatedAt.Unix(), got.CreatedAt.Unix())
			assert.Equal(t, tt.want.UpdatedAt.Unix(), got.UpdatedAt.Unix())
		})
	}
}

func TestAccountRepo_FindByNickname(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	a1, err := insertTestAccount(account.Account{Archived: true})
	assert.NoError(t, err)

	type args struct {
		ctx      context.Context
		nickname string
	}
	tests := []struct {
		name    string
		args    args
		want    account.Account
		wantErr bool
		err     error
	}{
		{
			name: "account found",
			args: args{
				ctx:      context.Background(),
				nickname: a.Nickname,
			},
			want:    a,
			wantErr: false,
			err:     nil,
		},
		{
			name: "account doesn't exist",
			args: args{
				ctx:      context.Background(),
				nickname: "yeet",
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
		{
			name: "archived account not found",
			args: args{
				ctx:      context.Background(),
				nickname: a1.Nickname,
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.FindByNickname(tt.args.ctx, tt.args.nickname)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			// hardcoded because assert.Equal is not working well when comparing structs with time.Time
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Id, got.Id)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.Equal(t, tt.want.Nickname, got.Nickname)
			assert.Equal(t, tt.want.Verified, got.Verified)
			assert.Equal(t, tt.want.Archived, got.Archived)
			assert.Equal(t, tt.want.CreatedAt.Unix(), got.CreatedAt.Unix())
			assert.Equal(t, tt.want.UpdatedAt.Unix(), got.UpdatedAt.Unix())
		})
	}
}

func TestAccountRepo_SetArchived(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	a1, err := insertTestAccount(account.Account{Archived: true})
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		accountId string
		archived  bool
		updatedAt time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "account archived",
			args: args{
				ctx:       context.Background(),
				accountId: a.Id,
				archived:  true,
				updatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "account restored",
			args: args{
				ctx:       context.Background(),
				accountId: a1.Id,
				archived:  false,
				updatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "account not found",
			args: args{
				ctx:       context.Background(),
				accountId: "1bede1cb-9d49-4eb9-b2d1-f053531c565a",
				archived:  true,
				updatedAt: time.Now(),
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := _accountRepo.SetArchived(tt.args.ctx, tt.args.accountId, tt.args.archived, tt.args.updatedAt)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			var archived bool
			err = _accountRepo.Pool.QueryRow(
				context.Background(),
				"SELECT is_archived FROM account WHERE id = $1",
				tt.args.accountId,
			).Scan(&archived)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.archived, archived)
		})
	}
}

func TestAccountRepo_FindPasswordById(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{Password: "password"})
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		accountId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		err     error
	}{
		{
			name:    "account password found",
			args:    args{ctx: context.Background(), accountId: a.Id},
			want:    a.Password,
			wantErr: false,
			err:     nil,
		},
		{
			name: "account not found",
			args: args{
				ctx:       context.Background(),
				accountId: "22e48a5a-aa64-42e7-9e23-335f77962555", // doesnt exist
			},
			want:    "",
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.FindPasswordById(tt.args.ctx, tt.args.accountId)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccountRepo_UpdatePassword(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{Password: "test"})
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		accountId string
		password  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name:    "account password successfully updated",
			args:    args{ctx: context.Background(), accountId: a.Id, password: a.Password},
			wantErr: false,
			err:     nil,
		},
		{
			name: "account not found",
			args: args{
				ctx:       context.Background(),
				accountId: "22e48a5a-aa64-42e7-9e23-335f77962555", // doesnt exist
				password:  "",
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := _accountRepo.UpdatePassword(tt.args.ctx, tt.args.accountId, tt.args.password)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestAccountRepo_SavePasswordToken(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	a1, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	token, err := strings.NewUnique(account.PasswordTokenLen)
	assert.NoError(t, err)

	token1, err := strings.NewUnique(account.PasswordTokenLen)
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		login     string
		token     string
		createdAt time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		err     error
	}{
		{
			name: "password token saved by email",
			args: args{
				ctx:       context.Background(),
				login:     a.Email,
				token:     token,
				createdAt: time.Now(),
			},
			want:    a.Email,
			wantErr: false,
			err:     nil,
		},
		{
			name: "password token saved by nickname",
			args: args{
				ctx:       context.Background(),
				login:     a1.Nickname,
				token:     token1,
				createdAt: time.Now(),
			},
			want:    a1.Email,
			wantErr: false,
			err:     nil,
		},
		{
			name: "password token already exist",
			args: args{
				ctx:   context.Background(),
				login: a.Email,
				token: token,
			},
			want:    a.Email,
			wantErr: true,
			err:     account.ErrPasswordTokenAlreadyExist,
		},
		{
			name: "account with email not found",
			args: args{
				ctx:   context.Background(),
				login: "doesntexist@mail.com",
				token: token,
			},
			want:    a.Email,
			wantErr: true,
			err:     account.ErrNotFound,
		},
		{
			name: "account with nickname not found",
			args: args{
				ctx:   context.Background(),
				login: "doesntexist",
				token: token,
			},
			want:    a.Email,
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.SavePasswordToken(tt.args.ctx, account.SavePasswordTokenDTO{
				Login: tt.args.login,
				Token: tt.args.token,
			})
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			var (
				tokenFromDB string
				createdAt   time.Time
			)
			err = _accountRepo.Pool.QueryRow(
				context.Background(),
				`SELECT token, created_at
FROM password_token
WHERE account_id = (SELECT id FROM account WHERE email = $1 OR nickname = $2)`,
				tt.args.login, tt.args.login,
			).Scan(&tokenFromDB, &createdAt)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.token, tokenFromDB)
			assert.Equal(t, tt.args.createdAt.Unix(), createdAt.Unix())
		})
	}
}

func TestAccountRepo_FindPasswordToken(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	now := time.Now()
	token, err := insertTestPasswordToken(a.Id, time.Now())
	assert.NoError(t, err)

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    account.PasswordToken
		wantErr bool
		err     error
	}{
		{
			name: "password token found",
			args: args{
				ctx:   context.Background(),
				token: token,
			},
			want: account.PasswordToken{
				AccountId: a.Id,
				Token:     token,
				CreatedAt: now,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "password token not found",
			args: args{
				ctx:   context.Background(),
				token: "doesntexist",
			},
			wantErr: true,
			err:     account.ErrPasswordTokenNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.FindPasswordToken(tt.args.ctx, tt.args.token)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.want.AccountId, got.AccountId)
			assert.Equal(t, tt.want.Token, got.Token)
			assert.Equal(t, tt.want.CreatedAt.Unix(), got.CreatedAt.Unix())
		})
	}
}

func TestAccountRepo_SetPassword(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	token, err := insertTestPasswordToken(a.Id, time.Now())
	assert.NoError(t, err)

	type args struct {
		ctx context.Context
		dto account.SetPasswordDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "password succesfully set",
			args: args{
				ctx: context.Background(),
				dto: account.SetPasswordDTO{
					AccountId: a.Id,
					Password:  "newsecret",
					Token:     token,
					UpdatedAt: time.Now(),
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "password not set with empty token",
			args: args{
				ctx: context.Background(),
				dto: account.SetPasswordDTO{
					AccountId: a.Id,
					Password:  "newsecret",
					UpdatedAt: time.Now(),
				},
			},
			wantErr: true,
			err:     account.ErrPasswordNotSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := _accountRepo.SetPassword(tt.args.ctx, tt.args.dto)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			var (
				password  string
				updatedAt time.Time
			)
			err = _accountRepo.Pool.QueryRow(
				context.Background(),
				"SELECT password, updated_at FROM account WHERE id = $1",
				tt.args.dto.AccountId,
			).Scan(&password, &updatedAt)
			assert.NoError(t, err)

			var exists bool
			err = _accountRepo.Pool.QueryRow(
				context.Background(),
				"SELECT EXISTS(SELECT 1 FROM password_token WHERE token = $1)",
				tt.args.dto.Token,
			).Scan(&exists)
			assert.NoError(t, err)
			assert.Equal(t, false, exists)
		})
	}
}

func TestAccountRepo_Verify(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	code, err := insertTestVerifCode(a.Id)
	assert.NoError(t, err)

	a1, err := insertTestAccount(account.Account{Verified: true})
	assert.NoError(t, err)

	code1, err := insertTestVerifCode(a1.Id)
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		code      string
		updatedAt time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "account successfully verified",
			args: args{
				ctx:       context.Background(),
				code:      code,
				updatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "account already verified",
			args: args{
				ctx:       context.Background(),
				code:      code1,
				updatedAt: time.Now(),
			},
			wantErr: true,
			err:     account.ErrAlreadyVerified,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := _accountRepo.Verify(tt.args.ctx, tt.args.code, tt.args.updatedAt)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			var (
				verified  bool
				updatedAt time.Time
			)
			err = _accountRepo.Pool.QueryRow(
				context.Background(),
				"SELECT is_verified, updated_at FROM account WHERE id = $1",
				a.Id,
			).Scan(&verified, &updatedAt)
			assert.NoError(t, err)
			assert.Equal(t, true, verified)
			assert.Equal(t, tt.args.updatedAt.Unix(), updatedAt.Unix())
		})
	}
}

func TestAccountRepo_FindVerification(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{Verified: true})
	assert.NoError(t, err)

	code, err := insertTestVerifCode(a.Id)
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		accountId string
	}
	tests := []struct {
		name    string
		args    args
		want    account.Verification
		wantErr bool
		err     error
	}{
		{
			name: "verification found",
			args: args{
				ctx:       context.Background(),
				accountId: a.Id,
			},
			want: account.Verification{
				Email:    a.Email,
				Code:     code,
				Verified: a.Verified,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "verification not found",
			args: args{
				ctx:       context.Background(),
				accountId: "6a3daeca-9036-4f37-8c2c-ee29c9bebc97", // doesnt exist
			},
			want: account.Verification{
				Email:    a.Email,
				Code:     code,
				Verified: a.Verified,
			},
			wantErr: true,
			err:     account.ErrVerificationNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _accountRepo.FindVerification(tt.args.ctx, tt.args.accountId)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
