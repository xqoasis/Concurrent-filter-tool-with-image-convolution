package concurrent

import (
	"sync"
	"sync/atomic"
)

type Task interface{}

type DEQueue interface {
	PushBottom(task Task)
	IsEmpty() bool //returns whether the queue is empty
	PopTop() Task
	PopBottom() Task
	Size()	int // added api to return the size of dequeue from other files
}

type Node struct {
	task 	Task
	prev 	*Node
	next 	*Node
}

type UnBoundedDEQueue struct {
	mutex 	*sync.Mutex
	top 	*Node
	bottom 	*Node
	size 	int64
}


// NewUnBoundedDEQueue returns an empty UnBoundedDEQueue
func NewUnBoundedDEQueue() DEQueue {
	return &UnBoundedDEQueue{mutex: &sync.Mutex{}, top: nil, bottom: nil, size: 0, }
}

func (queue *UnBoundedDEQueue) PushBottom(task Task) {
	queue.mutex.Lock() 
	defer queue.mutex.Unlock()

	newNode := &Node{task: task, prev: nil, next: nil}
	if queue.bottom == nil {
		// both top and bottom are the new node
		queue.bottom = newNode
		queue.top = newNode
	} else {
		newNode.prev = queue.bottom
		queue.bottom.next = newNode
		queue.bottom = newNode
	}
	queue.size ++
}

func (queue *UnBoundedDEQueue) IsEmpty() bool{
	return (queue.bottom == nil || queue.top == nil || queue.Size() == 0)
}

func (queue *UnBoundedDEQueue) PopTop() Task{
	queue.mutex.Lock() 
	defer queue.mutex.Unlock()

	if (queue.IsEmpty()) {
		//check again
		return nil  
	} 
	queue.size --
	temp := queue.top
	queue.top = queue.top.next
	
	if (queue.top == nil) {
		queue.bottom = nil
	} else {
		queue.top.prev = nil
	}
	return temp.task
	
}

func (queue *UnBoundedDEQueue) PopBottom() Task {
	queue.mutex.Lock() 
	defer queue.mutex.Unlock()

	if (queue.IsEmpty()) {
		return nil
	} 
	queue.size --
	temp := queue.bottom
	queue.bottom = queue.bottom.prev
	if (queue.bottom == nil) {
		queue.top = nil
	} else {
		queue.bottom.next = nil
	}
	return temp.task
	
}

func (queue *UnBoundedDEQueue) Size() int{
	return int(atomic.LoadInt64(&queue.size))
}
