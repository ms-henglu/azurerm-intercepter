package interceptor_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ms-henglu/azurerm-interceptor/interceptor"
)

func Test_HandleRequest(t *testing.T) {
	testcases := []struct {
		Name     string
		Request  *http.Request
		Response *http.Response
	}{
		{
			Name: "CheckNameAvailability",
			Request: &http.Request{
				Method: "POST",
				URL:    &url.URL{Path: "/subscriptions/123/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/checkNameAvailability"},
			},
			Response: &http.Response{
				StatusCode: 200,
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"nameAvailable":true}`)),
			},
		},

		{
			Name: "GetResource which does not exist",
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/subscriptions/123/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/myStorageAccount"},
			},
			Response: &http.Response{
				StatusCode: 404,
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
				Body: http.NoBody,
			},
		},

		{
			Name: "CreateResource",
			Request: &http.Request{
				Method: "PUT",
				URL:    &url.URL{Path: "/subscriptions/123/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/myStorageAccount"},
				Body:   io.NopCloser(strings.NewReader(`{"location":"westus"}`)),
			},
			Response: &http.Response{
				StatusCode: 400,
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"code":"InterceptedError","message":"InterceptedError","target":null,"details":null,"innererror":{"body":"{\"location\":\"westus\"}","url":"/subscriptions/123/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/myStorageAccount"},"additionalInfo":null}`)),
			},
		},
		{
			Name: "GetResource which exists",
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/subscriptions/123/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/myStorageAccount"},
			},
			Response: &http.Response{
				StatusCode: 200,
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"location":"westus"}`)),
			},
		},
	}

	for _, tc := range testcases {
		t.Logf("Running test case: %s", tc.Name)
		resp, err := interceptor.HandleRequest(tc.Request)
		if err != nil {
			t.Fatalf("Failed to handle request: %v", err)
		}

		if resp == nil {
			t.Fatalf("Expected a response, got nil")
		}

		if resp.StatusCode != tc.Response.StatusCode {
			t.Fatalf("Expected status code %d, got %d", tc.Response.StatusCode, resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != tc.Response.Header.Get("Content-Type") {
			t.Fatalf("Expected Content-Type %s, got %s", tc.Response.Header.Get("Content-Type"), resp.Header.Get("Content-Type"))
		}

		expectedBody, err := io.ReadAll(tc.Response.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		actualBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		if string(expectedBody) != string(actualBody) {
			t.Fatalf("Expected body %s, got %s", string(expectedBody), string(actualBody))
		}
	}
}
