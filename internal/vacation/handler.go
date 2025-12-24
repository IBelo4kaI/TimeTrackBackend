package vacation

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	userID := c.Params("user")
	if userID == "" {
		return h.respondError(c, http.StatusBadRequest, "user ID is required")
	}

	year, err := c.ParamsInt("year")
	if err != nil || year < 1900 || year > 2100 {
		return h.respondError(c, http.StatusBadRequest, "invalid year parameter")
	}

	vacations, err := h.service.List(c.Context(), userID, int32(year))
	if err != nil {
		h.logger.Error("failed to get vacations",
			slog.String("user_id", userID),
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to get vacations")
	}

	return c.JSON(vacations)
}

func (h *Handler) ListAll(c *fiber.Ctx) error {
	year, err := c.ParamsInt("year")
	if err != nil || year < 1900 || year > 2100 {
		return h.respondError(c, http.StatusBadRequest, "invalid year parameter")
	}

	vacations, err := h.service.ListAll(c.Context(), int32(year))
	if err != nil {
		h.logger.Error("failed to get vacations",
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to get vacations")
	}

	return c.JSON(vacations)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	type createRequest struct {
		UserID      string                    `json:"userId"`
		StartDate   time.Time                 `json:"startDate"`
		EndDate     time.Time                 `json:"endDate"`
		Year        int32                     `json:"year"`
		Description string                    `json:"description"`
		Status      repo.ReportVacationStatus `json:"status"`
	}

	var req createRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		return h.respondError(c, http.StatusBadRequest, "invalid request body")
	}

	description := sql.NullString{
		String: req.Description,
		Valid:  req.Description != "",
	}

	vacation, err := h.service.Create(c.Context(), repo.CreateVacationParams{
		UserID:      req.UserID,
		ID:          uuid.NewString(),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Year:        req.Year,
		Description: description,
		Status:      req.Status,
	})

	if err != nil {
		h.logger.Error("failed to get vacation stats",
			slog.String("user_id", req.UserID),
			slog.Time("startDate", req.StartDate),
			slog.Time("endDate", req.EndDate),
			slog.Int("year", int(req.Year)),
			slog.String("desc", description.String),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to create vacation")
	}

	return c.JSON(vacation)
}

func (h *Handler) Stats(c *fiber.Ctx) error {
	userID := c.Params("user")
	if userID == "" {
		return h.respondError(c, http.StatusBadRequest, "user ID is required")
	}

	year, err := c.ParamsInt("year")
	if err != nil || year < 1900 || year > 2100 {
		return h.respondError(c, http.StatusBadRequest, "invalid year parameter")
	}

	stats, err := h.service.Stats(c.Context(), userID, int32(year))
	if err != nil {
		h.logger.Error("failed to get vacation stats",
			slog.String("user_id", userID),
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to get vacation stats")
	}

	return c.JSON(stats)
}

func (h *Handler) ChangeStatus(c *fiber.Ctx) error {
	type changeStatus struct {
		ID     string                    `json:"id"`
		Status repo.ReportVacationStatus `json:"status"`
	}

	var req changeStatus

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		return h.respondError(c, http.StatusBadRequest, "invalid request body")
	}

	if err := h.service.ChangeStatus(c.Context(), repo.UpdateVacationStatusParams{
		ID:     req.ID,
		Status: req.Status,
	}); err != nil {
		h.logger.Warn("filed change status", slog.String("error", err.Error()))
		return h.respondError(c, http.StatusBadRequest, "filed change status")
	}

	c.Status(http.StatusOK)
	return nil
}

func (h *Handler) Years(c *fiber.Ctx) error {
	userID := c.Params("user")
	if userID == "" {
		return h.respondError(c, http.StatusBadRequest, "user ID is required")
	}

	years, err := h.service.Years(c.Context(), userID)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, "failed get years")
	}

	return c.JSON(years)
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
