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
		5*time.Second,
	), blnkgo.WithRetry(2))

	searchSvc := blnkgo.NewSearchService(client)

	// Example 1: Search for ledgers
	ledgerParams := blnkgo.SearchParams{
		Q:        "*",
		FilterBy: stringPtr("name:World"),
		SortBy:   stringPtr("created_at:desc"),
	}

	ledgerResults, _, err := searchSvc.SearchDocument(ledgerParams, blnkgo.ResourceLedgers)
	if err != nil {
		log.Fatalf("Failed to search ledgers: %v", err)
	}
	fmt.Printf("Found %d ledgers\n", ledgerResults.Found)

	// Example 2: Search for balances
	balanceParams := blnkgo.SearchParams{
		Q:        "*",
		FilterBy: stringPtr("balance:>1000"),
		SortBy:   stringPtr("created_at:desc"),
	}

	balanceResults, _, err := searchSvc.SearchDocument(balanceParams, blnkgo.ResourceBalances)
	if err != nil {
		log.Fatalf("Failed to search balances: %v", err)
	}
	fmt.Printf("Found %d balances\n", balanceResults.Found)

	// Example 3: Search for transactions
	txParams := blnkgo.SearchParams{
		Q:        "*",
		FilterBy: stringPtr("amount:[2000..100000]"),
		SortBy:   stringPtr("created_at:desc"),
	}

	txResults, _, err := searchSvc.SearchDocument(txParams, blnkgo.ResourceTransactions)
	if err != nil {
		log.Fatalf("Failed to search transactions: %v", err)
	}
	fmt.Printf("Found %d transactions\n", txResults.Found)
}

func stringPtr(s string) *string {
	return &s
}
