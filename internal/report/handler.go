package report

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

// ErrorResponse представляет стандартный формат ошибки
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse представляет стандартный формат успешного ответа
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *Handler) List(c *fiber.Ctx) error {
	userID := c.Params("user")
	if userID == "" {
		return h.respondError(c, http.StatusBadRequest, "user ID is required")
	}

	month, err := c.ParamsInt("month")
	if err != nil || month < 1 || month > 12 {
		return h.respondError(c, http.StatusBadRequest, "invalid month parameter")
	}

	year, err := c.ParamsInt("year")
	if err != nil || year < 1900 || year > 2100 {
		return h.respondError(c, http.StatusBadRequest, "invalid year parameter")
	}

	prm := repo.GetUserMonthReportParams{
		UserID: userID,
		Month:  int32(month),
		Year:   int32(year),
	}

	report, err := h.service.List(c.Context(), prm)
	if err != nil {
		h.logger.Error("failed to list reports",
			slog.String("user_id", userID),
			slog.Int("month", month),
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to retrieve reports")
	}

	return c.JSON(report)
}

type createRequest struct {
	UserID string  `json:"userId" validate:"required,uuid"`
	Day    int32   `json:"day" validate:"required,min=1,max=31"`
	Month  int32   `json:"month" validate:"required,min=1,max=12"`
	Year   int32   `json:"year" validate:"required,min=1900,max=2100"`
	Hours  float64 `json:"hours" validate:"required,min=0,max=24"`
	Type   string  `json:"typeSystemName" validate:"required"`
}

func (r *createRequest) validate() error {
	if r.UserID == "" {
		return errors.New("userId is required")
	}
	if _, err := uuid.Parse(r.UserID); err != nil {
		return errors.New("userId must be a valid UUID")
	}
	if r.Day < 1 || r.Day > 31 {
		return errors.New("day must be between 1 and 31")
	}
	if r.Month < 1 || r.Month > 12 {
		return errors.New("month must be between 1 and 12")
	}
	if r.Year < 1900 || r.Year > 2100 {
		return errors.New("year must be between 1900 and 2100")
	}
	if r.Hours < 0 || r.Hours > 24 {
		return errors.New("hours must be between 0 and 24")
	}
	if r.Type == "" {
		return errors.New("typeSystemName is required")
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

	report, err := h.service.Create(c.Context(), CreateReportParams{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		Day:    req.Day,
		Month:  req.Month,
		Year:   req.Year,
		Hours:  req.Hours,
		Type:   req.Type,
	})
	if err != nil {
		h.logger.Error("failed to create report",
			slog.String("user_id", req.UserID),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to create report")
	}

	return c.Status(http.StatusCreated).JSON(report)
}

type updateRequest struct {
	ID    string  `json:"id" validate:"required,uuid"`
	Hours float64 `json:"hours" validate:"required,min=0,max=24"`
	Type  string  `json:"typeSystemName" validate:"required"`
}

func (r *updateRequest) validate() error {
	if r.ID == "" {
		return errors.New("id is required")
	}
	if _, err := uuid.Parse(r.ID); err != nil {
		return errors.New("id must be a valid UUID")
	}
	if r.Hours < 0 || r.Hours > 24 {
		return errors.New("hours must be between 0 and 24")
	}
	if r.Type == "" {
		return errors.New("typeSystemName is required")
	}
	return nil
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		return h.respondError(c, http.StatusBadRequest, "invalid request body")
	}

	if err := req.validate(); err != nil {
		return h.respondError(c, http.StatusBadRequest, err.Error())
	}

	report, err := h.service.Update(c.Context(), UpdateReportParams(req))
	if err != nil {
		h.logger.Error("failed to update report",
			slog.String("report_id", req.ID),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to update report")
	}

	return c.JSON(report)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	userID := c.Params("user")
	if userID == "" {
		return h.respondError(c, http.StatusBadRequest, "user ID is required")
	}

	month, err := c.ParamsInt("month")
	if err != nil || month < 1 || month > 12 {
		return h.respondError(c, http.StatusBadRequest, "invalid month parameter")
	}

	year, err := c.ParamsInt("year")
	if err != nil || year < 1900 || year > 2100 {
		return h.respondError(c, http.StatusBadRequest, "invalid year parameter")
	}

	day, err := c.ParamsInt("day")
	if err != nil || day < 1 || day > 31 {
		return h.respondError(c, http.StatusBadRequest, "invalid day parameter")
	}

	err = h.service.Delete(c.Context(), repo.DeleteUserReportParams{
		UserID: userID,
		Day:    int32(day),
		Month:  int32(month),
		Year:   int32(year),
	})
	if err != nil {
		h.logger.Error("failed to delete report",
			slog.String("user_id", userID),
			slog.Int("day", day),
			slog.Int("month", month),
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return h.respondError(c, http.StatusInternalServerError, "failed to delete report")
	}

	return c.JSON(SuccessResponse{
		Message: "Report deleted successfully",
	})
}

// respondError - вспомогательный метод для отправки ошибок
func (h *Handler) respondError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
