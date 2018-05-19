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
	ctx        context.Context
	cancel     context.CancelFunc
	queue      *queue.Queue
	numWorkers int
}

// NewRunner is
func NewRunner(numWorkers int) *Runner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Runner{
		queue:      queue.NewQueue(),
		ctx:        ctx,
		cancel:     cancel,
		numWorkers: numWorkers,
	}
}

// Start is
func (r *Runner) Start() {
	r.launchWorkers()
	r.launchServer()
}

// Stop is
func (r *Runner) Stop() {
	r.cancel()
	r.wg.Wait()
	r.queue.Close()
}

func (r *Runner) launchWorkers() {
	for i := 0; i < r.numWorkers; i++ {
		r.wg.Add(1)
		go func(id int) {
			worker := &worker.Worker{Name: fmt.Sprintf("id%02d", id)}
			ch := r.queue.ReadChan()
			for {
				select {
				case item := <-ch:
					worker.Work(item)
				case <-r.ctx.Done():
					log.Printf("worker %s done\n", worker.Name)
					r.wg.Done()
					return
				}
			}
		}(i)
	}
}

func (r *Runner) launchServer() {
	r.wg.Add(1)
	go func() {
		server := &server.Server{Queue: r.queue}
		httpServer := &http.Server{Addr: ":8888", Handler: server}
		go httpServer.ListenAndServe()
		<-r.ctx.Done()
		httpServer.Shutdown(context.Background())
		log.Println("server done")
		r.wg.Done()
	}()
}
