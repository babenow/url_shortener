package format

import (
	"fmt"
	"log/slog"

	"github.com/babenow/url_shortener/intrernal/lib/logger/sl"
)

func Err(operator string, err error) error {
	return fmt.Errorf("%s: %w", operator, err)
}

func CheckErr(operator string, log *slog.Logger, fn func() error) {
	if err := fn(); err != nil {
		log.Error("error handled", sl.Op(operator), sl.Err(err))
	}
}
