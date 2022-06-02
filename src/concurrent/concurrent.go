package concurrent

import (
	"sync/atomic"
	"unsafe"
	"proj3/func_objs"
)


// wrap up Func and Arg, easier to store in the Deque
type Task struct {
	Func func_objs.Runnable
	Arg  interface{}
}

type node struct {
	value Task
	prev unsafe.Pointer
	next unsafe.Pointer
}

type Dequeue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	len  int32
}

type DEQueue interface {
	// not like project 1 only has one producer pushing task in queue
	// here need to worry about the concurrency when pushing task to Dequeue
	PushBottom(task func_objs.Runnable, taskArg interface{})
	IsEmpty() bool
	PopTop() *Task
	PopBottom() *Task
}

func NewDEQueue() (deq *Dequeue) {
	n := unsafe.Pointer(&node{})
	deq = &Dequeue{head: n, tail: n}
	return 
}

func (deq *Dequeue) PushBottom(task func_objs.Runnable, taskArg interface{}) {
	newTask := Task{Func: task, Arg: taskArg}
	n := &node{value: newTask}
	nullNode := &node{}

	oldTail := load(&deq.tail)
	prev := load(&oldTail.prev)
	next := load(&oldTail.next)

	if oldTail == load(&deq.tail) {
		if next == nil {
			// safely add the new node 
			if prev != nil {
				// update all nodes prev/next pointer
				store(&prev.next, n)
				store(&n.prev, prev)
				store(&n.next, nullNode)
				store(&nullNode.prev, n)
			} else {
				store(&n.next, nullNode)
				store(&nullNode.prev, n)
				cas(&deq.head, oldTail, n)
			}				
			// update len atomically
			atomic.AddInt32(&deq.len, 1)
			// update the deq tail pointer
			cas(&deq.tail, oldTail, nullNode)
		}
	}
	//fmt.Println("After Push Bottom: ", load(&deq.head).value)
}

func (deq *Dequeue) IsEmpty() bool {
	return deq.len == 0
}

func (deq *Dequeue) PopTop() *Task {
	if deq.len <= 0 {
		return nil
	}

	// keep failed thread to keep trying to PopTop()
	for {
		oldHead := load(&deq.head)
		newHead := load(&oldHead.next)

		task := oldHead.value

		if cas(&deq.head, oldHead, newHead) {
			atomic.AddInt32(&deq.len, -1)
			return &task
		} else {
			return nil
		}
	}
	
}

func (deq *Dequeue) PopBottom() *Task {
	// no nodes in dequeue, do nothing
	if atomic.CompareAndSwapInt32(&deq.len, 0, 0) {
		return nil
	}
	
	
	oldTail := load(&deq.tail)
	newTail := load(&oldTail.prev)
	oldHead := load(&deq.head)
	
	task := newTail.value
	// competing with PopTop()
	if atomic.CompareAndSwapInt32(&deq.len, 1, 1) {
		if cas(&deq.tail, oldTail, newTail) {
			atomic.AddInt32(&deq.len, -1)
			// now head and tail points to the same node
			return &task
		}
	}
	
	// safely PopBottom()
	if atomic.CompareAndSwapInt32(&deq.len, 1, deq.len-1) {
		// update next/prev pointers to connect with oldHead and oldTail
		store(&oldHead.next, oldTail)
		store(&oldTail.prev, oldHead)
		return &task
	}
	return nil
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

func (deq *Dequeue) GetLen() int32 {
	return deq.len
}


