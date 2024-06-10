package Wait

import (
	"bufio"
	"os"
	"slices"
	"sync"
)

var (
	// observers is the list of observers.
	observers []*Waitee

	// lastID is the last ID.
	lastID int

	// wg is the wait group.
	wg sync.WaitGroup

	// closeChan is the channel to close the pool.
	closeChan chan struct{}

	// mu is the mutex to synchronize access to the pool.
	mu sync.RWMutex
)

func init() {
	observers = make([]*Waitee, 0)
	lastID = 0
	wg = sync.WaitGroup{}
	closeChan = nil
	mu = sync.RWMutex{}

	// Start the pool
	Start()
}

// notifyAll is a helper function that notifies all observers.
//
// Parameters:
//   - isKeyPressed: A flag indicating whether a key was pressed.
func notifyAll(isKeyPressed bool) {
	mu.RLock()

	var wg sync.WaitGroup

	wg.Add(len(observers))

	for _, observer := range observers {
		go func(o *Waitee) {
			defer wg.Done()

			o.Notify(isKeyPressed)
		}(observer)
	}

	mu.RUnlock()

	wg.Wait()
}

// Start starts the pool wait.
func Start() {
	if closeChan != nil {
		return
	}

	closeChan = make(chan struct{})

	wg.Add(1)

	go enterKeyListener()
}

// Close closes the pool wait.
func Close() {
	close(closeChan)

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	for i := range observers {
		observers[i] = nil
	}

	observers = nil
	closeChan = nil
}

// enterKeyListener is a helper function that listens for a key press.
func enterKeyListener() {
	defer wg.Done()

	for {
		select {
		case <-closeChan:
			notifyAll(false)
			return
		default:
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')

			notifyAll(true)
		}
	}
}

// GetWaitee gets a new waitee.
//
// Returns:
//   - *Waitee: The new waitee.
func GetWaitee() *Waitee {
	mu.Lock()
	defer mu.Unlock()

	w := &Waitee{
		cond:         sync.NewCond(&sync.Mutex{}),
		isKeyPressed: false,
		isClosed:     false,
		id:           lastID,
	}

	lastID++
	observers = append(observers, w)

	return w
}

// IsRunning returns true if the pool wait is running, false otherwise.
//
// Returns:
//   - bool: True if the pool wait is running, false otherwise.
func IsRunning() bool {
	return closeChan != nil
}

// removeWaitee removes the waitee from the pool.
//
// Parameters:
//   - w: The waitee to remove.
func removeWaitee(w *Waitee) {
	if w == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	observers[w.id] = nil
	observers = slices.Delete(observers, w.id, w.id+1)
}
