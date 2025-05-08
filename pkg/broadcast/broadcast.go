package broadcast

import (
	"sync"
)

// Broadcaster provides a thread-safe mechanism to broadcast messages to multiple subscribers.
// It is generic over the message type T, allowing broadcasting of any data type.
type Broadcaster[T any] struct {
	subscribers map[chan T]any // Map of subscriber channels (using empty interface as value)
	mu          sync.RWMutex   // Mutex to protect concurrent access to subscribers
}

const (
	// defaultBroadcastCapacity is the buffer size for each subscriber's channel
	defaultBroadcastCapacity = 10
)

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		subscribers: make(map[chan T]any),
	}
}

// Subscribe adds a new subscriber to the broadcaster and returns a channel
// that will receive all broadcast messages. The channel has a buffer size
// defined by defaultBroadcastCapacity to prevent blocking the broadcaster.
//
// Returns:
//   - chan T: The channel that will receive broadcast messages
func (b *Broadcaster[T]) Subscribe() chan T {
	ch := make(chan T, defaultBroadcastCapacity)
	b.mu.Lock()
	b.subscribers[ch] = nil
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes a subscriber channel from the broadcaster and closes
// the channel. This is safe to call multiple times as channel closing is
// idempotent.
//
// Parameters:
//   - ch: The channel to unsubscribe and close
func (b *Broadcaster[T]) Unsubscribe(ch chan T) {
	b.mu.Lock()
	delete(b.subscribers, ch)
	close(ch)
	b.mu.Unlock()
}

// Broadcast sends a message to all current subscribers. The operation is
// non-blocking - if a subscriber's channel is full, the message will be
// skipped for that subscriber.
//
// Parameters:
//   - value: The message to broadcast to all subscribers
func (b *Broadcaster[T]) Broadcast(value T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subscribers {
		select {
		case ch <- value: // Try to send if channel has capacity
		default: // Skip if channel is full
		}
	}
}
