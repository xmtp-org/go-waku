package relay

import (
	"sync"

	"github.com/waku-org/go-waku/waku/v2/protocol"
)

// Subscription handles the subscrition to a particular pubsub topic
type Subscription struct {
	sync.RWMutex

	// C is channel used for receiving envelopes
	C chan *protocol.Envelope

	closed bool
	once   sync.Once
	quit   chan struct{}
	wg     sync.WaitGroup
}

// Unsubscribe will close a subscription from a pubsub topic. Will close the message channel
func (subs *Subscription) Unsubscribe() {
	subs.once.Do(func() {
		close(subs.quit)
	})
	subs.wg.Wait()
}

// IsClosed determine whether a Subscription is still open for receiving messages
func (subs *Subscription) IsClosed() bool {
	subs.RLock()
	defer subs.RUnlock()
	return subs.closed
}
