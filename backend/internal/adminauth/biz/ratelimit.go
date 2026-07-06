package biz

import (
	"sync"
	"time"
)

// loginLimiter is a small in-memory sliding-window rate limiter for login
// attempts, keyed by normalized email AND by client IP (either key tripping
// blocks the attempt). In-memory is deliberate: adminauth runs as a single
// container and the limiter only needs to blunt online guessing; a distributed
// limiter would be overkill here.
type loginLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	max      int
	attempts map[string][]time.Time
	now      func() time.Time
}

const (
	loginWindow      = time.Minute
	maxLoginAttempts = 5
)

func newLoginLimiter() *loginLimiter {
	return &loginLimiter{
		window:   loginWindow,
		max:      maxLoginAttempts,
		attempts: make(map[string][]time.Time),
		now:      time.Now,
	}
}

// allow records an attempt against every key and reports whether it is within
// the window cap for all of them.
func (l *loginLimiter) allow(keys ...string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)
	allowed := true

	for _, key := range keys {
		if key == "" {
			continue
		}

		kept := l.attempts[key][:0]

		for _, at := range l.attempts[key] {
			if at.After(cutoff) {
				kept = append(kept, at)
			}
		}

		if len(kept) >= l.max {
			allowed = false
		}

		l.attempts[key] = append(kept, now)
	}

	return allowed
}
