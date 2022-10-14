package cryptostr

type charset string

const (
	// Alphabet contains all lowercase characaters from english alphabet.
	AlphabetLower = charset("abcdefghijklmnopqrstuvwxyz")

	// Alphabet contains all uppercase characaters from english alphabet.
	AlphabetUpper = charset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Alphabet contains lowercase and uppercase characters from english alphabet.
	Alphabet = charset(AlphabetLower + AlphabetUpper)

	// Digits contains all digits.
	Digits = charset("0123456789")

	// AlphabetDigits contains uppercase, lowercase characters from english alphabet and digits.
	AlphabetDigits = charset(Alphabet + Digits)

	// Specials contains all special characters.
	Specials = charset("~=+-_%^&*/()[]{}<>.,:;/!@#$?|\"'")

	// AllChars contains uppercase, lowercase characters from english alphabet, digits and special characters.
	AllChars = charset(AlphabetDigits + Digits + Specials)
)
