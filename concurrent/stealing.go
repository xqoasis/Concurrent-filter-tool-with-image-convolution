package concurrent

import (
	"sync"
)

type WSService struct {
	capacity int
	queues []DEQueue
	shutdown bool
}
func NewWorkStealingExecutor(capacity, threshold int) ExecutorService {
	var queues []DEQueue
	for i := 0; i < capacity; i++ {
		queues = append(queues, NewUnBoundedDEQueue())
	}
	newWSService := &WSService{
		capacity: capacity,
		queues: queues,
		shutdown: false,
	}
	for i := 0; i < capacity; i++ {
		go newWSService.run(i)
	}
	return newWSService
}

func (WSService *WSService) run(id int) {
	cntTask := WSService.queues[id].PopBottom()
	for {
		for cntTask != nil {
			cntFuture, ok := cntTask.(*futureInterface)
			if ok {
				//In this project, all tasks (imgTask) are runnable instead of callable
				runnable, yes := cntFuture.task.(Runnable)
				if yes {
					runnable.Run()
				}
				cntFuture.wg.Done()
				cntTask = WSService.queues[id].PopBottom()
			} else {
				continue
			}
		}

		for cntTask == nil {
			if WSService.queues[id].Size() == 0 && WSService.shutdown {
				return
			}
			// steal other queue's task randomly
			victim := randId(WSService.capacity)
			if !WSService.queues[victim].IsEmpty() {
				cntTask = WSService.queues[victim].PopTop()
			}
		}
	}
}

func (WSService * WSService) Submit(task interface{}) Future {
	if WSService.shutdown {
		return nil
	}
	receiver := randId(WSService.capacity)
	var wg sync.WaitGroup
	wg.Add(1)
	cntFuture := &futureInterface{task: task, wg: &wg, result: nil}
	WSService.queues[receiver].PushBottom(cntFuture)
	return cntFuture
}

func (WSService *WSService) Shutdown() {
	WSService.shutdown = true
}

