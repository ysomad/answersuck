package email

type (
	withURL struct {
		URL string
	}

	sendEmailDTO struct {
		to         string
		template   string
		subject    string
		format     string
		formatArgs []any
	}
)
