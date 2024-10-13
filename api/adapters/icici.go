package adapters

import (
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type ICICI struct {
	cfg *Config
}

func newICICI(category string) *ICICI {
	cfg := &Config{
		DateColumn:    "Transaction Date",
		DateFormat:    "02/01/2006",
		RemarksColumn: "Transaction Remarks",
		CreditColumn:  "Deposit Amount (INR )",
		DebitColumn:   "Withdrawal Amount (INR )",
	}
	if category == "CC" {
		cfg.RemarksColumn = "Details"
		cfg.CreditColumn = "Amount (INR)"
		cfg.DebitColumn = "Amount (INR)"
	}

	return &ICICI{
		cfg: cfg,
	}
}

func (i *ICICI) GetTransactions(rows [][]string) []AdapterTransaction { //nolint: cyclop,funlen,gocognit
	dateIdx, remarkIdx, creditIdx, debitIdx := []int{-1, -1}, []int{-1, -1}, []int{-1, -1}, []int{-1, -1}

	for rowIdx, row := range rows {
		for colIdx, col := range row {
			if col == i.cfg.DateColumn {
				dateIdx = []int{rowIdx, colIdx}
			}

			if col == i.cfg.RemarksColumn {
				remarkIdx = []int{rowIdx, colIdx}
			}

			if col == i.cfg.CreditColumn {
				creditIdx = []int{rowIdx, colIdx}
			}

			if col == i.cfg.DebitColumn {
				debitIdx = []int{rowIdx, colIdx}
			}
		}
	}

	transactions := []AdapterTransaction{}

	idx := dateIdx[0] + 1

	for idx < len(rows) {
		date, err := time.Parse(i.cfg.DateFormat, rows[idx][dateIdx[1]])
		if err != nil {
			slog.Error("error parsing transaction time", "error", err)

			break
		}

		debit, credit := 0.0, 0.0

		if debitIdx[1] == creditIdx[1] { //nolint: nestif
			if strings.HasSuffix(strings.ToLower(rows[idx][debitIdx[1]]), "dr.") {
				debit, err = strconv.ParseFloat(strings.Split(strings.ReplaceAll(rows[idx][debitIdx[1]], ",", ""), " ")[0], 64)
				if err != nil {
					slog.Error("error parsing debit amount", "error", err)

					break
				}
			} else {
				credit, err = strconv.ParseFloat(strings.Split(strings.ReplaceAll(rows[idx][creditIdx[1]], ",", ""), " ")[0], 64)
				if err != nil {
					slog.Error("error parsing credit amount", "error", err)

					break
				}
			}
		} else {
			debit, err = strconv.ParseFloat(rows[idx][debitIdx[1]], 64)
			if err != nil {
				slog.Error("error parsing debit amount", "error", err)

				break
			}

			credit, err = strconv.ParseFloat(rows[idx][creditIdx[1]], 64)
			if err != nil {
				slog.Error("error parsing credit amount", "error", err)

				break
			}
		}

		transactions = append(transactions, AdapterTransaction{
			Date:    date,
			Remarks: rows[idx][remarkIdx[1]],
			Debit:   debit,
			Credit:  credit,
		})

		idx++
	}

	return transactions
}
