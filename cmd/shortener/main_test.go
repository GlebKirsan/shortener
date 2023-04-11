package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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
		name         string
		db           map[string]URL
		reverseIndex map[URL]string
		request      request
		want         want
	}{
		{
			name:         "Wrong path method",
			db:           map[string]URL{},
			reverseIndex: map[URL]string{},
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
			name:         "Get short URL",
			db:           map[string]URL{},
			reverseIndex: map[URL]string{},
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
		{
			name: "Shorten existing URL",
			db: map[string]URL{
				"short": URL("full url"),
			},
			reverseIndex: map[URL]string{
				URL("full url"): "short",
			},
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
			db = tt.db
			reverseIndex = tt.reverseIndex
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
		statusCode int
		location   string
	}
	type request struct {
		method string
		path   string
		vars   map[string]string
	}
	tests := []struct {
		name         string
		db           map[string]URL
		reverseIndex map[URL]string
		request      request
		want         want
	}{
		{
			name:         "Wrong request method",
			db:           map[string]URL{},
			reverseIndex: map[URL]string{},
			request: request{
				method: http.MethodPost,
				path:   "/someId",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				location:   "",
			},
		},
		{
			name:         "Miss id in path",
			db:           map[string]URL{},
			reverseIndex: map[URL]string{},
			request: request{
				method: http.MethodGet,
				path:   "/",
				vars:   map[string]string{},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:         "Id not found",
			db:           map[string]URL{},
			reverseIndex: map[URL]string{},
			request: request{
				method: http.MethodGet,
				path:   "/someId",
				vars:   map[string]string{"id": "trueId"},
			},
			want: want{
				statusCode: http.StatusNotFound,
				location:   "",
			},
		},
		{
			name: "Redirect is ok",
			db: map[string]URL{
				"randstrI": URL("https://practicum.yandex.ru/"),
			},
			reverseIndex: map[URL]string{
				URL("https://practicum.yandex.ru/"): "randstrI",
			},
			request: request{
				method: http.MethodGet,
				path:   "/randstrI",
				vars:   map[string]string{"id": "randstrI"},
			},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://practicum.yandex.ru/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.path, nil)

			w := httptest.NewRecorder()
			db = tt.db
			reverseIndex = tt.reverseIndex
			request = mux.SetURLVars(request, tt.request.vars)

			getURL(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
