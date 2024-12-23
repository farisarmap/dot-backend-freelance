package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/farisarmap/dot-backend-freelance/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	migrationDir := flag.String("dir", "migrations", "Path to migrations folder")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Please provide a migration command: up, down, drop, version, force, etc.")
	}
	command := flag.Arg(0)

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := config.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting DB: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting sql.DB from GORM: %v", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating Postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationDir),
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error on migrate up: %v", err)
		} else {
			log.Println("Migrate up done (or no change).")
		}

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error on migrate down: %v", err)
		} else {
			log.Println("Migrate down done (or no change).")
		}

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Error on migrate drop: %v", err)
		} else {
			log.Println("All tables dropped.")
		}

	case "force":
		if len(flag.Args()) < 2 {
			log.Fatal("force command requires a version number")
		}
		v := flag.Arg(1)
		ver, err := parseInt(v)
		if err != nil {
			log.Fatalf("Invalid version number: %v", v)
		}
		if err := m.Force(int(ver)); err != nil {
			log.Fatalf("Error on migrate force: %v", err)
		}
		log.Printf("Force set version to %d\n", ver)

	case "version":
		v, d, err := m.Version()
		if err != nil {
			log.Fatalf("Cannot get version: %v", err)
		}
		log.Printf("Current migration version: %d (dirty=%v)\n", v, d)

	default:
		log.Fatalf("Unknown command: %s. Valid commands: up, down, drop, force, version", command)
	}

	log.Println("Migration command completed.")
}

func parseInt(s string) (int64, error) {
	// parse string to int64
	var i int64
	_, err := fmt.Sscan(s, &i)
	return i, err
}
