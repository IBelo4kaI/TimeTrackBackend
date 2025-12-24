package calendar

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) List(c *fiber.Ctx) error {
	year, err := c.ParamsInt("year")
	if err != nil || year < 1900 || year > 2100 {
		return h.respondError(c, http.StatusBadRequest, "invalid year parameter")
	}
	
	month, err := c.ParamsInt("month")
	if err != nil || month < 1 || month > 12 {
		return h.respondError(c, http.StatusBadRequest, "invalid month parameter")
	}

	calendars, err := h.service.List(c.Context(), repo.GetDaysParams{
		Month: int32(month),
		Year:  int32(year),
	})

	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(calendars)
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// respondError - вспомогательный метод для отправки ошибок
func (h *Handler) respondError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
