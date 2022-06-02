package workqueue

// ref: https://github.com/nayuta87/queue/blob/master/queue.go
import (
	"proj3/func_objs"
	"sync/atomic"
	"sync"
	"unsafe"
)

type Task struct {
	Func func_objs.Runnable
	Arg  interface{}
}

// Queue is a lock-free unbounded queue.
type Queue struct {
	head unsafe.Pointer // *node
	tail unsafe.Pointer // *node
	size int32
	mutex *sync.Mutex
}

type node struct {
	value Task
	next  unsafe.Pointer // *node
}

// NewQueue returns a pointer to an empty queue.
func NewQueue() (q *Queue) {
	n := unsafe.Pointer(&node{})
	var mutex sync.Mutex
	q = &Queue{head: n, tail: n, mutex: &mutex}
	return
}

// Enq puts the given value v at the tail of the queue.
func (q *Queue) Enq(task func_objs.Runnable, taskArg interface{}) {
	newTask := Task{Func: task, Arg: taskArg}
	n := &node{value: newTask}
	for {
		last := load(&q.tail)
		next := load(&last.next)
		if last == load(&q.tail) {
			if next == nil {
				if cas(&last.next, next, n) {
					cas(&q.tail, last, n)
					atomic.AddInt32(&q.size, 1)
					return
				}
			} else {
				cas(&q.tail, last, next)
				atomic.AddInt32(&q.size, 1)
			}
		}
	}
}

// Deq removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *Queue) Deq() *Task {
	for {
		first := load(&q.head)
		last := load(&q.tail)
		next := load(&first.next)
		if first == load(&q.head) {
			if first == last {
				if next == nil {
					return nil
				}
				cas(&q.tail, last, next)
				
			} else {
				v := next.value
				if cas(&q.head, first, next) {
					atomic.AddInt32(&q.size, -1)
					return &v
				}
			}
		}
	}
}

func (q *Queue) Size() int {
	return int(q.size)
}

func (q *Queue) Lock() {
	q.mutex.Lock()
}

func (q *Queue) UnLock() {
	q.mutex.Unlock()
}

func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

func store(p *unsafe.Pointer, n *node) {
	atomic.StorePointer(p, unsafe.Pointer(n))
}

func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}