package algorand

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *CredentialData) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Round":
			err = z.Round.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Step":
			err = z.Step.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Quantity":
			err = z.Quantity.DecodeMsg(dc)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *CredentialData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Round"
	err = en.Append(0x83, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	if err != nil {
		return
	}
	err = z.Round.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Step"
	err = en.Append(0xa4, 0x53, 0x74, 0x65, 0x70)
	if err != nil {
		return
	}
	err = z.Step.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Quantity"
	err = en.Append(0xa8, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79)
	if err != nil {
		return
	}
	err = z.Quantity.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *CredentialData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Round"
	o = append(o, 0x83, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o, err = z.Round.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Step"
	o = append(o, 0xa4, 0x53, 0x74, 0x65, 0x70)
	o, err = z.Step.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Quantity"
	o = append(o, 0xa8, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79)
	o, err = z.Quantity.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CredentialData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Round":
			bts, err = z.Round.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Step":
			bts, err = z.Step.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Quantity":
			bts, err = z.Quantity.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *CredentialData) Msgsize() (s int) {
	s = 1 + 6 + z.Round.Msgsize() + 5 + z.Step.Msgsize() + 9 + z.Quantity.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *CredentialSig) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Round":
			err = z.Round.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Step":
			err = z.Step.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "R":
			err = z.R.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "S":
			err = z.S.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "V":
			err = z.V.DecodeMsg(dc)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *CredentialSig) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "Round"
	err = en.Append(0x85, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	if err != nil {
		return
	}
	err = z.Round.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Step"
	err = en.Append(0xa4, 0x53, 0x74, 0x65, 0x70)
	if err != nil {
		return
	}
	err = z.Step.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "R"
	err = en.Append(0xa1, 0x52)
	if err != nil {
		return
	}
	err = z.R.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "S"
	err = en.Append(0xa1, 0x53)
	if err != nil {
		return
	}
	err = z.S.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "V"
	err = en.Append(0xa1, 0x56)
	if err != nil {
		return
	}
	err = z.V.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *CredentialSig) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Round"
	o = append(o, 0x85, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o, err = z.Round.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Step"
	o = append(o, 0xa4, 0x53, 0x74, 0x65, 0x70)
	o, err = z.Step.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "R"
	o = append(o, 0xa1, 0x52)
	o, err = z.R.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "S"
	o = append(o, 0xa1, 0x53)
	o, err = z.S.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "V"
	o = append(o, 0xa1, 0x56)
	o, err = z.V.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CredentialSig) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Round":
			bts, err = z.Round.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Step":
			bts, err = z.Step.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "R":
			bts, err = z.R.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "S":
			bts, err = z.S.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "V":
			bts, err = z.V.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *CredentialSig) Msgsize() (s int) {
	s = 1 + 6 + z.Round.Msgsize() + 5 + z.Step.Msgsize() + 2 + z.R.Msgsize() + 2 + z.S.Msgsize() + 2 + z.V.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *SignatureVal) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "R":
			err = z.R.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "S":
			err = z.S.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "V":
			err = z.V.DecodeMsg(dc)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *SignatureVal) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "R"
	err = en.Append(0x83, 0xa1, 0x52)
	if err != nil {
		return
	}
	err = z.R.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "S"
	err = en.Append(0xa1, 0x53)
	if err != nil {
		return
	}
	err = z.S.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "V"
	err = en.Append(0xa1, 0x56)
	if err != nil {
		return
	}
	err = z.V.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *SignatureVal) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "R"
	o = append(o, 0x83, 0xa1, 0x52)
	o, err = z.R.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "S"
	o = append(o, 0xa1, 0x53)
	o, err = z.S.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "V"
	o = append(o, 0xa1, 0x56)
	o, err = z.V.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SignatureVal) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "R":
			bts, err = z.R.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "S":
			bts, err = z.S.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "V":
			bts, err = z.V.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *SignatureVal) Msgsize() (s int) {
	s = 1 + 2 + z.R.Msgsize() + 2 + z.S.Msgsize() + 2 + z.V.Msgsize()
	return
}
