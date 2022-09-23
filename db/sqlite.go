package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	_ "modernc.org/sqlite"
)

func NewRepository(file string) (*Repository, error) {

	// Create db file if not exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		file, err := os.Create(file)
		if err != nil {
			return nil, fmt.Errorf("db create file: %v", err)
		}
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("db close file: %v", err)
		}
	}

	// Open database connection
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, fmt.Errorf("db open: %v", err)
	}

	// Run migration
	var repo = &Repository{
		db: db,
	}
	if err := repo.Migrate(file); err != nil {
		return nil, fmt.Errorf("db migrate: %v", err)
	}

	return repo, nil
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) Migrate(file string) error {
	driver, err := sqlite.WithInstance(r.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("db migrate driver: %v", err)
	}

	source, err := bindata.WithInstance(bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		}))
	if err != nil {
		return fmt.Errorf("db migrate source: %s", err)
	}

	migrator, err := migrate.NewWithInstance("", source, fmt.Sprintf("sqlite://%s", file), driver)
	if err != nil {
		return fmt.Errorf("db migrate: %s", err)
	}

	if err := migrator.Up(); err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}
