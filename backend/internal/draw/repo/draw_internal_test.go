package repo

import (
	"application/internal/draw/biz"
	"errors"
	"testing"
)

// randomIndex uses crypto/rand; verify it stays in bounds, covers the range,
// and rejects a non-positive n.
func TestRandomIndex(t *testing.T) {
	const n = 8

	seen := make(map[int]bool)

	for range 500 {
		idx, err := randomIndex(n)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if idx < 0 || idx >= n {
			t.Fatalf("index %d out of range [0,%d)", idx, n)
		}

		seen[idx] = true
	}

	if len(seen) != n {
		t.Fatalf("expected all %d indices to appear, saw %d", n, len(seen))
	}
}

func TestRandomIndex_ZeroReturnsNoTickets(t *testing.T) {
	if _, err := randomIndex(0); !errors.Is(err, biz.ErrNoTickets) {
		t.Fatalf("expected ErrNoTickets, got %v", err)
	}
}
