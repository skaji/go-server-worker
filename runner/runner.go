package runner

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/skaji/go-server-worker/queue"
	"github.com/skaji/go-server-worker/server"
	"github.com/skaji/go-server-worker/worker"
)

// Runner is
type Runner struct {
	wg         sync.WaitGroup
	queue      *queue.Queue
	numWorkers int
}

// NewRunner is
func NewRunner(numWorkers int) *Runner {
	return &Runner{
		queue:      queue.NewQueue(),
		numWorkers: numWorkers,
	}
}

// Start is
func (r *Runner) Start(ctx context.Context) {
	r.launchWorkers(ctx)
	r.launchServer(ctx)
}

// Wait is
func (r *Runner) Wait() {
	r.wg.Wait()
	r.queue.Close()
}

func (r *Runner) launchWorkers(ctx context.Context) {
	for i := 0; i < r.numWorkers; i++ {
		r.wg.Add(1)
		go func(id int) {
			worker := &worker.Worker{Name: fmt.Sprintf("id%02d", id)}
			ch := r.queue.ReadChan()
			for {
				select {
				case item := <-ch:
					worker.Work(item)
				case <-ctx.Done():
					log.Printf("worker %s done\n", worker.Name)
					r.wg.Done()
					return
				}
			}
		}(i)
	}
}

func (r *Runner) launchServer(ctx context.Context) {
	r.wg.Add(1)
	go func() {
		server := &server.Server{Queue: r.queue}
		httpServer := &http.Server{Addr: ":8888", Handler: server}
		go httpServer.ListenAndServe()
		<-ctx.Done()
		httpServer.Shutdown(context.Background())
		log.Println("server done")
		r.wg.Done()
	}()
}
