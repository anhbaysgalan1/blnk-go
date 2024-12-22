package blnkgo_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	blnkgo "github.com/blnkfinance/blnk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupIdentityService() (*MockClient, *blnkgo.IdentityService) {
	mockClient := &MockClient{}
	svc := blnkgo.NewIdentityService(mockClient)
	return mockClient, svc
}

func TestIdentityService_Create(t *testing.T) {
	mockClient, svc := setupIdentityService()

	identity := blnkgo.Identity{
		IdentityType: blnkgo.Individual,
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
		PhoneNumber:  "1234567890",
		Category:     "customer",
		Street:       "123 Main St",
		Country:      "USA",
		State:        "CA",
		PostCode:     "90001",
		City:         "Los Angeles",
		DOB:          &time.Time{},
		Gender:       "Male",
		Nationality:  "Nigerian",
	}

	t.Run("successful creation", func(t *testing.T) {
		expectedResponse := &blnkgo.IdentityResponse{
			IdentityId: "12345",
			CreatedAt:  time.Now().Format(time.RFC3339),
			Identity:   identity,
		}

		mockClient.On("NewRequest", "identities", http.MethodPost, identity).Return(&http.Request{}, nil)
		mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp := args.Get(1).(*blnkgo.IdentityResponse)
			*resp = *expectedResponse
		}).Return(&http.Response{}, nil)

		resp, httpResp, err := svc.Create(identity)
		assert.NoError(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, expectedResponse, resp)
		mockClient.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		invalidIdentity := blnkgo.Identity{
			EmailAddress: "",
			PhoneNumber:  "1234567890",
		}

		resp, httpResp, err := svc.Create(invalidIdentity)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Nil(t, httpResp)
	})
}

func TestIdentityService_ServerError(t *testing.T) {
	mockClient, svc := setupIdentityService()

	identity := blnkgo.Identity{
		IdentityType:     blnkgo.Organization,
		OrganizationName: "ACME Inc",
	}

	mockClient.On("NewRequest", "identities", http.MethodPost, identity).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError}, errors.New("server error"))

	resp, httpResp, err := svc.Create(identity)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.NotNil(t, httpResp)
	assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
	assert.Contains(t, err.Error(), "server error")
	mockClient.AssertExpectations(t)
}
func TestIdentityService_Get(t *testing.T) {
	mockClient, svc := setupIdentityService()

	identityId := "12345"
	expectedResponse := &blnkgo.IdentityResponse{
		IdentityId: identityId,
		CreatedAt:  time.Now().Format(time.RFC3339),
		Identity: blnkgo.Identity{
			IdentityType: blnkgo.Individual,
			FirstName:    "John",
			LastName:     "Doe",
			EmailAddress: "john.doe@example.com",
			PhoneNumber:  "1234567890",
			Category:     "customer",
			Street:       "123 Main St",
			Country:      "USA",
			State:        "CA",
			PostCode:     "90001",
			City:         "Los Angeles",
			DOB:          &time.Time{},
			Gender:       "Male",
			Nationality:  "Nigerian",
		},
	}

	t.Run("successful get", func(t *testing.T) {
		mockClient.On("NewRequest", fmt.Sprintf("identities/%s", identityId), http.MethodGet, nil).Return(&http.Request{}, nil)
		mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp := args.Get(1).(*blnkgo.IdentityResponse)
			*resp = *expectedResponse
		}).Return(&http.Response{}, nil)

		resp, httpResp, err := svc.Get(identityId)
		assert.NoError(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, expectedResponse, resp)
		mockClient.AssertExpectations(t)
	})
}

func TestIdentityService_Get_NotFound(t *testing.T) {
	mockClient, svc := setupIdentityService()
	identityId := "12345"
	mockClient.On("NewRequest", fmt.Sprintf("identities/%s", identityId), http.MethodGet, nil).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusNotFound}, errors.New("not found"))

	resp, httpResp, err := svc.Get(identityId)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.NotNil(t, httpResp)
	assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
	assert.Contains(t, err.Error(), "not found")
	mockClient.AssertExpectations(t)
}

func TestIdentityService_Get_ServerError(t *testing.T) {
	mockClient, svc := setupIdentityService()
	identityId := "12345"
	mockClient.On("NewRequest", fmt.Sprintf("identities/%s", identityId), http.MethodGet, nil).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError}, errors.New("server error"))

	resp, httpResp, err := svc.Get(identityId)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.NotNil(t, httpResp)
	assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
	assert.Contains(t, err.Error(), "server error")
}

func TestIdentityService_List(t *testing.T) {
	mockClient, svc := setupIdentityService()

	t.Run("successful list", func(t *testing.T) {
		expectedResponse := []*blnkgo.IdentityResponse{
			{
				IdentityId: "12345",
				CreatedAt:  time.Now().Format(time.RFC3339),
				Identity: blnkgo.Identity{
					IdentityType: blnkgo.Individual,
					FirstName:    "John",
					LastName:     "Doe",
					EmailAddress: "john@example.com",
					PhoneNumber:  "1234567890",
					Category:     "customer",
					Street:       "123 Main St",
					Country:      "USA",
					State:        "CA",
					PostCode:     "90001",
					City:         "Los Angeles",
				},
			},
			{
				IdentityId: "67890",
				CreatedAt:  time.Now().Format(time.RFC3339),
				Identity: blnkgo.Identity{
					IdentityType:     blnkgo.Organization,
					OrganizationName: "ACME Inc",
					EmailAddress:     "contact@acme.com",
					PhoneNumber:      "0987654321",
					Category:         "business",
					Street:           "456 Corp Ave",
					Country:          "USA",
					State:            "NY",
					PostCode:         "10001",
					City:             "New York",
				},
			},
		}

		mockClient.On("NewRequest", "identities", http.MethodGet, nil).Return(&http.Request{}, nil)
		mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp := args.Get(1).(*[]*blnkgo.IdentityResponse)
			*resp = expectedResponse
		}).Return(&http.Response{}, nil)

		resp, httpResp, err := svc.List()
		assert.NoError(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, expectedResponse, resp)
		mockClient.AssertExpectations(t)
	})
}

func TestIdentityService_List_ServerError(t *testing.T) {
	mockClient, svc := setupIdentityService()

	mockClient.On("NewRequest", "identities", http.MethodGet, nil).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError}, errors.New("server error"))

	resp, httpResp, err := svc.List()
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.NotNil(t, httpResp)
	assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
	assert.Contains(t, err.Error(), "server error")
	mockClient.AssertExpectations(t)
}

func TestIdentityService_Update_Successful(t *testing.T) {
	mockClient, svc := setupIdentityService()
	identityId := "12345"

	updateIdentity := &blnkgo.Identity{
		FirstName:    "Jane",
		LastName:     "Doe",
		EmailAddress: "jane.doe@example.com",
		PhoneNumber:  "0987654321",
		Category:     "customer",
		Street:       "456 Oak St",
		Country:      "USA",
		State:        "NY",
		PostCode:     "10001",
		City:         "New York",
	}

	expectedResponse := &blnkgo.IdentityResponse{
		IdentityId: identityId,
		CreatedAt:  time.Now().Format(time.RFC3339),
		Identity:   *updateIdentity,
	}

	mockClient.On("NewRequest", fmt.Sprintf("identities/%s", identityId), http.MethodPut, updateIdentity).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		resp := args.Get(1).(**blnkgo.IdentityResponse)
		*resp = expectedResponse
	}).Return(&http.Response{}, nil)

	resp, httpResp, err := svc.Update(identityId, updateIdentity)
	assert.NoError(t, err)
	assert.NotNil(t, httpResp)
	assert.Equal(t, expectedResponse, resp)
	mockClient.AssertExpectations(t)
}

func TestIdentityService_Update_NotFound(t *testing.T) {
	mockClient, svc := setupIdentityService()
	identityId := "12345"

	updateIdentity := &blnkgo.Identity{
		FirstName: "Jane",
	}

	mockClient.On("NewRequest", fmt.Sprintf("identities/%s", identityId), http.MethodPut, updateIdentity).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusNotFound}, errors.New("not found"))

	resp, httpResp, err := svc.Update(identityId, updateIdentity)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.NotNil(t, httpResp)
	assert.Equal(t, http.StatusNotFound, httpResp.StatusCode)
	assert.Contains(t, err.Error(), "not found")
	mockClient.AssertExpectations(t)
}

func TestIdentityService_Update_ServerError(t *testing.T) {
	mockClient, svc := setupIdentityService()
	identityId := "12345"

	updateIdentity := &blnkgo.Identity{
		FirstName: "Jane",
	}

	mockClient.On("NewRequest", fmt.Sprintf("identities/%s", identityId), http.MethodPut, updateIdentity).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError}, errors.New("server error"))

	resp, httpResp, err := svc.Update(identityId, updateIdentity)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.NotNil(t, httpResp)
	assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
	assert.Contains(t, err.Error(), "server error")
	mockClient.AssertExpectations(t)
}
func TestIdentityService_Create_EmptyPayload(t *testing.T) {
	_, svc := setupIdentityService()

	emptyIdentity := blnkgo.Identity{}
	resp, httpResp, err := svc.Create(emptyIdentity)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Nil(t, httpResp)
}
