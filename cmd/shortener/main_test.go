package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

func Test_shortenURL(t *testing.T) {
	type want struct {
		statusCode    int
		response      string
		contentType   string
		contentLength string
	}
	type request struct {
		method      string
		contentType string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Wrong path method",
			request: request{
				method:      http.MethodGet,
				contentType: "text/plain; charset=utf-8",
			},
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Wrong content type",
			request: request{
				method:      http.MethodPost,
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusUnsupportedMediaType,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Get short URL",
			request: request{
				method:      http.MethodPost,
				contentType: "text/plain",
			},
			want: want{
				statusCode:    http.StatusCreated,
				contentType:   "text/plain",
				contentLength: "30",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, "/", nil)
			request.Header.Set("Content-Type", tt.request.contentType)
			w := httptest.NewRecorder()
			shortenURL(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.contentLength, result.Header.Get("Content-Length"))
		})
	}
}

func Test_getURL(t *testing.T) {
	type want struct {
		statusCode    int
		response      string
		contentType   string
		contentLength string
	}
	type request struct {
		method      string
		contentType string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
