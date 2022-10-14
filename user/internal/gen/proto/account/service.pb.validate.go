// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: user/account/service.proto

package account

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// define the regex for a UUID once up-front
var _service_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Account with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Account) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Account with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in AccountMultiError, or nil if none found.
func (m *Account) ValidateAll() error {
	return m.validate(true)
}

func (m *Account) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Email

	// no validation rules for Username

	// no validation rules for EmailVerified

	// no validation rules for Archived

	if all {
		switch v := interface{}(m.GetCreatedTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, AccountValidationError{
					field:  "CreatedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, AccountValidationError{
					field:  "CreatedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AccountValidationError{
				field:  "CreatedTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdatedTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, AccountValidationError{
					field:  "UpdatedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, AccountValidationError{
					field:  "UpdatedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AccountValidationError{
				field:  "UpdatedTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return AccountMultiError(errors)
	}

	return nil
}

// AccountMultiError is an error wrapping multiple validation errors returned
// by Account.ValidateAll() if the designated constraints aren't met.
type AccountMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AccountMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AccountMultiError) AllErrors() []error { return m }

// AccountValidationError is the validation error returned by Account.Validate
// if the designated constraints aren't met.
type AccountValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AccountValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AccountValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AccountValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AccountValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AccountValidationError) ErrorName() string { return "AccountValidationError" }

// Error satisfies the builtin error interface
func (e AccountValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAccount.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AccountValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AccountValidationError{}

// Validate checks the field values on CreateAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateAccountRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateAccountRequestMultiError, or nil if none found.
func (m *CreateAccountRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateAccountRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetEmail()) > 320 {
		err := CreateAccountRequestValidationError{
			field:  "Email",
			reason: "value length must be at most 320 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = CreateAccountRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_CreateAccountRequest_Username_Pattern.MatchString(m.GetUsername()) {
		err := CreateAccountRequestValidationError{
			field:  "Username",
			reason: "value does not match regex pattern \"^[a-zA-Z0-9][\\\\w]{3,24}$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 10 || l > 128 {
		err := CreateAccountRequestValidationError{
			field:  "Password",
			reason: "value length must be between 10 and 128 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateAccountRequestMultiError(errors)
	}

	return nil
}

func (m *CreateAccountRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *CreateAccountRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// CreateAccountRequestMultiError is an error wrapping multiple validation
// errors returned by CreateAccountRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateAccountRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateAccountRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateAccountRequestMultiError) AllErrors() []error { return m }

// CreateAccountRequestValidationError is the validation error returned by
// CreateAccountRequest.Validate if the designated constraints aren't met.
type CreateAccountRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateAccountRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateAccountRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateAccountRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateAccountRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateAccountRequestValidationError) ErrorName() string {
	return "CreateAccountRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateAccountRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateAccountRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateAccountRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateAccountRequestValidationError{}

var _CreateAccountRequest_Username_Pattern = regexp.MustCompile("^[a-zA-Z0-9][\\w]{3,24}$")

// Validate checks the field values on CreateAccountResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateAccountResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateAccountResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateAccountResponseMultiError, or nil if none found.
func (m *CreateAccountResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateAccountResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAccount()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateAccountResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateAccountResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccount()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateAccountResponseValidationError{
				field:  "Account",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreateAccountResponseMultiError(errors)
	}

	return nil
}

// CreateAccountResponseMultiError is an error wrapping multiple validation
// errors returned by CreateAccountResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateAccountResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateAccountResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateAccountResponseMultiError) AllErrors() []error { return m }

// CreateAccountResponseValidationError is the validation error returned by
// CreateAccountResponse.Validate if the designated constraints aren't met.
type CreateAccountResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateAccountResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateAccountResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateAccountResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateAccountResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateAccountResponseValidationError) ErrorName() string {
	return "CreateAccountResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateAccountResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateAccountResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateAccountResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateAccountResponseValidationError{}

// Validate checks the field values on GetAccountByIdRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAccountByIdRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAccountByIdRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAccountByIdRequestMultiError, or nil if none found.
func (m *GetAccountByIdRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAccountByIdRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetAccountId()); err != nil {
		err = GetAccountByIdRequestValidationError{
			field:  "AccountId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetAccountByIdRequestMultiError(errors)
	}

	return nil
}

func (m *GetAccountByIdRequest) _validateUuid(uuid string) error {
	if matched := _service_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// GetAccountByIdRequestMultiError is an error wrapping multiple validation
// errors returned by GetAccountByIdRequest.ValidateAll() if the designated
// constraints aren't met.
type GetAccountByIdRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAccountByIdRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAccountByIdRequestMultiError) AllErrors() []error { return m }

// GetAccountByIdRequestValidationError is the validation error returned by
// GetAccountByIdRequest.Validate if the designated constraints aren't met.
type GetAccountByIdRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAccountByIdRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAccountByIdRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAccountByIdRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAccountByIdRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAccountByIdRequestValidationError) ErrorName() string {
	return "GetAccountByIdRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetAccountByIdRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAccountByIdRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAccountByIdRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAccountByIdRequestValidationError{}

// Validate checks the field values on GetAccountByIdResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAccountByIdResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAccountByIdResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAccountByIdResponseMultiError, or nil if none found.
func (m *GetAccountByIdResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAccountByIdResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAccount()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetAccountByIdResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetAccountByIdResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccount()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetAccountByIdResponseValidationError{
				field:  "Account",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetAccountByIdResponseMultiError(errors)
	}

	return nil
}

// GetAccountByIdResponseMultiError is an error wrapping multiple validation
// errors returned by GetAccountByIdResponse.ValidateAll() if the designated
// constraints aren't met.
type GetAccountByIdResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAccountByIdResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAccountByIdResponseMultiError) AllErrors() []error { return m }

// GetAccountByIdResponseValidationError is the validation error returned by
// GetAccountByIdResponse.Validate if the designated constraints aren't met.
type GetAccountByIdResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAccountByIdResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAccountByIdResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAccountByIdResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAccountByIdResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAccountByIdResponseValidationError) ErrorName() string {
	return "GetAccountByIdResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetAccountByIdResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAccountByIdResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAccountByIdResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAccountByIdResponseValidationError{}

// Validate checks the field values on GetAccountByEmailRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAccountByEmailRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAccountByEmailRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAccountByEmailRequestMultiError, or nil if none found.
func (m *GetAccountByEmailRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAccountByEmailRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = GetAccountByEmailRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetAccountByEmailRequestMultiError(errors)
	}

	return nil
}

func (m *GetAccountByEmailRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *GetAccountByEmailRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// GetAccountByEmailRequestMultiError is an error wrapping multiple validation
// errors returned by GetAccountByEmailRequest.ValidateAll() if the designated
// constraints aren't met.
type GetAccountByEmailRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAccountByEmailRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAccountByEmailRequestMultiError) AllErrors() []error { return m }

// GetAccountByEmailRequestValidationError is the validation error returned by
// GetAccountByEmailRequest.Validate if the designated constraints aren't met.
type GetAccountByEmailRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAccountByEmailRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAccountByEmailRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAccountByEmailRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAccountByEmailRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAccountByEmailRequestValidationError) ErrorName() string {
	return "GetAccountByEmailRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetAccountByEmailRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAccountByEmailRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAccountByEmailRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAccountByEmailRequestValidationError{}

// Validate checks the field values on GetAccountByEmailResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAccountByEmailResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAccountByEmailResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAccountByEmailResponseMultiError, or nil if none found.
func (m *GetAccountByEmailResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAccountByEmailResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAccount()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetAccountByEmailResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetAccountByEmailResponseValidationError{
					field:  "Account",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccount()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetAccountByEmailResponseValidationError{
				field:  "Account",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetAccountByEmailResponseMultiError(errors)
	}

	return nil
}

// GetAccountByEmailResponseMultiError is an error wrapping multiple validation
// errors returned by GetAccountByEmailResponse.ValidateAll() if the
// designated constraints aren't met.
type GetAccountByEmailResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAccountByEmailResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAccountByEmailResponseMultiError) AllErrors() []error { return m }

// GetAccountByEmailResponseValidationError is the validation error returned by
// GetAccountByEmailResponse.Validate if the designated constraints aren't met.
type GetAccountByEmailResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAccountByEmailResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAccountByEmailResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAccountByEmailResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAccountByEmailResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAccountByEmailResponseValidationError) ErrorName() string {
	return "GetAccountByEmailResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetAccountByEmailResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAccountByEmailResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAccountByEmailResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAccountByEmailResponseValidationError{}

// Validate checks the field values on DeleteAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteAccountRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteAccountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteAccountRequestMultiError, or nil if none found.
func (m *DeleteAccountRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteAccountRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetAccountId()); err != nil {
		err = DeleteAccountRequestValidationError{
			field:  "AccountId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteAccountRequestMultiError(errors)
	}

	return nil
}

func (m *DeleteAccountRequest) _validateUuid(uuid string) error {
	if matched := _service_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// DeleteAccountRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteAccountRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteAccountRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteAccountRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteAccountRequestMultiError) AllErrors() []error { return m }

// DeleteAccountRequestValidationError is the validation error returned by
// DeleteAccountRequest.Validate if the designated constraints aren't met.
type DeleteAccountRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteAccountRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteAccountRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteAccountRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteAccountRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteAccountRequestValidationError) ErrorName() string {
	return "DeleteAccountRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteAccountRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteAccountRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteAccountRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteAccountRequestValidationError{}