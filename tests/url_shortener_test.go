package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/babenow/url_shortener/intrernal/http-server/handlers/url/save"
	"github.com/babenow/url_shortener/intrernal/lib/api"
	"github.com/babenow/url_shortener/intrernal/lib/random"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(6),
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")

}

func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		desc  string
		url   string
		alias string
		error string
	}{
		{
			desc:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			desc:  "Invalid URL",
			url:   "invalid",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			desc:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}
	for _, tC := range testCases {
		tc := tC
		t.Run(tC.desc, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirect(t, alias, tc.url)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
