package worker

import "log"

// Worker is
type Worker struct {
	Name string
}

// Work is
func (w *Worker) Work(item string) error {
	log.Printf("work %s %s\n", w.Name, item)
	return nil
}
