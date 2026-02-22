package obreron

import "testing"

func TestNewSegment_NoParamsSetsIndexToMinusOne(t *testing.T) {
    s := newSegment(10, 20, 7, 123, 0)
    if s.pCount != 0 {
        t.Fatalf("expected pCount=0, got=%d", s.pCount)
    }
    if s.pIndex != -1 {
        t.Fatalf("expected pIndex=-1 when pCount=0, got=%d", s.pIndex)
    }
}

func TestNewSegment_WithParamsKeepsIndex(t *testing.T) {
    s := newSegment(10, 20, 7, 5, 3)
    if s.pCount != 3 {
        t.Fatalf("expected pCount=3, got=%d", s.pCount)
    }
    if s.pIndex != 5 {
        t.Fatalf("expected pIndex=5 when pCount>0, got=%d", s.pIndex)
    }
}
