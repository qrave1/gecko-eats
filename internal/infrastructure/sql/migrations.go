package sql

import (
	"database/sql"
	"embed"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(path string) error {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return err
	}
	defer db.Close()

	source := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationFiles,
		Root:       "migrations",
	}

	n, err := migrate.Exec(db, "sqlite3", source, migrate.Up)

	if err != nil {
		return err
	}

	log.Printf("Применено миграций: %d\n", n)

	return nil
}
