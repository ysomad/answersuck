package psql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/answersuck/vault/internal/adapter/repository/psql"
	"github.com/answersuck/vault/internal/domain/account"

	"github.com/answersuck/vault/pkg/strings"
)

var accountRepo *psql.AccountRepo

type accountRepoTestSuite struct {
	suite.Suite
}

func TestAccountRepoTestSuite(t *testing.T) { suite.Run(t, new(accountRepoTestSuite)) }

func (s *accountRepoTestSuite) insertTestAccount(a account.Account) account.Account {
	err := postgresClient.Pool.QueryRow(
		context.Background(),
		"INSERT INTO account(email, nickname, password, is_verified, is_archived, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6,$7) RETURNING id",
		a.Email, a.Nickname, a.Password, a.Verified, a.Archived, a.CreatedAt, a.UpdatedAt,
	).Scan(&a.Id)
	s.NoError(err)
	return a
}

func (s *accountRepoTestSuite) insertTestVerifCode(accountId string) string {
	code, err := strings.NewUnique(account.VerifCodeLen)
	s.NoError(err)

	_, err = postgresClient.Pool.Exec(
		context.Background(),
		"INSERT INTO verification(code, account_id) VALUES ($1, $2)",
		code, accountId,
	)
	s.NoError(err)
	return code
}

func (s *accountRepoTestSuite) insertTestPasswordToken(accountId string, createdAt time.Time) string {
	t, err := strings.NewUnique(account.PasswordTokenLen)
	s.NoError(err)

	_, err = postgresClient.Pool.Exec(
		context.Background(),
		"INSERT INTO password_token(token, account_id, created_at) VALUES ($1, $2, $3)",
		t, accountId, createdAt,
	)
	s.NoError(err)
	return t
}

func (s *accountRepoTestSuite) TestSave() {
	code, err := strings.NewUnique(account.VerifCodeLen)
	s.NoError(err)

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
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := accountRepo.Save(tt.args.ctx, tt.args.a, tt.args.code)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NotEqual(t, "", got.Id)
			assert.Equal(t, tt.args.a.Email, got.Email)
			assert.Equal(t, tt.args.a.Nickname, got.Nickname)
			assert.Equal(t, tt.args.a.Verified, got.Verified)
			assert.Equal(t, now, got.CreatedAt)
			assert.Equal(t, now, got.UpdatedAt)

			var code string
			err = postgresClient.Pool.QueryRow(
				context.Background(),
				"SELECT code FROM verification WHERE account_id = $1",
				got.Id,
			).Scan(&code)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.code, code)
		})
	}
}

func (s *accountRepoTestSuite) TestFindById() {
	now := time.Now()
	a := s.insertTestAccount(account.Account{
		Email:     "findbyid@test.com",
		Nickname:  "findbyid",
		Password:  "secret",
		Verified:  false,
		Archived:  false,
		CreatedAt: now,
		UpdatedAt: now,
	})
	a2 := s.insertTestAccount(account.Account{
		Email:     "findbyid1@test.com",
		Nickname:  "findbyid1",
		Password:  "secret",
		Verified:  false,
		Archived:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

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
				accountId: a2.Id,
			},
			wantErr: true,
			err:     account.ErrNotFound,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := accountRepo.FindById(tt.args.ctx, tt.args.accountId)
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

func (s *accountRepoTestSuite) TestFindByEmail() {
	now := time.Now()
	a := s.insertTestAccount(account.Account{
		Email:     "findbyemail@test.com",
		Nickname:  "findbyemail",
		Password:  "secret",
		Verified:  false,
		Archived:  false,
		CreatedAt: now,
		UpdatedAt: now,
	})
	a1 := s.insertTestAccount(account.Account{
		Email:     "findbyemail1@test.com",
		Nickname:  "findbyemail1",
		Password:  "secret",
		Verified:  false,
		Archived:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

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
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := accountRepo.FindByEmail(tt.args.ctx, tt.args.email)
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

func (s *accountRepoTestSuite) TestFindByNickname() {
	now := time.Now()
	a := s.insertTestAccount(account.Account{
		Email:     "findbynickname@test.com",
		Nickname:  "findbynickname",
		Password:  "secret",
		Verified:  false,
		Archived:  false,
		CreatedAt: now,
		UpdatedAt: now,
	})
	a1 := s.insertTestAccount(account.Account{
		Email:     "findbynickname1@test.com",
		Nickname:  "findbynickname1",
		Password:  "secret",
		Verified:  false,
		Archived:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

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
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := accountRepo.FindByNickname(tt.args.ctx, tt.args.nickname)
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

func (s *accountRepoTestSuite) TestArchive() {
	a := s.insertTestAccount(account.Account{
		Email:    "archivetest@test.com",
		Nickname: "archivetest",
		Archived: false,
	})

	type args struct {
		ctx       context.Context
		accountId string
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
				updatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := accountRepo.Archive(tt.args.ctx, tt.args.accountId, tt.args.updatedAt)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			var archived bool
			err = postgresClient.Pool.QueryRow(
				context.Background(),
				"SELECT is_archived FROM account WHERE id = $1",
				tt.args.accountId,
			).Scan(&archived)
			assert.NoError(t, err)
			assert.Equal(t, true, archived)
		})
	}
}

func (s *accountRepoTestSuite) TestSavePasswordToken() {
	a := s.insertTestAccount(account.Account{
		Email:    "savetokentest@test.com",
		Nickname: "savetokentest",
	})
	a1 := s.insertTestAccount(account.Account{
		Email:    "savetokentest1@test.com",
		Nickname: "savetokentest1",
	})

	token, err := strings.NewUnique(account.PasswordTokenLen)
	s.NoError(err)
	token1, err := strings.NewUnique(account.PasswordTokenLen)
	s.NoError(err)

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
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := accountRepo.SavePasswordToken(tt.args.ctx, account.SavePasswordTokenDTO{
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
			err = postgresClient.Pool.QueryRow(
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

func (s *accountRepoTestSuite) TestFindPasswordToken() {
	a := s.insertTestAccount(account.Account{
		Email:     "findpaswordtoken@mail.com",
		Nickname:  "findpasswordtoken",
		CreatedAt: time.Now(),
	})
	now := time.Now()
	t := s.insertTestPasswordToken(a.Id, now)

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
				token: t,
			},
			want: account.PasswordToken{
				AccountId: a.Id,
				Token:     t,
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
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := accountRepo.FindPasswordToken(tt.args.ctx, tt.args.token)
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

func (s *accountRepoTestSuite) TestSetPassword() {
	a := s.insertTestAccount(account.Account{
		Email:    "setpasswordtest@mail.com",
		Nickname: "setpasswordtest@mail.com",
		Password: "forgotten",
	})
	t := s.insertTestPasswordToken(a.Id, time.Now())

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
					Token:     t,
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
		s.T().Run(tt.name, func(t *testing.T) {
			err := accountRepo.SetPassword(tt.args.ctx, tt.args.dto)
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
			err = postgresClient.Pool.QueryRow(
				context.Background(),
				"SELECT password, updated_at FROM account WHERE id = $1",
				tt.args.dto.AccountId,
			).Scan(&password, &updatedAt)
			assert.NoError(t, err)

			var exists bool
			err = postgresClient.Pool.QueryRow(
				context.Background(),
				"SELECT EXISTS(SELECT 1 FROM password_token WHERE token = $1)",
				tt.args.dto.Token,
			).Scan(&exists)
			assert.NoError(t, err)
			assert.Equal(t, false, exists)
		})
	}
}

func (s *accountRepoTestSuite) TestVerify() {
	a := s.insertTestAccount(account.Account{
		Email:    "testverify@mail.com",
		Nickname: "testverify@mail.com",
		Verified: false,
	})
	code := s.insertTestVerifCode(a.Id)

	a1 := s.insertTestAccount(account.Account{
		Email:    "testverify1@mail.com",
		Nickname: "testverify1@mail.com",
		Verified: true,
	})
	code1 := s.insertTestVerifCode(a1.Id)

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
		s.T().Run(tt.name, func(t *testing.T) {
			err := accountRepo.Verify(tt.args.ctx, tt.args.code, tt.args.updatedAt)
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
			err = postgresClient.Pool.QueryRow(
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

func (s *accountRepoTestSuite) TestFindVerification() {
	a := s.insertTestAccount(account.Account{
		Email:    "findverification@mail.com",
		Nickname: "findverification",
		Verified: true,
	})
	code := s.insertTestVerifCode(a.Id)

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
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := accountRepo.FindVerification(tt.args.ctx, tt.args.accountId)
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
