// Package apperror is used for storing generic application errors so other services can reuse them.
package apperror

import (
	"fmt"
)

const format = "%s: %s, %w"

// New creates new wrapped erros in format: "{msg}: {original_error}, {error_for_client}".
func New(msg string, err, clientErr error) error {
	return fmt.Errorf(format, msg, err.Error(), clientErr)
}
