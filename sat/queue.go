package sat

type boundedQueue struct {
	elems      []int
	first      int
	last       int
	sumOfQueue int
	maxSize    int
	queueSize  int
}

func newBoundedQueue() *boundedQueue {
	return &boundedQueue{make([]int, 0), 0, 0, 0, 0, 0}
}

func (q *boundedQueue) initSize(size int) {
	q.growTo(size)
}

func (q *boundedQueue) growTo(size int) {
	if len(q.elems) >= size {
		return
	}
	numberNew := size - len(q.elems)
	for range numberNew {
		q.elems = append(q.elems, 0)
	}
	q.first = 0
	q.maxSize = size
	q.queueSize = 0
	q.last = 0
}

func (q *boundedQueue) push(x int) {
	if q.queueSize == q.maxSize {
		q.sumOfQueue -= q.elems[q.last]
		q.last++
		if q.last == q.maxSize {
			q.last = 0
		}
	} else {
		q.queueSize++
	}
	q.sumOfQueue += x
	q.elems[q.first] = x
	q.first++
	if q.first == q.maxSize {
		q.first = 0
		q.last = 0
	}
}

func (q *boundedQueue) avg() int {
	return q.sumOfQueue / int(q.queueSize)
}

func (q *boundedQueue) valid() bool {
	return q.queueSize == q.maxSize
}

func (q *boundedQueue) fastClear() {
	q.first = 0
	q.last = 0
	q.queueSize = 0
	q.sumOfQueue = 0
}
