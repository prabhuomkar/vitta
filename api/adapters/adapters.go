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
)

type Adapter interface {
	GetTransactions(rows [][]string) []AdapterTransaction
}

func New(name string, _ *Config) Adapter { //nolint: ireturn
	if name == "icici" {
		return newICICI()
	}

	return nil // TODO(omkar): Add generic adapter
}
