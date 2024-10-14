package adapters

import (
	"fmt"
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

	return &ICICI{cfg: cfg}
}

func (i *ICICI) GetTransactions(rows [][]string) []AdapterTransaction {
	dateIdx, remarkIdx, creditIdx, debitIdx := findIndices(i.cfg, rows)
	slog.Info("indices", "date", dateIdx, "remark", remarkIdx, "credit", creditIdx, "debit", debitIdx)

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
			amountStr := rows[idx][debitIdx[1]]
			if strings.HasSuffix(strings.ToLower(amountStr), "dr.") {
				debit, err = parseAmount(amountStr)
				if err != nil {
					slog.Error("error parsing debit amount", "error", err)

					break
				}
			} else {
				credit, err = parseAmount(amountStr)
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

func parseAmount(value string) (float64, error) {
	cleanedValue := strings.Split(strings.ReplaceAll(value, ",", ""), " ")[0]

	res, err := strconv.ParseFloat(cleanedValue, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing amount: %w", err)
	}

	return res, nil
}
