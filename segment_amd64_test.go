//go:build amd64

package obreron_test

import (
	"testing"
	"unsafe"

	"github.com/profe-ajedrez/obreron/v3"
)

func TestSegmentSizeAMD64(t *testing.T) {
	if got, want := unsafe.Sizeof(obreron.NewEmptySegment()), uintptr(16); got != want {
		t.Fatalf("segment size changed: got=%d want=%d", got, want)
	}
}
