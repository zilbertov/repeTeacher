package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run ./tools/migrate <up|down|force>")
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://repe_teacher:repe_teacher@localhost:5432/repe_teacher?sslmode=disable"
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	instance, err := newMigrate(databaseURL, migrationsDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer instance.Close()

	switch os.Args[1] {
	case "up":
		err = instance.Up()
	case "down":
		err = instance.Down()
	case "force":
		if len(os.Args) < 3 {
			fmt.Println("usage: go run ./tools/migrate force <version>")
			os.Exit(1)
		}
		var version int
		_, err = fmt.Sscanf(os.Args[2], "%d", &version)
		if err == nil {
			err = instance.Force(version)
		}
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("migrations: no change")
		return
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("migrations:", os.Args[1], "done")
}

func newMigrate(databaseURL string, migrationsDir string) (*migrate.Migrate, error) {
	absoluteDir, err := filepath.Abs(filepath.Clean(migrationsDir))
	if err != nil {
		return nil, err
	}
	return migrate.New("file://"+absoluteDir, databaseURL)
}
