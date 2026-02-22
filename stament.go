package obreron

// NOTE: The core internal engine (stament/scanner/build) is implemented in
// subsequent issues. This file anchors the internal types and keeps the package
// structure stable.

// segment describes a single SQL fragment stored inside stament.buf.
//
// Layout is intentionally compact (16 bytes on amd64) to improve cache locality
// during sorting and building.
//
// Invariants:
//   - If pCount == 0 then pIndex must be -1.
//   - pIndex is the base index into stament.params.
//   - start/length address a slice window within stament.buf.
//
// The actual assembly order is controlled by sType (stable-sorted).
//
// Total: 16 bytes
type segment struct {
    start  uint32 // offset in buf (max 4GiB)
    length uint32 // fragment length (max 4GiB)
    pIndex int32  // base index into params; -1 when pCount == 0
    pCount uint16 // number of params used by this fragment (max 65535)
    sType  uint8  // clause/type ordering key
    _      uint8  // explicit padding / reserved
}

func newSegment(start, length uint32, sType uint8, pIndex int32, pCount uint16) segment {
    if pCount == 0 {
        pIndex = -1
    }
    return segment{
        start:  start,
        length: length,
        pIndex: pIndex,
        pCount: pCount,
        sType:  sType,
    }
}
