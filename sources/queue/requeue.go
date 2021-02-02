package queue

import "sync"

type requeue struct {
	sync.Mutex
	maxRequeue int
	m          map[string]int
}

func newRequeue(maxValue int) *requeue {
	return &requeue{
		Mutex:      sync.Mutex{},
		maxRequeue: maxValue,
		m:          map[string]int{},
	}
}

func (r *requeue) isRequeue(msgId string) bool {
	if r.maxRequeue == 0 {
		return false
	}
	r.Lock()
	defer r.Unlock()
	val, ok := r.m[msgId]
	if ok {
		if val >= r.maxRequeue {
			delete(r.m, msgId)
			return false
		} else {
			r.m[msgId] = val + 1
			return true
		}
	} else {
		r.m[msgId] = 1
		return true
	}
}
func (r *requeue) remove(msgId string) {
	r.Lock()
	defer r.Unlock()
	delete(r.m, msgId)
}
