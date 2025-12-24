package main

import (
	"TimeTrack/internal/env"
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env := env.Env{}
	env.Init()

	cfg := config{
		addr: env.GetAddr(),
		db: dbConfig{
			dsn: env.GetDbString(),
		},
	}

	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		panic("error con database")
	}

	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	app := application{
		config: cfg,
		db:     db,
		logger: logger,
	}

	// Run the application
	app.run(app.mount())
}
