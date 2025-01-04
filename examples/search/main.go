package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	blnkgo "github.com/anhbaysgalan1/blnk-go"
)

func main() {
	baseURL, _ := url.Parse("http://localhost:5002/")
	client := blnkgo.NewClient(baseURL, nil, blnkgo.WithTimeout(
		10*time.Second,
	), blnkgo.WithRetry(3))

	searchSvc := blnkgo.NewSearchService(client)

	// Search for ledgers with pagination
	ledgerParams := blnkgo.SearchParams{
		Q:        "*",
		FilterBy: stringPtr("name:fanhubble-main-ledger"),
		SortBy:   stringPtr("created_at:desc"),
		Page:     intPtr(1),
		PerPage:  intPtr(10),
	}

	ledgerResults, resp, err := searchSvc.SearchDocument(ledgerParams, blnkgo.ResourceLedgers)
	if err != nil {
		log.Printf("Failed to search ledgers: %v", err)
		if resp != nil {
			log.Printf("Response status: %d", resp.StatusCode)
		}
		return
	}

	fmt.Printf("Found %d ledgers (Page %d/%d)\n",
		ledgerResults.Found,
		ledgerResults.Page,
		(ledgerResults.OutOf+9)/10)

	for _, hit := range ledgerResults.Hits {
		if ledger, ok := hit.Document.(*blnkgo.Ledger); ok {
			fmt.Printf("Ledger: %s (ID: %s)\n", ledger.Name, ledger.LedgerID)
			if metadata := ledger.GetMetaData(); metadata != nil {
				fmt.Printf("  Metadata: %v\n", metadata)
			}
		}
	}

	// Search for high-value balances
	balanceParams := blnkgo.SearchParams{
		Q:        "*",
		FilterBy: stringPtr("balance:>1"),
		SortBy:   stringPtr("balance:desc,created_at:desc"),
		Page:     intPtr(1),
		PerPage:  intPtr(10),
	}

	balanceResults, resp, err := searchSvc.SearchDocument(balanceParams, blnkgo.ResourceBalances)
	if err != nil {
		log.Printf("Failed to search balances: %v", err)
		if resp != nil {
			log.Printf("Response status: %d", resp.StatusCode)
		}
		return
	}

	fmt.Printf("\nFound %d high-value balances (Page %d/%d)\n",
		balanceResults.Found,
		balanceResults.Page,
		(balanceResults.OutOf+9)/10)

	for _, hit := range balanceResults.Hits {
		if balance, ok := hit.Document.(*blnkgo.LedgerBalance); ok {
			fmt.Printf("Balance: %.2f %s (ID: %s)\n",
				float64(balance.Balance)/100.0,
				balance.Currency,
				balance.BalanceID)
			if metadata := balance.GetMetaData(); metadata != nil {
				fmt.Printf("  Metadata: %v\n", metadata)
			}
		}
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
