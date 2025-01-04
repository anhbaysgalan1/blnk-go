package blnkgo_test

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) NewRequest(endpoint string, method string, body interface{}) (*http.Request, error) {
	args := m.Called(endpoint, method, body)
	if req, ok := args.Get(0).(*http.Request); ok || args.Get(0) == nil {
		return req, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClient) CallWithRetry(req *http.Request, v interface{}) (*http.Response, error) {
	args := m.Called(req, v)
	if resp, ok := args.Get(0).(*http.Response); ok || args.Get(0) == nil {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClient) NewFileUploadRequest(endpoint string, fileParam string, file interface{}, fileName string, fields map[string]string) (*http.Request, error) {
	args := m.Called(endpoint, fileParam, file, fileName, fields)
	if req, ok := args.Get(0).(*http.Request); ok || args.Get(0) == nil {
		return req, args.Error(1)
	}
	return nil, args.Error(1)
}
