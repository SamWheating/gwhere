package main

import (
	"testing"
)

func TestPlaceholder(t *testing.T) {
	got := 2 - 1
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}
