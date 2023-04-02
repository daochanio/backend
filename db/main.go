package main

import (
	"embed"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	if env := os.Getenv("ENV"); env == "dev" {
		if err := godotenv.Load(".env/.env.db.dev"); err != nil {
			panic(err)
		}
	}

	db, err := goose.OpenDBWithDriver("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
