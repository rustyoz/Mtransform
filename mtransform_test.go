package mtransform

import (
	"testing"
)

func TestMultiply(t *testing.T) {
	a := Transform{{1, 2, 3}, {-1, -2, -3}, {4, 5, 6}}
	b := Transform{{0, 1, 2}, {0, -1, -2}, {2, 1, 0}}
	want := Transform{{6, 2, -2}, {-6, -2, 2}, {12, 5, -2}}
	got := MultiplyTransforms(a, b)
	if got != want {
		t.Errorf("Multiplying: got %v, want %v", got, want)
	}
}
