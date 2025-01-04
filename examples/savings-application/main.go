package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	blnkgo "github.com/anhbaysgalan1/blnk-go"
)

func main() {
	baseURL, _ := url.Parse("http://localhost:5001/")
	client := blnkgo.NewClient(baseURL, nil, blnkgo.WithTimeout(
		5*time.Second,
	), blnkgo.WithRetry(2))

	savingsLedgerBody := blnkgo.CreateLedgerRequest{
		Name: "Savings",
	}

	savingsLedger, resp, err := client.Ledger.Create(savingsLedgerBody)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(savingsLedger)

	savingsBody := blnkgo.CreateLedgerBalanceRequest{
		LedgerID: savingsLedger.LedgerID,
		Currency: "USD",
	}

	savingsBalance, resp, err := client.LedgerBalance.Create(savingsBody)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(savingsBalance)

	transactionBody := blnkgo.CreateTransactionRequest{
		ParentTransaction: blnkgo.ParentTransaction{
			Amount:      1000,
			Reference:   "ref-04",
			Precision:   100,
			Currency:    "USD",
			Source:      "@World",
			Destination: savingsBalance.BalanceID,
			Description: "Savings",
		},
		AllowOverdraft: true,
	}

	transaction, resp, err := client.Transaction.Create(transactionBody)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(transaction)
}
