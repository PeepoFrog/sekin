package subscriptionmanager

import (
	"sync"

	"github.com/kiracore/sekin/src/shidai/internal/logger"
	"go.uber.org/zap"
)

type SubscriptionManager struct {
	subscribers map[chan bool]struct{}
	lock        sync.Mutex
	log         *zap.Logger
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[chan bool]struct{}),
		log:         logger.GetLogger(),
	}
}

func (m *SubscriptionManager) Subscribe(c chan bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.subscribers[c] = struct{}{}
	m.log.Debug("Subscriber added", zap.Int("current_subscribers", len(m.subscribers)))
}

func (m *SubscriptionManager) Unsubscribe(c chan bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exists := m.subscribers[c]; exists {
		delete(m.subscribers, c)
		m.log.Debug("Subscriber removed", zap.Int("current_subscribers", len(m.subscribers)))
	}
}

func (m *SubscriptionManager) NotifySubscribers() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.log.Debug("Notifying subscribers", zap.Int("count", len(m.subscribers)))
	for c := range m.subscribers {
		select {
		case c <- true:
			m.log.Debug("Notification sent to a subscriber")
		default:
			m.log.Debug("Skipped a subscriber, channel was full")
		}
	}
}

func NewSubscriberChannel() chan bool {
	return make(chan bool, 1) // Buffer size of 1 to prevent blocking on first send
}
