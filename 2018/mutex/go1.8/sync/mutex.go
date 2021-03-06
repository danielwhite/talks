// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sync provides basic synchronization primitives such as mutual
// exclusion locks. Other than the Once and WaitGroup types, most are intended
// for use by low-level library routines. Higher-level synchronization is
// better done via channels and communication.
//
// Values containing the types defined in this package should not be copied.
package sync

import (
	"sync/atomic"
)

// A Mutex is a mutual exclusion lock.
// Mutexes can be created as part of other structures;
// the zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.\
// 1MUX OMIT
type Mutex struct {
	state int32
	sema  uint32 // wait queue
}

// 2MUX OMIT

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexWaiterShift = iota
)

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
	// 1CAS OMIT
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}
	// 2CAS OMIT

	awoke := false
	iter := 0
	for {
		old := m.state
		new := old | mutexLocked
		if old&mutexLocked != 0 {
			if runtime_canSpin(iter) {
				// Active spinning makes sense.
				// Try to set mutexWoken flag to inform Unlock
				// to not wake other blocked goroutines.
				if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
					atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
					awoke = true
				}
				runtime_doSpin()
				iter++
				continue
			}
			new = old + 1<<mutexWaiterShift
		}
		if awoke {
			// The goroutine has been woken from sleep,
			// so we need to reset the flag in either case.
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken
		}
		// 1SEMA OMIT
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&mutexLocked == 0 {
				break
			}
			runtime_SemacquireMutex(&m.sema, false) // HL
			awoke = true
			iter = 0
		}
		// 2SEMA OMIT
	}
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	// 1DROP OMIT
	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked)
	// 2DROP OMIT
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}

	old := new
	for {
		// 1HOFF OMIT
		// If there are no waiters or a goroutine has already
		// been woken or grabbed the lock, no need to wake anyone.
		if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0 {
			return
		}
		// 2HOFF OMIT
		// Grab the right to wake someone.
		new = (old - 1<<mutexWaiterShift) | mutexWoken
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			runtime_Semrelease(&m.sema, false) // HL
			return
		}
		// 3HOFF OMIT
		old = m.state
	}
}
