package save

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/babenow/url_shortener/intrernal/config"
	resp "github.com/babenow/url_shortener/intrernal/lib/api/response"
	"github.com/babenow/url_shortener/intrernal/lib/logger/sl"
	"github.com/babenow/url_shortener/intrernal/lib/random"
	"github.com/babenow/url_shortener/intrernal/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const pkg = "handlers.url.save."

type URLSaver interface {
	Save(context.Context, model.Url) (int64, error)
	GetURLByAlias(ctx context.Context, alias string) (*model.Url, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	op := pkg + "New"
	var req Request

	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())), // TODO: привязка к роутеру chi
		)

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		if req.Alias == "" {
			req.Alias = generateAlias(r.Context(), log, saver)
		}
		m := model.Url{Alias: req.Alias, URL: req.URL}

		id, err := saver.Save(r.Context(), m)
		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: *resp.OK(),
			Alias:    m.Alias,
		})
	}
}

func generateAlias(ctx context.Context, log *slog.Logger, saver URLSaver) string {
	a := random.NewRandomString(config.Instance().AliasLength)
	if _, err := saver.GetURLByAlias(ctx, a); err != nil {
		return a
	}
	return generateAlias(ctx, log, saver)
}
