package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/babenow/url_shortener/intrernal/lib/helper/format"
	"github.com/babenow/url_shortener/intrernal/model"
	"github.com/babenow/url_shortener/intrernal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

var pkg = "storage.sqlite.SqliteURLStorage."

const (
	tableName = "url"
)

// dbUrl структура прослойки
type dbUrl struct {
	ID    int64  `db:"id"`
	Alias string `db:"alias"`
	URL   string `db:"url"`
}

type SqliteURLStorage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func newSqliteURLStorage(db *sqlx.DB, log *slog.Logger) *SqliteURLStorage {
	return &SqliteURLStorage{db, log}
}

func (s *SqliteURLStorage) Save(ctx context.Context, model model.Url) (int64, error) {
	op := pkg + "Save"
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, format.Err(op, err)
	}
	defer format.CheckErr(op, s.log, conn.Close)

	// TODO: Возможность обновлять URL в случае, если что-то изменилось

	var id int64

	row := s.db.QueryRowContext(ctx, fmt.Sprintf(`INSERT INTO %s(alias,url)VALUES($1,$2) RETURNING id;`, tableName),
		model.Alias,
		model.URL,
	)

	if err := row.Err(); err != nil {
		return 0, format.Err(op, err)
	}

	if err := row.Scan(&id); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, format.Err(op, storage.ErrURLExists)
		}
		return 0, format.Err(op, err)
	}

	return id, nil
}

// GetURLByAlias
func (s *SqliteURLStorage) GetURLByAlias(ctx context.Context, alias string) (*model.Url, error) {
	op := pkg + "GetAlias"

	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, format.Err(op, err)
	}
	defer format.CheckErr(op, s.log, conn.Close)
	//
	var url dbUrl

	if err := s.db.GetContext(ctx, &url, fmt.Sprintf(`SELECT * FROM %s WHERE alias=$1`, tableName), alias); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, format.Err(op, storage.ErrURLNotFound)
		}
		return nil, format.Err(op, err)
	}

	return (*model.Url)(&url), nil
}

func (s *SqliteURLStorage) DeleteURLByID(ctx context.Context, id int64) error {
	op := pkg + "DeleteURLByID"
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return format.Err(op, err)
	}
	defer format.CheckErr(op, s.log, conn.Close)

	if _, err := s.db.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, tableName), id); err != nil {
		return format.Err(op, err)
	}

	return nil
}

func (s *SqliteURLStorage) DeleteURLByAlias(ctx context.Context, alias string) error {
	op := pkg + "DeleteURLByAlias"
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return format.Err(op, err)
	}
	defer format.CheckErr(op, s.log, conn.Close)

	if _, err := s.db.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s WHERE alias=$1`, tableName), alias); err != nil {
		return format.Err(op, err)
	}

	return nil
}
