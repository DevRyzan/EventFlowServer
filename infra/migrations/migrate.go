package migrations

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var embedFS embed.FS

// Run executes all embedded SQL migrations in order.
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := fs.Glob(embedFS, "*.sql")
	if err != nil {
		return err
	}
	sort.Strings(entries)
	for _, name := range entries {
		sql, err := fs.ReadFile(embedFS, name)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return err
		}
		slog.Info("migration applied", "file", name)
	}
	return nil
}
