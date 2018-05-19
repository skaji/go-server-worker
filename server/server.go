package server

import (
	"fmt"
	"net/http"

	"github.com/skaji/go-server-worker/queue"
)

// Server is
type Server struct {
	Queue *queue.Queue
}

// ServeHTTP is
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if items, ok := r.URL.Query()["item"]; ok {
		for _, item := range items {
			s.Queue.WriteChan() <- item
		}
	}
	fmt.Fprintf(w, "OK\n")
}
