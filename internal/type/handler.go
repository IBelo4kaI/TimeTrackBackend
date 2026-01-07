package types

import (
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
	types, err := h.service.List(c.Context())
	if err != nil {
		h.logger.Error("failed to get vacations",
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to get vacations")
	}

	return c.JSON(types)
}

// ErrorResponse представляет стандартный формат ошибки
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
