package Wait

import "sync"

// Waitee is a type that waits for a key press.
type Waitee struct {
	// id is the ID of the waitee.
	id int

	// cond is the condition variable.
	cond *sync.Cond

	// isKeyPressed is a flag indicating whether a key was pressed.
	isKeyPressed bool

	// isClosed is a flag indicating whether the waitee is closed.
	isClosed bool
}

// Notify notifies the waitee.
//
// Parameters:
//   - isKeyPressed: A flag indicating whether a key was pressed.
func (w *Waitee) Notify(isKeyPressed bool) {
	w.cond.L.Lock()
	defer w.cond.L.Unlock()

	if isKeyPressed {
		w.isKeyPressed = true
	} else {
		w.isClosed = true
	}

	w.cond.Broadcast()
}

// Wait waits for a key press.
func (w *Waitee) Wait() {
	if w.isClosed {
		// No need to wait
		return
	}

	w.cond.L.Lock()

	for !w.isKeyPressed && !w.isClosed {
		w.cond.Wait()
	}

	if w.isClosed || w.isKeyPressed {
		// Reset the state
		w.isKeyPressed = false
	}

	w.cond.L.Unlock()
}

func (w *Waitee) Clean() {
	// remove from the pool
	removeWaitee(w)
}
