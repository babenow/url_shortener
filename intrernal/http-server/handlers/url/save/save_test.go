package save_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/babenow/url_shortener/intrernal/http-server/handlers/url/save"
	"github.com/babenow/url_shortener/intrernal/http-server/handlers/url/save/mocks"
	"github.com/babenow/url_shortener/intrernal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := save.ConfigureConfigPath(); err != nil {
		fmt.Println("can not configure config path: %w", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestSaveHandler(t *testing.T) {
	testCases := []struct {
		desc      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			desc:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			desc:  "Empty Alias",
			alias: "",
			url:   "https://google.com",
			// mockError: errors.New("unexpected error"),
		},
		{
			desc:      "Empty URL",
			alias:     "some_alias",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			desc:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			desc:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
	}
	for _, tC := range testCases {
		tc := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.alias == "" {
				urlSaverMock.On("GetURLByAlias", context.Background(), mock.AnythingOfType("string")).
					Return(nil, errors.New("unexpected error")).
					Once()
			}

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("Save", context.Background(), mock.AnythingOfType("model.Url")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewSlogDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"alias": "%s", "url": "%s"}`, tc.alias, tc.url)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
