package adapters

import (
	"strings"
	"time"
)

type (
	// Config model.
	Config struct {
		DateColumn, DateFormat, RemarksColumn, CreditColumn, DebitColumn string
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

func findIndices(cfg *Config, rows [][]string) ([2]int, [2]int, [2]int, [2]int) {
	results := map[string][2]int{}

	for rowIdx, row := range rows {
		for colIdx, col := range row {
			switch strings.TrimSpace(col) {
			case cfg.DateColumn:
				results["date"] = [2]int{rowIdx, colIdx}
			case cfg.RemarksColumn:
				results["remarks"] = [2]int{rowIdx, colIdx}
			case cfg.CreditColumn:
				results["credit"] = [2]int{rowIdx, colIdx}
				if cfg.CreditColumn == cfg.DebitColumn {
					results["debit"] = [2]int{rowIdx, colIdx}
				}
			case cfg.DebitColumn:
				results["debit"] = [2]int{rowIdx, colIdx}
			}
		}

		if len(results) == 4 { //nolint: gomnd,mnd
			break
		}
	}

	return results["date"], results["remarks"], results["credit"], results["debit"]
}
