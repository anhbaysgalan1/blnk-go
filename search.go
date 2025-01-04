package blnkgo

import (
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
	Document Document `json:"document"`
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

func (s *SearchService) SearchDocument(params SearchParams, resource ResourceType) (*SearchResponse, *http.Response, error) {
	// Validate search query
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

	return searchResp, resp, nil
}
