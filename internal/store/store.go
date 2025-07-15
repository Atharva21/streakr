package store

import (
	"database/sql"
	"embed"
	"os"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations
var migrationsFS embed.FS

var (
	bootstrapStoreOnce sync.Once
	db                 *sql.DB
)

func BootstrapStore(dbPath string) {
	bootstrapStoreOnce.Do(func() {
		file, err := os.OpenFile(dbPath, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			panic("Failed to create streakr store: " + err.Error())
		}
		err = file.Close()
		if err != nil {
			panic("Failed to create streakr store: " + err.Error())
		}
		db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
		if err != nil {
			panic("Failed to open streakr store: " + err.Error())
		}
		if err = db.Ping(); err != nil {
			panic("Failed to connect to streakr store: " + err.Error())
		}

		d, err := iofs.New(migrationsFS, "migrations")
		if err != nil {
			panic("Failed to initialize streakr store migrations: " + err.Error())
		}
		m, err := migrate.NewWithSourceInstance("iofs", d, "sqlite3://"+dbPath+"?_foreign_keys=on")
		if err != nil {
			panic("Failed to initialize streakr store: " + err.Error())
		}
		defer m.Close()

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			panic("Failed to apply streakr store configurations: " + err.Error())
		}
	})
}
