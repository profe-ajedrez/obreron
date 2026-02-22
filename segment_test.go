package obreron_test

import (
	"testing"

	"github.com/profe-ajedrez/obreron/v3"
)

func TestNewSegment_NoParamsSetsIndexToMinusOne(t *testing.T) {
	s := obreron.NewSegment(10, 20, 7, 123, 0)
	if obreron.GetParamsCount(s) != 0 {
		t.Fatalf("expected pCount=0, got=%d", obreron.GetParamsCount(s))
	}
	if obreron.GetParamIndex(s) != -1 {
		t.Fatalf("expected pIndex=-1 when pCount=0, got=%d", obreron.GetParamIndex(s))
	}
}

func TestNewSegment_WithParamsKeepsIndex(t *testing.T) {
	s := obreron.NewSegment(10, 20, 7, 5, 3)
	if obreron.GetParamsCount(s) != 3 {
		t.Fatalf("expected pCount=3, got=%d", obreron.GetParamsCount(s))
	}
	if obreron.GetParamIndex(s) != 5 {
		t.Fatalf("expected pIndex=5 when pCount>0, got=%d", obreron.GetParamIndex(s))
	}
}
