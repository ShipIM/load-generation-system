package attack

import (
	"load-generation-system/internal/core"
	"load-generation-system/pkg/broadcast"
	"sync"
	"time"
)

// attackService implements the core.AttackService interface and manages the
// complete lifecycle of load test attacks across distributed nodes.
//
// Fields:
//   - nodes: Map of registered nodes by their names
//   - removingCancels: Channels for cancellation signals during node removal
//   - attacks: Active attacks indexed by attack ID
//   - attackSeq: Sequence counter for generating unique attack IDs
//   - incrementSeqs: Sequence counters for generating increment IDs per attack
//   - recoveryInterval: Duration between recovery attempts for failed operations
//   - mu: Read-write mutex for concurrent access protection
type attackService struct {
	nodes            map[string]core.Node // Active worker nodes
	removingCancels  map[string]chan any  // Node removal cancellation channels
	attacks          map[int64]attack     // Active attacks
	attackSeq        int64                // Attack ID sequence counter
	incrementSeqs    map[int64]int64      // Increment ID sequences per attack
	recoveryInterval time.Duration        // Recovery retry interval
	mu               sync.RWMutex         // Concurrency control
}

// attack represents a single load test attack with its configuration and control mechanisms.
//
// Fields:
//   - details: Configuration and metadata of the attack
//   - stopBr: Broadcast channel for stopping the attack across all nodes
type attack struct {
	details core.AttackDetails          // Attack parameters and state
	stopBr  *broadcast.Broadcaster[any] // Attack stop signal broadcaster
}

func NewService(recoveryIntervalSec int64) core.AttackService {
	return &attackService{
		nodes:            make(map[string]core.Node),
		removingCancels:  make(map[string]chan any),
		attacks:          make(map[int64]attack),
		incrementSeqs:    make(map[int64]int64),
		recoveryInterval: time.Duration(recoveryIntervalSec) * time.Second,
	}
}
