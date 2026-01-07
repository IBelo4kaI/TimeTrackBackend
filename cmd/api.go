package main

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"TimeTrack/internal/calendar"
	"TimeTrack/internal/report"
	"TimeTrack/internal/standard"
	types "TimeTrack/internal/type"
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

	standardService := standard.NewService(repo.New(app.db), app.db)
	standardHandler := standard.NewHandler(standardService, app.logger)

	typesService := types.NewService(repo.New(app.db), app.db)
	typesHandler := types.NewHandler(typesService, app.logger)

	v1 := fiber.Group("v1")
	// admin := v1.Group("/admin")

	report := v1.Group("/report")
	calendar := v1.Group("/calendar")
	vacation := v1.Group("/vacation")
	standard := v1.Group("/standard")
	types := v1.Group("/type")

	report.Get("/list/:user/:month/:year", reportHandler.List)
	report.Get("/monthstats/:user/:month/:year", reportHandler.MonthStats)
	report.Post("/create", reportHandler.Create)
	report.Post("/update", reportHandler.Update)
	report.Delete("/delete/:user/:day/:month/:year", reportHandler.Delete)

	vacation.Get("/list/:year", vacationHandler.ListAll)
	vacation.Get("/list/:user/:year", vacationHandler.List)
	vacation.Get("/stats/:user/:year", vacationHandler.Stats)
	vacation.Get("/years/:user", vacationHandler.Years)
	vacation.Post("/create", vacationHandler.Create)
	vacation.Post("/change-status", vacationHandler.ChangeStatus)
	vacation.Delete("/delete/:vacation", vacationHandler.Delete)

	calendar.Get("/list/:month/:year", calendarHandler.ListMonth)
	calendar.Get("/list/:year", calendarHandler.ListYear)
	calendar.Post("/create", calendarHandler.Create)

	types.Get("/list", typesHandler.List)

	standard.Post("/create", standardHandler.Create)
	standard.Post("/update", standardHandler.Update)
	standard.Get("/listforsetting/:year", standardHandler.ListForSetting)

	return fiber
}

func (app *application) run(f *fiber.App) error {
	return f.ListenTLS(app.config.addr, "cert.pem", "key.pem")
}
