package queue

// Queue is
type Queue struct {
	rChan chan string
	wChan chan string
}

// NewQueue is
func NewQueue() *Queue {
	ch := make(chan string, 1000)
	return &Queue{rChan: ch, wChan: ch}
}

// WriteChan is
func (q *Queue) WriteChan() chan string {
	return q.wChan
}

// ReadChan is
func (q *Queue) ReadChan() <-chan string {
	return q.rChan
}

// Close is
func (q *Queue) Close() {
	if q.wChan != q.rChan {
		close(q.wChan)
	}
	close(q.rChan)
}
