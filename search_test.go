package blnkgo_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	blnkgo "github.com/anhbaysgalan1/blnk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func stringPtr(s string) *string {
	return &s
}

func setupSearchService() (*MockClient, *blnkgo.SearchService) {
	mockClient := new(MockClient)
	svc := blnkgo.NewSearchService(mockClient)
	return mockClient, svc
}

func TestSearchService_SearchDocument(t *testing.T) {
	tests := []struct {
		name        string
		params      blnkgo.SearchParams
		resource    blnkgo.ResourceType
		expectError bool
		errorMsg    string
		statusCode  int
		setupMocks  func(*MockClient)
	}{
		{
			name: "successful ledger search",
			params: blnkgo.SearchParams{
				Q:        "*",
				FilterBy: stringPtr("name:World"),
				SortBy:   stringPtr("created_at:desc"),
			},
			resource:    blnkgo.ResourceLedgers,
			expectError: false,
			statusCode:  http.StatusOK,
			setupMocks: func(m *MockClient) {
				fixedTime := time.Date(2024, time.February, 20, 5, 28, 3, 0, time.UTC)
				expectedResponse := &blnkgo.SearchResponse{
					Found:        1,
					OutOf:        1,
					Page:         1,
					SearchTimeMs: 1,
					Hits: []blnkgo.SearchHit{
						{
							Document: &blnkgo.Ledger{
								LedgerID:  "ldg_073f7ffe-9dfd-42ce-aa50-d1dca1788adc",
								Name:      "World Ledger",
								CreatedAt: fixedTime,
								MetaData: map[string]interface{}{
									"type": "main",
								},
							},
						},
					},
				}

				m.On("NewRequest", "search/ledgers", http.MethodPost, mock.Anything).
					Return(&http.Request{}, nil)
				m.On("CallWithRetry", mock.Anything, mock.Anything).
					Return(&http.Response{StatusCode: http.StatusOK}, nil).
					Run(func(args mock.Arguments) {
						resp := args.Get(1).(*blnkgo.SearchResponse)
						*resp = *expectedResponse
					})
			},
		},
		{
			name: "empty search query",
			params: blnkgo.SearchParams{
				Q: "",
			},
			resource:    blnkgo.ResourceLedgers,
			expectError: true,
			errorMsg:    "search query is required",
			setupMocks:  func(m *MockClient) {},
		},
		{
			name: "request creation failure",
			params: blnkgo.SearchParams{
				Q: "*",
			},
			resource:    blnkgo.ResourceLedgers,
			expectError: true,
			errorMsg:    "failed to create request",
			setupMocks: func(m *MockClient) {
				m.On("NewRequest", "search/ledgers", http.MethodPost, mock.Anything).
					Return(nil, errors.New("failed to create request"))
			},
		},
		{
			name: "server error",
			params: blnkgo.SearchParams{
				Q: "*",
			},
			resource:    blnkgo.ResourceLedgers,
			expectError: true,
			errorMsg:    "server error",
			statusCode:  http.StatusInternalServerError,
			setupMocks: func(m *MockClient) {
				m.On("NewRequest", "search/ledgers", http.MethodPost, mock.Anything).
					Return(&http.Request{}, nil)
				m.On("CallWithRetry", mock.Anything, mock.Anything).
					Return(&http.Response{StatusCode: http.StatusInternalServerError},
						errors.New("server error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient, svc := setupSearchService()
			tt.setupMocks(mockClient)

			searchResp, resp, err := svc.SearchDocument(tt.params, tt.resource)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				if tt.params.Q == "" {
					assert.Nil(t, resp)
					assert.Nil(t, searchResp)
					return
				}
				if tt.name == "request creation failure" {
					assert.Nil(t, searchResp)
					assert.Nil(t, resp)
					mockClient.AssertNotCalled(t, "CallWithRetry")
				} else {
					assert.Equal(t, tt.statusCode, resp.StatusCode)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, searchResp)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.statusCode, resp.StatusCode)
				assert.Greater(t, len(searchResp.Hits), 0)

				hit := searchResp.Hits[0]
				switch tt.resource {
				case blnkgo.ResourceLedgers:
					ledger, ok := hit.Document.(*blnkgo.Ledger)
					assert.True(t, ok)
					assert.NotEmpty(t, ledger.LedgerID)
				}
			}
			mockClient.AssertExpectations(t)
		})
	}
}
