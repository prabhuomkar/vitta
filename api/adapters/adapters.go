package adapters

import "time"

type (
	// Config model.
	Config struct {
		DateColumn    string
		DateFormat    string
		RemarksColumn string
		CreditColumn  string
		DebitColumn   string
	}

	// AdapterTransaction model.
	AdapterTransaction struct {
		Date    time.Time
		Remarks string
		Credit  float64
		Debit   float64
	}

	Adapter interface {
		GetTransactions(rows [][]string) []AdapterTransaction
	}
)

const (
	AdapterICICI = "icici"
)

func New(name, category string, _ *Config) Adapter {
	if name == AdapterICICI {
		return newICICI(category)
	}

	return nil // TODO(omkar): Add generic adapter
}
