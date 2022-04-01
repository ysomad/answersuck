package validation

type ErrorTranslator interface {
	// TranslateError translates received error into map of errors
	TranslateError(err error) map[string]string
}
