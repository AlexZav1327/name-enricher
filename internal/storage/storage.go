package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

//go:embed migrations
var migrations embed.FS

type Postgres struct {
	db  *pgx.Conn
	dsn string
	log *logrus.Entry
}

func ConnectDB(ctx context.Context, dsn string, log *logrus.Logger) (*Postgres, error) {
	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect(ctx, dsn): %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db.Ping(ctx): %w", err)
	}

	return &Postgres{
		db:  db,
		dsn: dsn,
		log: log.WithField("module", "postgres"),
	}, nil
}

func (p *Postgres) Migrate(direction migrate.MigrationDirection) error {
	conn, err := sql.Open("pgx", p.dsn)
	if err != nil {
		return fmt.Errorf(`sql.Open("pgx", p.dsn): %w`, err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			p.log.Warningf("conn.Close(): %s", err)
		}
	}()

	assetDir := func() func(string) ([]string, error) {
		return func(path string) ([]string, error) {
			dirEntry, err := migrations.ReadDir(path)
			if err != nil {
				return nil, fmt.Errorf("migrations.ReadDir(path): %w", err)
			}

			entries := make([]string, 0)

			for _, e := range dirEntry {
				entries = append(entries, e.Name())
			}

			return entries, nil
		}
	}()

	asset := migrate.AssetMigrationSource{
		Asset:    migrations.ReadFile,
		AssetDir: assetDir,
		Dir:      "migrations",
	}

	_, err = migrate.Exec(conn, "postgres", asset, direction)
	if err != nil {
		return fmt.Errorf(`migrate.Exec(conn, "postgres", asset, direction): %w`, err)
	}

	return nil
}

func (p *Postgres) TruncateTable(ctx context.Context, table string) error {
	query := fmt.Sprintf(`TRUNCATE TABLE %s`, table)

	_, err := p.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("p.db.Exec(ctx, query): %w", err)
	}

	return nil
}
