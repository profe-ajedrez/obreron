//go:build amd64

package obreron

import (
    "testing"
    "unsafe"
)

func TestSegmentSizeAMD64(t *testing.T) {
    if got, want := unsafe.Sizeof(segment{}), uintptr(16); got != want {
        t.Fatalf("segment size changed: got=%d want=%d", got, want)
    }
}
