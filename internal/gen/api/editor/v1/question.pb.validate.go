// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: editor/v1/question.proto

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

// Validate checks the field values on Answer with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Answer) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Answer with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in AnswerMultiError, or nil if none found.
func (m *Answer) ValidateAll() error {
	return m.validate(true)
}

func (m *Answer) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Text

	// no validation rules for MediaUrl

	if len(errors) > 0 {
		return AnswerMultiError(errors)
	}

	return nil
}

// AnswerMultiError is an error wrapping multiple validation errors returned by
// Answer.ValidateAll() if the designated constraints aren't met.
type AnswerMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AnswerMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AnswerMultiError) AllErrors() []error { return m }

// AnswerValidationError is the validation error returned by Answer.Validate if
// the designated constraints aren't met.
type AnswerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AnswerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AnswerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AnswerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AnswerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AnswerValidationError) ErrorName() string { return "AnswerValidationError" }

// Error satisfies the builtin error interface
func (e AnswerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAnswer.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AnswerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AnswerValidationError{}

// Validate checks the field values on Question with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Question) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Question with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in QuestionMultiError, or nil
// if none found.
func (m *Question) ValidateAll() error {
	return m.validate(true)
}

func (m *Question) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Text

	if all {
		switch v := interface{}(m.GetAnswer()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, QuestionValidationError{
					field:  "Answer",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, QuestionValidationError{
					field:  "Answer",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAnswer()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return QuestionValidationError{
				field:  "Answer",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Author

	// no validation rules for MediaUrl

	if all {
		switch v := interface{}(m.GetCreateTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, QuestionValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, QuestionValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return QuestionValidationError{
				field:  "CreateTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return QuestionMultiError(errors)
	}

	return nil
}

// QuestionMultiError is an error wrapping multiple validation errors returned
// by Question.ValidateAll() if the designated constraints aren't met.
type QuestionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QuestionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QuestionMultiError) AllErrors() []error { return m }

// QuestionValidationError is the validation error returned by
// Question.Validate if the designated constraints aren't met.
type QuestionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QuestionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QuestionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QuestionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QuestionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QuestionValidationError) ErrorName() string { return "QuestionValidationError" }

// Error satisfies the builtin error interface
func (e QuestionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQuestion.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QuestionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QuestionValidationError{}

// Validate checks the field values on CreateQuestionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateQuestionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateQuestionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateQuestionRequestMultiError, or nil if none found.
func (m *CreateQuestionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateQuestionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetQuestion()); l < 3 || l > 200 {
		err := CreateQuestionRequestValidationError{
			field:  "Question",
			reason: "value length must be between 3 and 200 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetQuestionMediaUrl() != "" {

		if uri, err := url.Parse(m.GetQuestionMediaUrl()); err != nil {
			err = CreateQuestionRequestValidationError{
				field:  "QuestionMediaUrl",
				reason: "value must be a valid URI",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else if !uri.IsAbs() {
			err := CreateQuestionRequestValidationError{
				field:  "QuestionMediaUrl",
				reason: "value must be absolute",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if l := utf8.RuneCountInString(m.GetAnswer()); l < 3 || l > 100 {
		err := CreateQuestionRequestValidationError{
			field:  "Answer",
			reason: "value length must be between 3 and 100 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetAnswerMediaUrl() != "" {

		if uri, err := url.Parse(m.GetAnswerMediaUrl()); err != nil {
			err = CreateQuestionRequestValidationError{
				field:  "AnswerMediaUrl",
				reason: "value must be a valid URI",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else if !uri.IsAbs() {
			err := CreateQuestionRequestValidationError{
				field:  "AnswerMediaUrl",
				reason: "value must be absolute",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return CreateQuestionRequestMultiError(errors)
	}

	return nil
}

// CreateQuestionRequestMultiError is an error wrapping multiple validation
// errors returned by CreateQuestionRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateQuestionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateQuestionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateQuestionRequestMultiError) AllErrors() []error { return m }

// CreateQuestionRequestValidationError is the validation error returned by
// CreateQuestionRequest.Validate if the designated constraints aren't met.
type CreateQuestionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateQuestionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateQuestionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateQuestionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateQuestionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateQuestionRequestValidationError) ErrorName() string {
	return "CreateQuestionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateQuestionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateQuestionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateQuestionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateQuestionRequestValidationError{}

// Validate checks the field values on CreateQuestionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateQuestionResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateQuestionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateQuestionResponseMultiError, or nil if none found.
func (m *CreateQuestionResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateQuestionResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for QuestionId

	if len(errors) > 0 {
		return CreateQuestionResponseMultiError(errors)
	}

	return nil
}

// CreateQuestionResponseMultiError is an error wrapping multiple validation
// errors returned by CreateQuestionResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateQuestionResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateQuestionResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateQuestionResponseMultiError) AllErrors() []error { return m }

// CreateQuestionResponseValidationError is the validation error returned by
// CreateQuestionResponse.Validate if the designated constraints aren't met.
type CreateQuestionResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateQuestionResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateQuestionResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateQuestionResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateQuestionResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateQuestionResponseValidationError) ErrorName() string {
	return "CreateQuestionResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateQuestionResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateQuestionResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateQuestionResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateQuestionResponseValidationError{}

// Validate checks the field values on GetQuestionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetQuestionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetQuestionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetQuestionRequestMultiError, or nil if none found.
func (m *GetQuestionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetQuestionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for QuestionId

	if len(errors) > 0 {
		return GetQuestionRequestMultiError(errors)
	}

	return nil
}

// GetQuestionRequestMultiError is an error wrapping multiple validation errors
// returned by GetQuestionRequest.ValidateAll() if the designated constraints
// aren't met.
type GetQuestionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetQuestionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetQuestionRequestMultiError) AllErrors() []error { return m }

// GetQuestionRequestValidationError is the validation error returned by
// GetQuestionRequest.Validate if the designated constraints aren't met.
type GetQuestionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetQuestionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetQuestionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetQuestionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetQuestionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetQuestionRequestValidationError) ErrorName() string {
	return "GetQuestionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetQuestionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetQuestionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetQuestionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetQuestionRequestValidationError{}

// Validate checks the field values on GetQuestionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetQuestionResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetQuestionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetQuestionResponseMultiError, or nil if none found.
func (m *GetQuestionResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetQuestionResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetQuestion()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetQuestionResponseValidationError{
					field:  "Question",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetQuestionResponseValidationError{
					field:  "Question",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetQuestion()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetQuestionResponseValidationError{
				field:  "Question",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetQuestionResponseMultiError(errors)
	}

	return nil
}

// GetQuestionResponseMultiError is an error wrapping multiple validation
// errors returned by GetQuestionResponse.ValidateAll() if the designated
// constraints aren't met.
type GetQuestionResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetQuestionResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetQuestionResponseMultiError) AllErrors() []error { return m }

// GetQuestionResponseValidationError is the validation error returned by
// GetQuestionResponse.Validate if the designated constraints aren't met.
type GetQuestionResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetQuestionResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetQuestionResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetQuestionResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetQuestionResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetQuestionResponseValidationError) ErrorName() string {
	return "GetQuestionResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetQuestionResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetQuestionResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetQuestionResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetQuestionResponseValidationError{}
