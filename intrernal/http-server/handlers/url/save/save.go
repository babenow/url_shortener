package save

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	resp "github.com/babenow/url_shortener/intrernal/lib/api/response"
	"github.com/babenow/url_shortener/intrernal/lib/logger/sl"
	"github.com/babenow/url_shortener/intrernal/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const pkg = "handlers.url.save."

type URLSaver interface {
	Save(context.Context, model.Url) (int64, error)
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
			req.Alias = fmt.Sprint(rand.Int63n(999999)) // FIXME: генерация ссылки
		}
		m := model.Url{Alias: req.Alias, URL: req.URL}

		_, err := saver.Save(r.Context(), m)
		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		render.JSON(w, r, resp.OK())
	}
}
