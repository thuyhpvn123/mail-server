package network

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gomail/pkg/logger"
	"gomail/types/network"
)

type Handler struct {
	routes   map[string]func(network.Request) error
	limits   map[string]int // Giới hạn số lần gọi mỗi giây cho từng route
	requests map[string]int // Số lần gọi theo từng route
	mutex    sync.Mutex     // Để tránh điều kiện race
}

func NewHandler(
	routes map[string]func(network.Request) error,
	limits map[string]int, // Giới hạn request cho từng route
) *Handler {
	// Khởi tạo danh sách giới hạn nếu chưa có
	if limits == nil {
		limits = make(map[string]int)
	}
	return &Handler{
		routes:   routes,
		limits:   limits,
		requests: make(map[string]int),
	}
}

func (h *Handler) HandleRequest(r network.Request) error {
	start := time.Now()
	cmd := r.Message().Command()

	logger.Trace("Handling command " + cmd)
	defer func() {
		logger.Trace(
			fmt.Sprintf(
				"Handled command %v from %v took %v",
				cmd,
				r.Connection().Address(),
				time.Since(start).String(),
			),
		)
	}()

	// Kiểm tra giới hạn số lần gọi
	h.mutex.Lock()
	if limit, exists := h.limits[cmd]; exists {
		if h.requests[cmd] >= limit {
			h.mutex.Unlock()
			return errors.New("rate limit exceeded for command: " + cmd)
		}
		h.requests[cmd]++
	}
	h.mutex.Unlock()

	// Reset bộ đếm sau mỗi giây
	go func() {
		time.Sleep(time.Second)
		h.mutex.Lock()
		h.requests[cmd]--
		h.mutex.Unlock()
	}()

	// Xử lý request
	if route, ok := h.routes[cmd]; ok {
		return route(r)
	} else {
		return errors.New("command not found: " + cmd)
	}
}
