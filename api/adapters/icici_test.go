package adapters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestICICIGetTransactions(t *testing.T) {
	tests := []struct {
		name     string
		category string
		rows     [][]string
		expected []AdapterTransaction
	}{
		{
			"error parsing transaction time", "CASH",
			[][]string{
				{"Transaction Date", "Transaction Remarks", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"invalid-date", "NA", "0.0", "0.0"},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing debit amount", "CASH",
			[][]string{
				{"Transaction Date", "Transaction Remarks", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "0.0", "NA"},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing credit amount", "CASH",
			[][]string{
				{"Transaction Date", "Transaction Remarks", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "NA", "0.0"},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing debit amount when same column", "CC",
			[][]string{
				{"Transaction Date", "Details", "Amount (INR)"},
				{"18/10/2024", "NA", "No Dr."},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing credit amount when same column", "CC",
			[][]string{
				{"Transaction Date", "Details", "Amount (INR)"},
				{"18/10/2024", "NA", "No Cr."},
			},
			[]AdapterTransaction{},
		},
		{
			"success parsing debit amount", "CASH",
			[][]string{
				{"Transaction Date", "Transaction Remarks", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "0.0", "4.20"},
			},
			[]AdapterTransaction{
				{time.Date(2024, time.October, 18, 0, 0, 0, 0, time.UTC), "NA", 0.0, 4.2},
			},
		},
		{
			"success parsing credit amount", "CASH",
			[][]string{
				{"Transaction Date", "Transaction Remarks", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "4.20", "0.0"},
			},
			[]AdapterTransaction{
				{time.Date(2024, time.October, 18, 0, 0, 0, 0, time.UTC), "NA", 4.2, 0.0},
			},
		},
		{
			"success parsing debit amount when same column", "CC",
			[][]string{
				{"Transaction Date", "Details", "Amount (INR)"},
				{"18/10/2024", "NA", "4.20 Dr."},
			},
			[]AdapterTransaction{
				{time.Date(2024, time.October, 18, 0, 0, 0, 0, time.UTC), "NA", 0.0, 4.2},
			},
		},
		{
			"success parsing credit amount when same column", "CC",
			[][]string{
				{"Transaction Date", "Details", "Amount (INR)"},
				{"18/10/2024", "NA", "4.20 Cr."},
			},
			[]AdapterTransaction{
				{time.Date(2024, time.October, 18, 0, 0, 0, 0, time.UTC), "NA", 4.2, 0.0},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			transactions := newICICI(tc.category).GetTransactions(tc.rows)
			assert.Equal(t, tc.expected, transactions)
		})
	}
}
