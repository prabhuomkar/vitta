package adapters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		rows     [][]string
		expected []AdapterTransaction
	}{
		{
			"error parsing transaction time", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "",
				Credit: "Deposit Amount (INR )", Debit: "Withdrawal Amount (INR )", TransactionDiff: nil,
			},
			[][]string{
				{"Transaction Date", "", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"invalid-date", "NA", "0.0", "0.0"},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing debit amount", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "",
				Credit: "Deposit Amount (INR )", Debit: "Withdrawal Amount (INR )", TransactionDiff: nil,
			},
			[][]string{
				{"Transaction Date", "", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "0.0", "NA"},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing credit amount", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "",
				Credit: "Deposit Amount (INR )", Debit: "Withdrawal Amount (INR )", TransactionDiff: nil,
			},
			[][]string{
				{"Transaction Date", "", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "NA", "0.0"},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing debit amount when same column", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "Details",
				Credit: "Amount (INR)", Debit: "Amount (INR)", TransactionDiff: []string{"cr.", "dr."},
			},
			[][]string{
				{"Transaction Date", "Details", "Amount (INR)"},
				{"18/10/2024", "NA", "No Dr."},
			},
			[]AdapterTransaction{},
		},
		{
			"error parsing credit amount when same column", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "Details",
				Credit: "Amount (INR)", Debit: "Amount (INR)", TransactionDiff: []string{"cr.", "dr."},
			},
			[][]string{
				{"Transaction Date", "Details", "Amount (INR)"},
				{"18/10/2024", "NA", "No Cr."},
			},
			[]AdapterTransaction{},
		},
		{
			"success parsing debit amount", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "",
				Credit: "Deposit Amount (INR )", Debit: "Withdrawal Amount (INR )", TransactionDiff: nil,
			},
			[][]string{
				{"Transaction Date", "", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "0.0", "4.20"},
			},
			[]AdapterTransaction{
				{time.Date(2024, time.October, 18, 0, 0, 0, 0, time.UTC), "NA", 0.0, 4.2},
			},
		},
		{
			"success parsing credit amount", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "",
				Credit: "Deposit Amount (INR )", Debit: "Withdrawal Amount (INR )", TransactionDiff: nil,
			},
			[][]string{
				{"Transaction Date", "", "Deposit Amount (INR )", "Withdrawal Amount (INR )"},
				{"18/10/2024", "NA", "4.20", "0.0"},
			},
			[]AdapterTransaction{
				{time.Date(2024, time.October, 18, 0, 0, 0, 0, time.UTC), "NA", 4.2, 0.0},
			},
		},
		{
			"success parsing debit amount when same column", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "Details",
				Credit: "Amount (INR)", Debit: "Amount (INR)", TransactionDiff: []string{"cr.", "dr."},
			},
			[][]string{
				{"Transaction Date", "Details", "Amount (INR)"},
				{"18/10/2024", "NA", "4.20 Dr."},
			},
			[]AdapterTransaction{
				{time.Date(2024, time.October, 18, 0, 0, 0, 0, time.UTC), "NA", 0.0, 4.2},
			},
		},
		{
			"success parsing credit amount when same column", Config{DateName: "Transaction Date",
				DateFormats: []string{"02/01/2006"}, Remarks: "Details",
				Credit: "Amount (INR)", Debit: "Amount (INR)", TransactionDiff: []string{"cr.", "dr."},
			},
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
			transactions := GetTransactions(tc.config, tc.rows)
			assert.Equal(t, tc.expected, transactions)
		})
	}
}
