package network

import (
	"errors"
	"fmt"
	"time"

	"gomail/mtn/logger"
	"gomail/mtn/types/network"
)

type Handler struct {
	routes map[string]func(network.Request) error
}

func NewHandler(
	routes map[string]func(network.Request) error,
) *Handler {
	return &Handler{
		routes,
	}
}

func (h *Handler) HandleRequest(r network.Request) error {
	start := time.Now()
	logger.Trace("Handling command " + r.Message().Command())
	defer func() {
		logger.Trace(
			fmt.Sprintf(
				"Handled command %v from %v took %v",
				r.Message().Command(),
				r.Connection().Address(),
				time.Since(start).String(),
			),
		)
	}()
	if route, ok := h.routes[r.Message().Command()]; ok {
		return route(r)
	} else {
		return errors.New("command not found: " + r.Message().Command())
	}
}
