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

type SearchHit struct {
	RawDocument json.RawMessage `json:"document"`
	Document    Document        `json:"-"`
}

type Document interface {
	GetCreatedAt() time.Time
	GetMetaData() map[string]interface{}
}

func (s *SearchService) SearchDocument(params SearchParams, resource ResourceType) (*SearchResponse, *http.Response, error) {
	if params.Q == "" {
		return nil, nil, fmt.Errorf("search query is required")
	}

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

	// Process each hit to convert RawDocument to proper type
	for i := range searchResp.Hits {
		if err := unmarshalDocument(&searchResp.Hits[i], resource); err != nil {
			return nil, resp, fmt.Errorf("failed to unmarshal document: %w", err)
		}
	}

	return searchResp, resp, nil
}

func unmarshalDocument(hit *SearchHit, resource ResourceType) error {
	switch resource {
	case ResourceLedgers:
		var doc Ledger
		if err := json.Unmarshal(hit.RawDocument, &doc); err != nil {
			return err
		}
		hit.Document = &doc
	case ResourceBalances:
		var doc LedgerBalance
		if err := json.Unmarshal(hit.RawDocument, &doc); err != nil {
			return err
		}
		hit.Document = &doc
	case ResourceTransactions:
		var doc Transaction
		if err := json.Unmarshal(hit.RawDocument, &doc); err != nil {
			return err
		}
		hit.Document = &doc
	default:
		return fmt.Errorf("unsupported resource type: %s", resource)
	}
	return nil
}
