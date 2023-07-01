package concurrent

import (
	"math/rand"
	"sync"
)

type WBService struct {
	capacity int
	thresholdBalance int
	queues []DEQueue
	shutdown bool
}
func NewWorkBalancingExecutor(capacity, thresholdQueue, thresholdBalance int) ExecutorService {
	var queues []DEQueue
	for i:=0; i < capacity; i++ {
		queues = append(queues, NewUnBoundedDEQueue())
	}
	newWBService := &WBService{
		capacity: capacity, 
		thresholdBalance: thresholdBalance,
		queues: queues,
		shutdown: false,}
	for i := 0; i < capacity; i++ {
		go newWBService.run(i)
	}
	return newWBService
}
// ref: Art of Multiprocessor Programming
func (WBService *WBService) run(id int) {
	for {
		if WBService.queues[id].Size() == 0 && WBService.shutdown {
			return
		}
		//In this project, all tasks (imgTask) are only runnable
		cntTask := WBService.queues[id].PopBottom()
		if cntTask != nil {
			cntFuture, ok := cntTask.(*futureInterface)
			if ok {
				runnable, yes := cntFuture.task.(Runnable)
				if yes {
					runnable.Run()
				}
				cntFuture.wg.Done()
			}else {
				continue
			}
		}
		size := WBService.queues[id].Size()
		// With probability 1/|queue|
		if (rand.Intn(size+1) == size) {
			victim := randId(WBService.capacity) // Choose random victim
			minId, maxId := id, victim
			if id > victim {
				minId, maxId = victim, id
			}
			// lock in canonical order
			WBService.balance(minId, maxId)
		}

	}
}

func (WBService *WBService) balance(minId int, maxId int) {
	qMin, qMax := minId, maxId
	if WBService.queues[minId].Size() > WBService.queues[maxId].Size() {
		qMin, qMax = maxId, minId
	}
	diff := WBService.queues[qMax].Size() - WBService.queues[qMin].Size()
	if (diff > WBService.thresholdBalance) {
		for (WBService.queues[qMax].Size() > WBService.queues[qMin].Size()) {
			WBService.queues[qMin].PushBottom(WBService.queues[qMax].PopTop)
		}
	}
}

func (WBService *WBService) Submit(task interface{}) Future {
	if WBService.shutdown {
		return nil
	}
	receiver := randId(WBService.capacity)
	var wg sync.WaitGroup
	wg.Add(1)
	cntFuture := &futureInterface{task: task, wg: &wg, result: nil}
	WBService.queues[receiver].PushBottom(cntFuture)
	return cntFuture
}
func (WBService *WBService) Shutdown() {
	WBService.shutdown = true
}