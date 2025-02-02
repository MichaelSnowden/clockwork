package clockwork

import (
	"sync"
	"testing"
	"time"
)

// myFunc is an example of a time-dependent function, using an injected clock.
func myFunc(clock Clock, i *int) {
	clock.Sleep(3 * time.Second)
	*i += 1
}

// assertState is an example of a state assertion in a test.
func assertState(t *testing.T, i, j int) {
	if i != j {
		t.Fatalf("i %d, j %d", i, j)
	}
}

// TestMyFunc tests myFunc's behaviour with a FakeClock.
func TestMyFunc(t *testing.T) {
	var i int
	c := NewFakeClock()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		myFunc(c, &i)
		wg.Done()
	}()

	// Wait until myFunc is actually sleeping on the clock.
	c.BlockUntil(1)

	// Assert the initial state.
	assertState(t, i, 0)

	// Now advance the clock forward in time.
	c.Advance(1 * time.Hour)

	// Wait until the function completes.
	wg.Wait()

	// Assert the final state.
	assertState(t, i, 1)
}

// myFunc2 is an example of a time-dependent function which uses AfterFunc.
func myFunc2(clock Clock, i *int) {
	clock.AfterFunc(3*time.Second, func() {
		*i += 1
	})
}

func TestMyFunc2(t *testing.T) {
	var i int

	// Use WithSynchronousAfterFunc to ensure that the AfterFunc callback is called by the time Advance returns.
	c := NewFakeClock(WithSynchronousAfterFunc())

	// Call myFunc2, which will schedule a callback to increment i.
	myFunc2(c, &i)

	// Assert the initial state.
	assertState(t, i, 0)

	// Now advance the clock forward in time to trigger the callback.
	c.Advance(3 * time.Second)

	// Assert the final state.
	assertState(t, i, 1)
}
