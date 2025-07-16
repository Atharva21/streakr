package store

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"sync"

	"github.com/Atharva21/streakr/internal/shutdown"
	"github.com/Atharva21/streakr/internal/store/generated"
	"github.com/Atharva21/streakr/internal/util"
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
	queries            *generated.Queries
)

func BootstrapStore(dbPath string) {
	bootstrapStoreOnce.Do(func() {
		file, err := os.OpenFile(dbPath, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			util.ErrorAndExitGeneric(err)
		}
		err = file.Close()
		if err != nil {
			util.ErrorAndExitGeneric(err)
		}
		db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
		if err != nil {
			util.ErrorAndExitGeneric(err)
		}
		shutdown.RegisterCleanupHook(func() error {
			return db.Close()
		})
		if err = db.Ping(); err != nil {
			util.ErrorAndExitGeneric(err)
		}

		d, err := iofs.New(migrationsFS, "migrations")
		if err != nil {
			util.ErrorAndExitGeneric(err)
		}
		m, err := migrate.NewWithSourceInstance("iofs", d, "sqlite3://"+dbPath+"?_foreign_keys=on")
		if err != nil {
			util.ErrorAndExitGeneric(err)
		}
		defer m.Close()
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			util.ErrorAndExitGeneric(err)
		}

		queries = generated.New(db)

	})
}

func GetDB() *sql.DB {
	if db == nil {
		util.ErrorAndExitGeneric(fmt.Errorf("DB not initialized, GetDB called before bootstrapStore"))
	}
	return db
}

func GetQueries() *generated.Queries {
	if queries == nil {
		util.ErrorAndExitGeneric(fmt.Errorf("Queries instance not found, GetQueries called before code is generated"))
	}
	return queries
}
