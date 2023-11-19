package redirect_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/babenow/url_shortener/intrernal/http-server/handlers/redirect"
	"github.com/babenow/url_shortener/intrernal/http-server/handlers/redirect/mocks"
	"github.com/babenow/url_shortener/intrernal/lib/api"
	"github.com/babenow/url_shortener/intrernal/lib/logger/handlers/slogdiscard"
	"github.com/babenow/url_shortener/intrernal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirect(t *testing.T) {
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
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			urlGetter := mocks.NewURLGetter(t)

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(context.Background(), slogdiscard.NewSlogDiscardLogger(), urlGetter))

			ts := httptest.NewServer(r)
			defer ts.Close()

			if tc.respError == "" || tc.mockError != nil {
				urlGetter.On("GetURLByAlias", context.Background(), tc.alias).
					Return(&model.Url{ID: int64(1), Alias: tc.alias, URL: tc.url}, nil).
					Once()
			}

			redirectedToUrl, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			assert.Equal(t, tc.url, redirectedToUrl)
		})
	}
}
