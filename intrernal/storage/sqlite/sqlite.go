package sqlite

import (
	"context"
	"embed"
	"fmt"

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
	urlStorage *SqliteURLStorage
}

func New(ctx context.Context) (*SqliteStorage, error) {
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

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return nil, format.Err(op, err)
	}

	return &SqliteStorage{db, nil}, nil
}

func (s *SqliteStorage) UrlStorage() *SqliteURLStorage {
	if s.urlStorage == nil {
		s.urlStorage = newSqliteURLStorage(s.db)
	}
	return s.urlStorage
}

// // AllURL получить все URL
// func (s *Sqlite) AllURL(ctx context.Context) ([]model.Url, error) {
// 	op := "storage.sqlite.AllURL"
// 	conn, err := s.db.Connx(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	defer conn.Close()

// 	var urls []dbUrl
// 	var models []model.Url

// 	if err := s.db.SelectContext(ctx, &urls, "SELECT * FROM url"); err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	for _, url := range urls {
// 		u := model.Url(url)
// 		models = append(models, u)
// 	}

// 	return models, nil
// }
