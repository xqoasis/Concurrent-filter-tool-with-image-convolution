// Used by both balance and stealing
package concurrent

import (
	"math/rand"
	"sync"
	"time"
)

type futureInterface struct {
	task interface{}
	wg *sync.WaitGroup
	result interface{}
}

func randId(capacity int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(capacity)
}

func (future *futureInterface) Get() interface{} {
	future.wg.Wait()
	return future.result
}