package standard

import (
	repo "TimeTrack/internal/adapter/mysql/sqlc"
	"errors"
	"log/slog"
	"net/http"

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

func (h *Handler) ListForSetting(c *fiber.Ctx) error {
	year, err := c.ParamsInt("year")
	if err != nil || year < 1900 || year > 2100 {
		return h.respondError(c, http.StatusBadRequest, "invalid year parameter")
	}

	standards, err := h.service.ListForSetting(c.Context(), int32(year))
	if err != nil {
		h.logger.Error("failed to get vacations",
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to get vacations")
	}

	return c.JSON(standards)
}

type createRequest struct {
	Month    int32 `json:"month"`
	Year     int32 `json:"year"`
	Hours    int32 `json:"hours"`
	GenderID int32 `json:"genderId"`
}

func (r *createRequest) validate() error {
	if r.Hours < 0 {
		return errors.New("hours must be > 0")
	}
	if r.Month < 1 || r.Month > 12 {
		return errors.New("month must be between 1 and 12")
	}
	if r.Year < 1900 || r.Year > 2100 {
		return errors.New("year must be between 1900 and 2100")
	}
	return nil
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		return h.respondError(c, http.StatusBadRequest, "invalid request body")
	}

	if err := req.validate(); err != nil {
		return h.respondError(c, http.StatusBadRequest, err.Error())
	}

	report, err := h.service.Create(c.Context(), repo.CreateStandardParams{
		ID:       uuid.NewString(),
		Month:    req.Month,
		Year:     req.Year,
		Hours:    req.Hours,
		GenderID: req.GenderID,
	})
	if err != nil {
		h.logger.Error("failed to create report",
			slog.Int64("Month", int64(req.Month)),
			slog.Int64("Year", int64(req.Year)),
			slog.Int64("GenderID", int64(req.GenderID)),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to create report")
	}

	return c.Status(http.StatusCreated).JSON(report)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var req repo.UpdateStandardParams
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		return h.respondError(c, http.StatusBadRequest, "invalid request body")
	}

	if req.ID == "" {
		return errors.New("id is required")
	}
	if _, err := uuid.Parse(req.ID); err != nil {
		return errors.New("id must be a valid UUID")
	}
	if req.Hours < 0 {
		return errors.New("hours no valid")
	}

	err := h.service.Update(c.Context(), req)
	if err != nil {
		h.logger.Error("failed to update report",
			slog.String("report_id", req.ID),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to update report")
	}

	c.Status(http.StatusOK)
	return nil
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
