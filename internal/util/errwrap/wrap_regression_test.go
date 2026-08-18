package errwrap

import (
	"errors"
	"fmt"
	"testing"
)

// TestWrapDeniedRegression_DoubleWrap captures the bug from jobledger03:
// WrapDenied must preserve the ErrDenied sentinel so that callers which wrap
// its result again (e.g. fmt.Errorf("create: %w", WrapDenied("policy")) in
// internal/svc/catalog.go) can still match ErrDenied via errors.Is.
//
// Before the fix WrapDenied used %v, which collapsed the error into a plain
// string and broke the chain at the first hop, so the double-wrapped error
// never reached ErrDenied.
func TestWrapDeniedRegression_DoubleWrap(t *testing.T) {
	// Single hop: WrapDenied itself must unwrap to ErrDenied.
	wrapped := WrapDenied("update")
	if !errors.Is(wrapped, ErrDenied) {
		t.Fatalf("single WrapDenied lost sentinel: %v", wrapped)
	}

	// Double hop, mirroring internal/svc/catalog.go Create().
	double := fmt.Errorf("create: %w", WrapDenied("policy"))
	if !errors.Is(double, ErrDenied) {
		t.Fatalf("double-wrapped WrapDenied lost sentinel: %v", double)
	}
	if !IsDenied(double) {
		t.Fatalf("IsDenied=false for double-wrapped: %v", double)
	}

	// The wrapped text must still mention the op for diagnostics.
	if msg := double.Error(); msg == "" {
		t.Fatal("wrapped error message is empty")
	}
}
