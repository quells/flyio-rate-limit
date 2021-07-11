package limit

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewMapCounter() *MapCounter {
	return &MapCounter{
		counts: make(map[string]int),
	}
}

type MapCounter struct {
	l      sync.RWMutex
	counts map[string]int
}

func (mc *MapCounter) GetCount(ctx context.Context, ip string) (int, error) {
	if mc == nil {
		return 0, nil
	}

	mc.l.RLock()
	count := mc.counts[ip]
	mc.l.RUnlock()
	log.WithField("ip", ip).Debugf("got count: %v", count)

	return count, nil
}

func (mc *MapCounter) Increment(ctx context.Context, ip string, ttl int) error {
	if mc == nil {
		return fmt.Errorf("counter is nil")
	}

	mc.l.Lock()
	count := mc.counts[ip]
	mc.counts[ip] = count + 1
	mc.l.Unlock()
	log.WithField("ip", ip).Debugf("incremented to: %v", count+1)

	if ttl > 0 {
		go mc.scheduleDecrement(ip, ttl)
	}

	return nil
}

func (mc *MapCounter) scheduleDecrement(ip string, ttl int) {
	timeout := time.Duration(ttl) * time.Second
	<-time.NewTimer(timeout).C

	mc.l.Lock()
	count := mc.counts[ip]
	if count > 0 {
		mc.counts[ip] = count - 1
	}
	mc.l.Unlock()
	log.WithField("ip", ip).Debug("decremented")
}
