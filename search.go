package blnkgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SearchService struct {
	client ClientInterface
}

func NewSearchService(c ClientInterface) *SearchService {
	return &SearchService{client: c}
}

type SearchParams struct {
	Q        string  `json:"q"`
	QueryBy  *string `json:"query_by,omitempty"`
	FilterBy *string `json:"filter_by,omitempty"`
	SortBy   *string `json:"sort_by,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PerPage  *int    `json:"per_page,omitempty"`
}

type SearchResponse struct {
	Found        int         `json:"found"`
	OutOf        int         `json:"out_of"`
	Page         int         `json:"page"`
	SearchTimeMs int         `json:"search_time_ms"`
	Hits         []SearchHit `json:"hits"`
}

type Document interface {
	GetCreatedAt() time.Time
	GetMetaData() map[string]interface{}
}

func (l *Ledger) GetCreatedAt() time.Time             { return l.CreatedAt }
func (l *Ledger) GetMetaData() map[string]interface{} { return l.MetaData }

func (b *LedgerBalance) GetCreatedAt() time.Time             { return b.CreatedAt }
func (b *LedgerBalance) GetMetaData() map[string]interface{} { return b.MetaData }

func (t *Transaction) GetCreatedAt() time.Time             { return t.CreatedAt }
func (t *Transaction) GetMetaData() map[string]interface{} { return t.MetaData }

type SearchHit struct {
	Document json.RawMessage `json:"document"`
}

func (s *SearchService) SearchDocument(params SearchParams, resource ResourceType) (*SearchResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("search/%s", resource)
	req, err := s.client.NewRequest(endpoint, http.MethodPost, params)
	if err != nil {
		return nil, nil, err
	}

	searchResp := new(SearchResponse)
	resp, err := s.client.CallWithRetry(req, searchResp)
	if err != nil {
		return nil, resp, err
	}

	// Convert raw messages to proper types based on resource
	for i, hit := range searchResp.Hits {
		var doc Document
		switch resource {
		case ResourceLedgers:
			var ledger Ledger
			if err := json.Unmarshal(hit.Document, &ledger); err != nil {
				return nil, resp, err
			}
			doc = &ledger
		case ResourceBalances:
			var balance LedgerBalance
			if err := json.Unmarshal(hit.Document, &balance); err != nil {
				return nil, resp, err
			}
			doc = &balance
		case ResourceTransactions:
			var tx Transaction
			if err := json.Unmarshal(hit.Document, &tx); err != nil {
				return nil, resp, err
			}
			doc = &tx
		}
		searchResp.Hits[i].Document = doc
	}

	return searchResp, resp, nil
}
