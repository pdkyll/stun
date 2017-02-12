package stun

import "errors"

// UnknownAttributes represents UNKNOWN-ATTRIBUTES
type UnknownAttributes struct {
	Types []AttrType
}

func (a UnknownAttributes) String() string {
	s := ""
	if len(a.Types) == 0 {
		return s + "<nil>"
	}
	last := len(a.Types) - 1
	for i, t := range a.Types {
		s += t.String()
		if i != last {
			s += ", "
		}
	}
	return s
}

// type size is 16 bit.
const attrTypeSize = 4

// AddTo adds UNKNOWN-ATTRIBUTES attribute to message.
func (a *UnknownAttributes) AddTo(m *Message) error {
	v := make([]byte, 0, attrTypeSize*20) // 20 should be enough
	// If len(a.Types) > 20, there will be allocations
	for i, t := range a.Types {
		v = append(v, 0, 0, 0, 0) // 4 times by 0 (16 bits)1
		first := attrTypeSize * i
		last := first + attrTypeSize
		bin.PutUint16(v[first:last], t.Value())
	}
	m.Add(AttrUnknownAttributes, v)
	return nil
}

// ErrBadUnknownAttrsSize means that UNKNOWN-ATTRIBUTES attribute value
// has invalid length.
var ErrBadUnknownAttrsSize = errors.New("bad UNKNOWN-ATTRIBUTES size")

// GetFrom appends UNKNOWN-ATTRIBUTES from m to a.Types,
// returning error if any.
func (a *UnknownAttributes) GetFrom(m *Message) error {
	v, err := m.Get(AttrUnknownAttributes)
	if err != nil {
		return err
	}
	if len(v)%attrTypeSize != 0 {
		return ErrBadUnknownAttrsSize
	}
	first := 0
	for first < len(v) {
		last := first + attrTypeSize
		a.Types = append(a.Types,
			AttrType(bin.Uint16(v[first:last])),
		)
		first = last
	}
	return nil
}