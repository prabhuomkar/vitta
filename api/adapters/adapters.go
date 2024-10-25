package adapters

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// Config for the adapter.
	Config struct {
		DateName        string
		DateFormats     []string
		Remarks         string
		Credit          string
		Debit           string
		TransactionDiff []string
	}

	// AdapterTransaction model.
	AdapterTransaction struct {
		Date    time.Time
		Remarks string
		Credit  float64
		Debit   float64
	}
)

func New(filePath string) (map[string]Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("error opening file", "error", err)

		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		slog.Error("error reading rows", "error", err)

		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	adapters := map[string]Config{}

	for _, row := range rows {
		adapters[row[0]+"-"+row[1]] = Config{
			DateName:        row[2],
			DateFormats:     strings.Split(row[3], ";"),
			Remarks:         row[4],
			Credit:          row[5],
			Debit:           row[6],
			TransactionDiff: strings.Split(row[7], ";"),
		}
	}

	return adapters, nil
}

func findDataIndices(cfg Config, rows [][]string) ([2]int, [2]int, [2]int, [2]int) {
	results := map[string][2]int{}

	for rowIdx, row := range rows {
		for colIdx, col := range row {
			switch strings.TrimSpace(col) {
			case cfg.DateName:
				results["date"] = [2]int{rowIdx, colIdx}
			case cfg.Remarks:
				results["remarks"] = [2]int{rowIdx, colIdx}
			case cfg.Credit:
				results["credit"] = [2]int{rowIdx, colIdx}
				if cfg.Credit == cfg.Debit {
					results["debit"] = [2]int{rowIdx, colIdx}
				}
			case cfg.Debit:
				results["debit"] = [2]int{rowIdx, colIdx}
			}
		}

		if len(results) == 4 { //nolint: gomnd,mnd
			break
		}
	}

	return results["date"], results["remarks"], results["credit"], results["debit"]
}

func parseTransactionAmount(value string) (float64, error) {
	cleanedValue := strings.Split(strings.ReplaceAll(value, ",", ""), " ")[0]

	res, err := strconv.ParseFloat(cleanedValue, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing transaction amount: %w", err)
	}

	return res, nil
}

func parseTransactionDateTime(formats []string, value string) (time.Time, error) {
	var (
		date time.Time
		err  error
	)

	for _, format := range formats {
		date, err = time.Parse(format, value)
		if err == nil {
			return date, nil
		}
	}

	return time.Now(), err
}

func GetTransactions(cfg Config, rows [][]string) []AdapterTransaction {
	dateIdx, remarkIdx, creditIdx, debitIdx := findDataIndices(cfg, rows)

	transactions := []AdapterTransaction{}
	idx := dateIdx[0] + 1

	for idx < len(rows) {
		date, err := parseTransactionDateTime(cfg.DateFormats, rows[idx][dateIdx[1]])
		if err != nil {
			slog.Error("error parsing transaction date and time", "error", err)

			break
		}

		debit, credit := 0.0, 0.0

		if debitIdx[1] == creditIdx[1] { //nolint: nestif
			amountStr := rows[idx][debitIdx[1]]
			if strings.Contains(strings.ToLower(amountStr), cfg.TransactionDiff[1]) {
				debit, err = parseTransactionAmount(amountStr)
				if err != nil {
					slog.Error("error parsing debit amount", "error", err)

					break
				}
			} else {
				credit, err = parseTransactionAmount(amountStr)
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
