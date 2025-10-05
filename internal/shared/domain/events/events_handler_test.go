package shared_events

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Mock Event para testes
type mockEvent struct {
	name string
}

func (m *mockEvent) EventName() string {
	return m.name
}

func TestNewEventDispatcher(t *testing.T) {
	dispatcher := NewEventDispatcher()

	if dispatcher == nil {
		t.Fatal("NewEventDispatcher() returned nil")
	}

	if dispatcher.handlers == nil {
		t.Error("handlers map is nil")
	}

	if len(dispatcher.handlers) != 0 {
		t.Errorf("handlers map not empty, got %d items", len(dispatcher.handlers))
	}
}

func TestEventDispatcher_Register(t *testing.T) {
	dispatcher := NewEventDispatcher()

	handler1 := func(event Event) {}
	handler2 := func(event Event) {}

	t.Run("register single handler", func(t *testing.T) {
		dispatcher.Register("test.event", handler1)

		dispatcher.mu.RLock()
		defer dispatcher.mu.RUnlock()

		if len(dispatcher.handlers["test.event"]) != 1 {
			t.Errorf("expected 1 handler, got %d", len(dispatcher.handlers["test.event"]))
		}
	})

	t.Run("register multiple handlers for same event", func(t *testing.T) {
		dispatcher.Register("test.event", handler2)

		dispatcher.mu.RLock()
		defer dispatcher.mu.RUnlock()

		if len(dispatcher.handlers["test.event"]) != 2 {
			t.Errorf("expected 2 handlers, got %d", len(dispatcher.handlers["test.event"]))
		}
	})

	t.Run("register handler for different event", func(t *testing.T) {
		dispatcher.Register("another.event", handler1)

		dispatcher.mu.RLock()
		defer dispatcher.mu.RUnlock()

		if len(dispatcher.handlers["another.event"]) != 1 {
			t.Errorf("expected 1 handler for another.event, got %d", len(dispatcher.handlers["another.event"]))
		}
	})
}

func TestEventDispatcher_Dispatch(t *testing.T) {
	t.Run("dispatch with registered handler", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		event := &mockEvent{name: "test.event"}

		var called atomic.Bool
		var wg sync.WaitGroup
		wg.Add(1)

		handler := func(e Event) {
			defer wg.Done()
			called.Store(true)

			if e.EventName() != "test.event" {
				t.Errorf("handler received wrong event: %v", e.EventName())
			}
		}

		dispatcher.Register("test.event", handler)
		dispatcher.Dispatch("test.event", event)

		// Aguardar handler ser chamado (goroutine)
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			if !called.Load() {
				t.Error("handler was not called")
			}
		case <-time.After(1 * time.Second):
			t.Error("handler took too long to execute")
		}
	})

	t.Run("dispatch without registered handler", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		event := &mockEvent{name: "unregistered.event"}

		// Deve apenas printar mensagem, não causar panic
		dispatcher.Dispatch("unregistered.event", event)
	})

	t.Run("dispatch with multiple handlers", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		event := &mockEvent{name: "multi.event"}

		var counter atomic.Int32
		var wg sync.WaitGroup
		wg.Add(3)

		handler := func(e Event) {
			defer wg.Done()
			counter.Add(1)
		}

		// Registrar 3 handlers
		dispatcher.Register("multi.event", handler)
		dispatcher.Register("multi.event", handler)
		dispatcher.Register("multi.event", handler)

		dispatcher.Dispatch("multi.event", event)

		// Aguardar todos os handlers
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			if counter.Load() != 3 {
				t.Errorf("expected 3 handlers to be called, got %d", counter.Load())
			}
		case <-time.After(1 * time.Second):
			t.Error("handlers took too long to execute")
		}
	})
}

func TestEventDispatcher_ConcurrentRegisterAndDispatch(t *testing.T) {
	dispatcher := NewEventDispatcher()
	var wg sync.WaitGroup

	// Registrar handlers concorrentemente
	numGoroutines := 50
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			handler := func(e Event) {}
			dispatcher.Register("concurrent.event", handler)
		}(i)
	}

	wg.Wait()

	// Verificar que não houve race conditions
	dispatcher.mu.RLock()
	handlersCount := len(dispatcher.handlers["concurrent.event"])
	dispatcher.mu.RUnlock()

	if handlersCount != numGoroutines {
		t.Errorf("expected %d handlers, got %d", numGoroutines, handlersCount)
	}
}

func TestEventDispatcher_ConcurrentDispatch(t *testing.T) {
	dispatcher := NewEventDispatcher()
	var counter atomic.Int32

	handler := func(e Event) {
		counter.Add(1)
	}

	dispatcher.Register("dispatch.event", handler)

	// Disparar eventos concorrentemente
	numDispatches := 100
	var wg sync.WaitGroup
	wg.Add(numDispatches)

	for i := 0; i < numDispatches; i++ {
		go func() {
			defer wg.Done()
			event := &mockEvent{name: "dispatch.event"}
			dispatcher.Dispatch("dispatch.event", event)
		}()
	}

	wg.Wait()

	// Aguardar um pouco para handlers assíncronos terminarem
	time.Sleep(100 * time.Millisecond)

	if counter.Load() != int32(numDispatches) {
		t.Errorf("expected %d handler calls, got %d", numDispatches, counter.Load())
	}
}

// Benchmark
func BenchmarkEventDispatcher_Register(b *testing.B) {
	dispatcher := NewEventDispatcher()
	handler := func(e Event) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.Register("bench.event", handler)
	}
}

func BenchmarkEventDispatcher_Dispatch(b *testing.B) {
	dispatcher := NewEventDispatcher()
	event := &mockEvent{name: "bench.event"}

	handler := func(e Event) {}
	dispatcher.Register("bench.event", handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.Dispatch("bench.event", event)
	}
}

func BenchmarkEventDispatcher_RegisterAndDispatch(b *testing.B) {
	event := &mockEvent{name: "bench.event"}
	handler := func(e Event) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher := NewEventDispatcher()
		dispatcher.Register("bench.event", handler)
		dispatcher.Dispatch("bench.event", event)
	}
}
