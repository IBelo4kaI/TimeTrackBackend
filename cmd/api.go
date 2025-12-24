package main

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"TimeTrack/internal/calendar"
	"TimeTrack/internal/report"
	"TimeTrack/internal/vacation"
	"database/sql"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type application struct {
	config config
	db     *sql.DB
	logger *slog.Logger
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() *fiber.App {
	fiber := fiber.New(fiber.Config{
		Prefork: true,
		// EnablePrintRoutes: true,
	})

	fiber.Use(cors.New())
	fiber.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} | ${latency} | ${status} - ${method} ${path} \n",
	}))

	reportService := report.NewService(repo.New(app.db), app.db)
	reportHandler := report.NewHandler(reportService, app.logger)

	vacationService := vacation.NewService(repo.New(app.db), app.db)
	vacationHandler := vacation.NewHandler(vacationService, app.logger)

	calendarService := calendar.NewService(repo.New(app.db), app.db)
	calendarHandler := calendar.NewHandler(calendarService, app.logger)

	v1 := fiber.Group("v1")
	// admin := v1.Group("/admin")

	report := v1.Group("/report")
	calendar := v1.Group("/calendar")
	vacation := v1.Group("/vacation")

	report.Get("/list/:user/:month/:year", reportHandler.List)
	report.Post("/create", reportHandler.Create)
	report.Post("/update", reportHandler.Update)
	report.Delete("/delete/:user/:day/:month/:year", reportHandler.Delete)

	vacation.Get("/list/:year", vacationHandler.ListAll)
	vacation.Get("/list/:user/:year", vacationHandler.List)
	vacation.Get("/stats/:user/:year", vacationHandler.Stats)
	vacation.Get("/years/:user", vacationHandler.Years)
	vacation.Post("/create", vacationHandler.Create)
	vacation.Post("/change-status", vacationHandler.ChangeStatus)

	calendar.Get("/list/:month/:year", calendarHandler.List)

	return fiber
}

func (app *application) run(f *fiber.App) error {
	return f.ListenTLS(app.config.addr, "cert.pem", "key.pem")
}
