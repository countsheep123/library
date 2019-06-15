package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/countsheep123/library/obj"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type db interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type Handler struct {
	readDB  db
	writeDB db
}

type Transaction struct {
	Handler
}

func New(r, w *sql.DB) (*Handler, error) {
	return &Handler{
		readDB:  r,
		writeDB: w,
	}, nil
}

func (h *Handler) Transact(ctx context.Context, txFunc func(*Transaction) error) (err error) {
	switch h.writeDB.(type) {
	case *sql.DB:
		// ok
	case *sql.Tx:
		err = fmt.Errorf("handler has already begun transaction")
		return
	default:
		err = fmt.Errorf("invalid type")
		return
	}

	tx, err := h.writeDB.(*sql.DB).BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	return txFunc(&Transaction{
		Handler: Handler{
			readDB:  tx,
			writeDB: tx,
		},
	})
}

var (
	replacer = strings.NewReplacer(`\`, `\\\`, `%`, `\%`, `_`, `\_`)
)

// sanitize characters for LIKE clause
// \, %, _
func sanitize(input string) string {
	return replacer.Replace(input)
}

func (h *Handler) Insert(ctx context.Context, table string, columns []string, values []interface{}) error {
	query, args, err := psql.Insert(table).Columns(columns...).Values(values...).ToSql()
	if err != nil {
		return err
	}

	zap.S().Debug(query, args)

	if _, err := h.writeDB.ExecContext(ctx, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique violation
				zap.S().Debug(pqErr)
				return obj.Duplicate{}
			default:
				return err
			}
		}
	}
	return nil
}

func (h *Handler) Update(ctx context.Context, table string, kv map[string]interface{}, filters map[string]string) error {
	f := squirrel.And{}

	for k, v := range filters {
		f = append(f, squirrel.Eq{
			k: v,
		})
	}

	query, args, err := psql.Update(table).SetMap(kv).Where(f).ToSql()
	if err != nil {
		return err
	}

	zap.S().Debug(query, args)

	if _, err := h.writeDB.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Delete(ctx context.Context, table string, filters map[string]string) error {
	f := squirrel.And{}

	for k, v := range filters {
		f = append(f, squirrel.Eq{
			k: v,
		})
	}

	query, args, err := psql.Delete(table).Where(f).ToSql()
	if err != nil {
		return err
	}

	zap.S().Debug(query, args)

	result, err := h.writeDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return obj.NotFound{}
	}

	return nil
}
