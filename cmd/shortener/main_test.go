package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generateString(t *testing.T) {
	tests := []struct {
		name         string
		length       int
		wantedLength int
	}{
		{
			name:         "Create random string with length 6",
			length:       6,
			wantedLength: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateString(tt.length); len(got) != tt.wantedLength {
				t.Errorf("generateString() = %v, want %v", got, tt.wantedLength)
			}
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func Test_shortenURL(t *testing.T) {
	ts := httptest.NewServer(URLRouter())
	defer ts.Close()

	type want struct {
		statusCode    int
		contentType   string
		contentLength string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Get short URL",
			want: want{
				statusCode:    http.StatusCreated,
				contentType:   "text/plain",
				contentLength: "30",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, http.MethodPost, "/")
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.contentLength, resp.Header.Get("Content-Length"))
		})
	}
}

func Test_getURL(t *testing.T) {
	ts := httptest.NewServer(URLRouter())
	defer ts.Close()

	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name         string
		db           map[string]URL
		reverseIndex map[URL]string
		path         string
		want         want
	}{
		{
			name: "Id not found",
			path: "/someId",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUp(tt.db, tt.reverseIndex)
			resp, _ := testRequest(t, ts, http.MethodGet, tt.path)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
		})
	}
}
