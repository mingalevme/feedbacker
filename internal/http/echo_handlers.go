package http

import (
	"github.com/labstack/echo/v4"
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/internal/app/interactor/health"
	"github.com/mingalevme/feedbacker/internal/app/interactor/leave_feedback"
	"github.com/mingalevme/feedbacker/internal/app/interactor/ping"
	"github.com/mingalevme/feedbacker/internal/app/interactor/view_feedback"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type EchoHandlerBag struct {
	env           app.Env
	ping          *ping.Ping
	health        *health.Health
	leaveFeedback *leave_feedback.LeaveFeedback
	viewFeedback  *view_feedback.ViewFeedback
}

func NewEchoHandlerBag(env app.Env) *EchoHandlerBag {
	return &EchoHandlerBag{
		env:           env,
		ping:          ping.New(env),
		health:        health.New(env),
		leaveFeedback: leave_feedback.New(env),
		viewFeedback:  view_feedback.New(env),
	}
}

func (s *EchoHandlerBag) Ping(c echo.Context) error {
	return c.String(http.StatusOK, s.ping.Ping())
}

func (s *EchoHandlerBag) Health(c echo.Context) error {
	// Content-Type: application/health+json
	h := s.health.Health()
	if h.Status != health.HealthStatusPass {
		//return c.JSON(http.StatusFailedDependency, h)
		return c.JSON(http.StatusBadGateway, h)
	}
	return c.JSON(http.StatusOK, h)
}

func (s *EchoHandlerBag) LeaveFeedback(c echo.Context) error {
	input := &leave_feedback.LeaveFeedbackData{}
	if err := c.Bind(input); err != nil {
		return err
	}
	f, err := s.leaveFeedback.LeaveFeedback(*input)
	if errors.Is(err, leave_feedback.ErrUnprocessableEntity) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if err != nil {
		return err
	}
	if err := c.JSON(http.StatusCreated, f); err != nil {
		return err
	}
	return nil
}

func (s *EchoHandlerBag) ViewFeedback(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return echo.ErrNotFound
	}
	f, err := s.viewFeedback.ViewFeedback(id)
	if err == view_feedback.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	if err := c.JSON(http.StatusCreated, f); err != nil {
		return err
	}
	return nil
}
