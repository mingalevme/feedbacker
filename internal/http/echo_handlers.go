package http

import (
	"github.com/labstack/echo/v4"
	"github.com/mingalevme/feedbacker/internal/app/di"
	"github.com/mingalevme/feedbacker/internal/app/interactor"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type EchoHandlerBag struct {
	container di.Container
	services  *interactor.Interactor
}

func NewEchoHandlerBag(container di.Container) *EchoHandlerBag {
	return &EchoHandlerBag{
		container: container,
		services:  interactor.New(container),
	}
}

func (s *EchoHandlerBag) Ping(c echo.Context) error {
	return c.String(http.StatusOK, s.services.Ping())
}

func (s *EchoHandlerBag) LeaveFeedback(c echo.Context) error {
	input := &interactor.LeaveFeedbackData{}
	if err := c.Bind(input); err != nil {
		return err
	}
	f, err := s.services.LeaveFeedback(*input)
	if errors.Is(err, interactor.ErrUnprocessableEntity) {
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
	f, err := s.services.ViewFeedback(id)
	if err == interactor.ErrNotFound {
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
