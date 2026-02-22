package obreron

import "slices"

// buildInto assembles the accumulated SQL fragments into the provided buffers.
//
// Contract (spec Rev 2.0):
//   - Segments are stable-sorted by sType before assembly. Equal sType
//     preserves insertion order (stable sort).
//   - Each placeholder marker (0x1F) in buf is replaced by the
//     dialect-specific placeholder token, and the corresponding arg is
//     appended to args in assembly order (i.e. after segment sort).
//   - If dialect.MaxParams() > 0 and the total placeholder count exceeds
//     that limit, ErrTooManyParams is returned.
//   - If st.err is already set (accumulated error), it is returned
//     immediately without touching buf or args.
//   - buf and args may be nil; the function allocates as needed.
//
// Callers of the public Build() pass nil/nil to get fresh allocations
// (defensive-copy semantics are satisfied because the output slices are
// newly created by this function).
//
// Callers of the public BuildInto() pass pre-allocated slices to amortise
// allocations across repeated builds.
func (st *stament) buildInto(buf []byte, args []any) ([]byte, []any, error) {
	if st.err != nil {
		return buf, args, st.err
	}

	// Stable sort: equal sType values keep their original insertion order.
	slices.SortStableFunc(st.segs, func(a, b segment) int {
		return int(a.sType) - int(b.sType)
	})

	totalParams := 0

	for i, seg := range st.segs {
		// Separate clauses with a single space.
		if i > 0 {
			buf = append(buf, ' ')
		}

		frag := st.buf[seg.start : seg.start+seg.length]
		localIdx := 0 // index within this segment's params

		for _, b := range frag {
			if b != placeholderMarker {
				buf = append(buf, b)
				continue
			}
			// Expand placeholder marker → dialect token, append arg.
			buf = st.dialect.AppendPlaceholder(buf, totalParams)
			args = append(args, st.params[int(seg.pIndex)+localIdx])
			totalParams++
			localIdx++
		}
	}

	// Dialect-level MaxParams validation (0 = unknown / skip).
	if maxP := st.dialect.MaxParams(); maxP > 0 && totalParams > maxP {
		st.setErr("Build", ErrTooManyParams)
		return nil, nil, st.err
	}

	return buf, args, nil
}
