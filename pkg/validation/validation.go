package validation

type ErrorTranslator interface {
	// Translate translates received error into map of errors
	TranslateError(err error) map[string]string
}
