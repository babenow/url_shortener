package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/babenow/url_shortener/intrernal/lib/api/response"
	"github.com/babenow/url_shortener/intrernal/lib/logger/sl"
	"github.com/babenow/url_shortener/intrernal/model"
	"github.com/babenow/url_shortener/intrernal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const pkg = "http_server.handlers.redirect."

//go:generate go run github.com/vektra/mockery/v2@v2.37.1 --name=URLGetter
type URLGetter interface {
	GetURLByAlias(ctx context.Context, alias string) (*model.Url, error)
	AddRedirect(ctx context.Context, alias string) error
}

func New(ctx context.Context, log *slog.Logger, finder URLGetter) http.HandlerFunc {
	op := pkg + ".New"
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		m, err := finder.GetURLByAlias(ctx, alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", m.URL))
		if err := finder.AddRedirect(ctx, alias); err != nil {
			log.Error("can not add redirect", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		http.Redirect(w, r, m.URL, http.StatusFound)
	}
}
