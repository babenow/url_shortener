package sqlite

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/babenow/url_shortener/intrernal/config"
	"github.com/babenow/url_shortener/intrernal/lib/helper/format"
	goose "github.com/pressly/goose/v3"

	// "github.com/babenow/url_shortener/intrernal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type SqliteStorage struct {
	db         *sqlx.DB
	log        *slog.Logger
	urlStorage *SqliteURLStorage
}

func New(ctx context.Context, log *slog.Logger) (*SqliteStorage, error) {
	const op = "storage.sqlite.New"

	db, err := sqlx.Connect("sqlite3", config.Instance().StoragePath)
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, format.Err(op, err)
	}
	if config.Instance().Goose.PrintStatus {
		if err := goose.Status(db.DB, "migrations"); err != nil {
			return nil, format.Err(op, err)
		}
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return nil, format.Err(op, err)
	}

	return &SqliteStorage{db, log, nil}, nil
}

func (s *SqliteStorage) UrlStorage() *SqliteURLStorage {
	if s.urlStorage == nil {
		s.urlStorage = newSqliteURLStorage(s.db, s.log)
	}
	return s.urlStorage
}
