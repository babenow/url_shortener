package save

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/babenow/url_shortener/intrernal/config"
	resp "github.com/babenow/url_shortener/intrernal/lib/api/response"
	"github.com/babenow/url_shortener/intrernal/lib/logger/sl"
	"github.com/babenow/url_shortener/intrernal/lib/random"
	"github.com/babenow/url_shortener/intrernal/model"
	"github.com/babenow/url_shortener/intrernal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const pkg = "handlers.url.save."

//go:generate go run github.com/vektra/mockery/v2@v2.37.1 --name=URLSaver
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

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(config.Instance().AliasLength)
		}
		m := model.Url{Alias: alias, URL: req.URL}

		id, err := saver.Save(r.Context(), m)
		if errors.Is(err, storage.ErrURLExists) {
			log.Error("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

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

// func generateAlias(ctx context.Context, saver URLSaver) string {
// 	a := random.NewRandomString(config.Instance().AliasLength)
// 	if _, err := saver.GetURLByAlias(ctx, a); err != nil {
// 		return a
// 	}
// 	return generateAlias(ctx, saver)
// }
