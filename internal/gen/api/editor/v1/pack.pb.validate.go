// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: editor/v1/pack.proto

package editorv1

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

// Validate checks the field values on PackRound with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PackRound) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PackRound with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PackRoundMultiError, or nil
// if none found.
func (m *PackRound) ValidateAll() error {
	return m.validate(true)
}

func (m *PackRound) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Name

	if len(errors) > 0 {
		return PackRoundMultiError(errors)
	}

	return nil
}

// PackRoundMultiError is an error wrapping multiple validation errors returned
// by PackRound.ValidateAll() if the designated constraints aren't met.
type PackRoundMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PackRoundMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PackRoundMultiError) AllErrors() []error { return m }

// PackRoundValidationError is the validation error returned by
// PackRound.Validate if the designated constraints aren't met.
type PackRoundValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PackRoundValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PackRoundValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PackRoundValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PackRoundValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PackRoundValidationError) ErrorName() string { return "PackRoundValidationError" }

// Error satisfies the builtin error interface
func (e PackRoundValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPackRound.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PackRoundValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PackRoundValidationError{}

// Validate checks the field values on Pack with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Pack) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Pack with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PackMultiError, or nil if none found.
func (m *Pack) ValidateAll() error {
	return m.validate(true)
}

func (m *Pack) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Name

	// no validation rules for Author

	// no validation rules for IsPublished

	// no validation rules for CoverUrl

	for idx, item := range m.GetRounds() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PackValidationError{
						field:  fmt.Sprintf("Rounds[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PackValidationError{
						field:  fmt.Sprintf("Rounds[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PackValidationError{
					field:  fmt.Sprintf("Rounds[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetCreateTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PackValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PackValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PackValidationError{
				field:  "CreateTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetPublishTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PackValidationError{
					field:  "PublishTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PackValidationError{
					field:  "PublishTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPublishTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PackValidationError{
				field:  "PublishTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PackMultiError(errors)
	}

	return nil
}

// PackMultiError is an error wrapping multiple validation errors returned by
// Pack.ValidateAll() if the designated constraints aren't met.
type PackMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PackMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PackMultiError) AllErrors() []error { return m }

// PackValidationError is the validation error returned by Pack.Validate if the
// designated constraints aren't met.
type PackValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PackValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PackValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PackValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PackValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PackValidationError) ErrorName() string { return "PackValidationError" }

// Error satisfies the builtin error interface
func (e PackValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPack.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PackValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PackValidationError{}

// Validate checks the field values on PackStats with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PackStats) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PackStats with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PackStatsMultiError, or nil
// if none found.
func (m *PackStats) ValidateAll() error {
	return m.validate(true)
}

func (m *PackStats) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for RoundCount

	// no validation rules for TopicCount

	// no validation rules for QuestionCount

	// no validation rules for VideoCount

	// no validation rules for AudioCount

	// no validation rules for ImageCount

	if len(errors) > 0 {
		return PackStatsMultiError(errors)
	}

	return nil
}

// PackStatsMultiError is an error wrapping multiple validation errors returned
// by PackStats.ValidateAll() if the designated constraints aren't met.
type PackStatsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PackStatsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PackStatsMultiError) AllErrors() []error { return m }

// PackStatsValidationError is the validation error returned by
// PackStats.Validate if the designated constraints aren't met.
type PackStatsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PackStatsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PackStatsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PackStatsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PackStatsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PackStatsValidationError) ErrorName() string { return "PackStatsValidationError" }

// Error satisfies the builtin error interface
func (e PackStatsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPackStats.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PackStatsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PackStatsValidationError{}

// Validate checks the field values on PackWithStats with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PackWithStats) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PackWithStats with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PackWithStatsMultiError, or
// nil if none found.
func (m *PackWithStats) ValidateAll() error {
	return m.validate(true)
}

func (m *PackWithStats) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPack()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PackWithStatsValidationError{
					field:  "Pack",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PackWithStatsValidationError{
					field:  "Pack",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPack()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PackWithStatsValidationError{
				field:  "Pack",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetStats()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PackWithStatsValidationError{
					field:  "Stats",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PackWithStatsValidationError{
					field:  "Stats",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetStats()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PackWithStatsValidationError{
				field:  "Stats",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PackWithStatsMultiError(errors)
	}

	return nil
}

// PackWithStatsMultiError is an error wrapping multiple validation errors
// returned by PackWithStats.ValidateAll() if the designated constraints
// aren't met.
type PackWithStatsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PackWithStatsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PackWithStatsMultiError) AllErrors() []error { return m }

// PackWithStatsValidationError is the validation error returned by
// PackWithStats.Validate if the designated constraints aren't met.
type PackWithStatsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PackWithStatsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PackWithStatsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PackWithStatsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PackWithStatsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PackWithStatsValidationError) ErrorName() string { return "PackWithStatsValidationError" }

// Error satisfies the builtin error interface
func (e PackWithStatsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPackWithStats.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PackWithStatsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PackWithStatsValidationError{}

// Validate checks the field values on GetPackRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetPackRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetPackRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetPackRequestMultiError,
// or nil if none found.
func (m *GetPackRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetPackRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PackId

	if len(errors) > 0 {
		return GetPackRequestMultiError(errors)
	}

	return nil
}

// GetPackRequestMultiError is an error wrapping multiple validation errors
// returned by GetPackRequest.ValidateAll() if the designated constraints
// aren't met.
type GetPackRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetPackRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetPackRequestMultiError) AllErrors() []error { return m }

// GetPackRequestValidationError is the validation error returned by
// GetPackRequest.Validate if the designated constraints aren't met.
type GetPackRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetPackRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetPackRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetPackRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetPackRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetPackRequestValidationError) ErrorName() string { return "GetPackRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetPackRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetPackRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetPackRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetPackRequestValidationError{}

// Validate checks the field values on GetPackResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetPackResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetPackResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetPackResponseMultiError, or nil if none found.
func (m *GetPackResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetPackResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPack()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetPackResponseValidationError{
					field:  "Pack",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetPackResponseValidationError{
					field:  "Pack",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPack()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetPackResponseValidationError{
				field:  "Pack",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetPackResponseMultiError(errors)
	}

	return nil
}

// GetPackResponseMultiError is an error wrapping multiple validation errors
// returned by GetPackResponse.ValidateAll() if the designated constraints
// aren't met.
type GetPackResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetPackResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetPackResponseMultiError) AllErrors() []error { return m }

// GetPackResponseValidationError is the validation error returned by
// GetPackResponse.Validate if the designated constraints aren't met.
type GetPackResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetPackResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetPackResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetPackResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetPackResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetPackResponseValidationError) ErrorName() string { return "GetPackResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetPackResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetPackResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetPackResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetPackResponseValidationError{}

// Validate checks the field values on CreatePackRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CreatePackRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePackRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePackRequestMultiError, or nil if none found.
func (m *CreatePackRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePackRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetPackName()); l < 3 || l > 50 {
		err := CreatePackRequestValidationError{
			field:  "PackName",
			reason: "value length must be between 3 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCoverUrl() != "" {

		if uri, err := url.Parse(m.GetCoverUrl()); err != nil {
			err = CreatePackRequestValidationError{
				field:  "CoverUrl",
				reason: "value must be a valid URI",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else if !uri.IsAbs() {
			err := CreatePackRequestValidationError{
				field:  "CoverUrl",
				reason: "value must be absolute",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(m.GetTags()) > 5 {
		err := CreatePackRequestValidationError{
			field:  "Tags",
			reason: "value must contain no more than 5 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	_CreatePackRequest_Tags_Unique := make(map[string]struct{}, len(m.GetTags()))

	for idx, item := range m.GetTags() {
		_, _ = idx, item

		if _, exists := _CreatePackRequest_Tags_Unique[item]; exists {
			err := CreatePackRequestValidationError{
				field:  fmt.Sprintf("Tags[%v]", idx),
				reason: "repeated value must contain unique items",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {
			_CreatePackRequest_Tags_Unique[item] = struct{}{}
		}

		// no validation rules for Tags[idx]
	}

	if len(errors) > 0 {
		return CreatePackRequestMultiError(errors)
	}

	return nil
}

// CreatePackRequestMultiError is an error wrapping multiple validation errors
// returned by CreatePackRequest.ValidateAll() if the designated constraints
// aren't met.
type CreatePackRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePackRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePackRequestMultiError) AllErrors() []error { return m }

// CreatePackRequestValidationError is the validation error returned by
// CreatePackRequest.Validate if the designated constraints aren't met.
type CreatePackRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePackRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePackRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePackRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePackRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePackRequestValidationError) ErrorName() string {
	return "CreatePackRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePackRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePackRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePackRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePackRequestValidationError{}

// Validate checks the field values on CreatePackResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreatePackResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePackResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePackResponseMultiError, or nil if none found.
func (m *CreatePackResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePackResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PackId

	if len(errors) > 0 {
		return CreatePackResponseMultiError(errors)
	}

	return nil
}

// CreatePackResponseMultiError is an error wrapping multiple validation errors
// returned by CreatePackResponse.ValidateAll() if the designated constraints
// aren't met.
type CreatePackResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePackResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePackResponseMultiError) AllErrors() []error { return m }

// CreatePackResponseValidationError is the validation error returned by
// CreatePackResponse.Validate if the designated constraints aren't met.
type CreatePackResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePackResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePackResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePackResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePackResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePackResponseValidationError) ErrorName() string {
	return "CreatePackResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePackResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePackResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePackResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePackResponseValidationError{}

// Validate checks the field values on PublishPackRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PublishPackRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PublishPackRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PublishPackRequestMultiError, or nil if none found.
func (m *PublishPackRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PublishPackRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PackageId

	if len(errors) > 0 {
		return PublishPackRequestMultiError(errors)
	}

	return nil
}

// PublishPackRequestMultiError is an error wrapping multiple validation errors
// returned by PublishPackRequest.ValidateAll() if the designated constraints
// aren't met.
type PublishPackRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PublishPackRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PublishPackRequestMultiError) AllErrors() []error { return m }

// PublishPackRequestValidationError is the validation error returned by
// PublishPackRequest.Validate if the designated constraints aren't met.
type PublishPackRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PublishPackRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PublishPackRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PublishPackRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PublishPackRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PublishPackRequestValidationError) ErrorName() string {
	return "PublishPackRequestValidationError"
}

// Error satisfies the builtin error interface
func (e PublishPackRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPublishPackRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PublishPackRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PublishPackRequestValidationError{}

// Validate checks the field values on PublishPackResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PublishPackResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PublishPackResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PublishPackResponseMultiError, or nil if none found.
func (m *PublishPackResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *PublishPackResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPack()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PublishPackResponseValidationError{
					field:  "Pack",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PublishPackResponseValidationError{
					field:  "Pack",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPack()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PublishPackResponseValidationError{
				field:  "Pack",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PublishPackResponseMultiError(errors)
	}

	return nil
}

// PublishPackResponseMultiError is an error wrapping multiple validation
// errors returned by PublishPackResponse.ValidateAll() if the designated
// constraints aren't met.
type PublishPackResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PublishPackResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PublishPackResponseMultiError) AllErrors() []error { return m }

// PublishPackResponseValidationError is the validation error returned by
// PublishPackResponse.Validate if the designated constraints aren't met.
type PublishPackResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PublishPackResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PublishPackResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PublishPackResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PublishPackResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PublishPackResponseValidationError) ErrorName() string {
	return "PublishPackResponseValidationError"
}

// Error satisfies the builtin error interface
func (e PublishPackResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPublishPackResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PublishPackResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PublishPackResponseValidationError{}
