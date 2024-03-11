package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"sso/internal/storage/sqlite"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migration table")
	flag.Parse()

	validateFlags(storagePath, migrationsPath)

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte("adminadmin"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	_, err = storage.SaveUser(context.Background(), "admin", passHash)
	if err != nil {
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}

func validateFlags(storagePath, migrationsPath string) {
	if storagePath == "" {
		panic("storage-path is required")
	}

	if migrationsPath == "" {
		panic("migrations-path is required")
	}
}
