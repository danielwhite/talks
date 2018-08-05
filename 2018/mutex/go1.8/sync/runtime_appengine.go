// +build appengine

package sync

// SemacquireMutex is like Semacquire, but for profiling contended Mutexes.
// If lifo is true, queue waiter at the head of wait queue.
func runtime_SemacquireMutex(s *uint32, lifo bool) {
	panic("runtime_SemacquireMutex: not supported on GAE")
}

// Semrelease atomically increments *s and notifies a waiting goroutine
// if one is blocked in Semacquire.
// It is intended as a simple wakeup primitive for use by the synchronization
// library and should not be used directly.
// If handoff is true, pass count directly to the first waiter.
func runtime_Semrelease(s *uint32, handoff bool) {
	panic("runtime_Semrelease: not supported on GAE")
}

// Active spinning runtime support.
// runtime_canSpin returns true is spinning makes sense at the moment.
func runtime_canSpin(i int) bool {
	panic("runtime_canSpin: not supported on GAE")
}

// runtime_doSpin does active spinning.
func runtime_doSpin() {
	panic("runtime_doSpin: not supported on GAE")
}
