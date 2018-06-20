package apos

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Credential) DecodeMsg(dc *msgp.Reader) (err error) {
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
			z.Round, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Step":
			z.Step, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Quantity":
			z.Quantity, err = dc.ReadBytes(z.Quantity)
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
func (z *Credential) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Round"
	err = en.Append(0x83, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Round)
	if err != nil {
		return
	}
	// write "Step"
	err = en.Append(0xa4, 0x53, 0x74, 0x65, 0x70)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Step)
	if err != nil {
		return
	}
	// write "Quantity"
	err = en.Append(0xa8, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Quantity)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Credential) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Round"
	o = append(o, 0x83, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendUint64(o, z.Round)
	// string "Step"
	o = append(o, 0xa4, 0x53, 0x74, 0x65, 0x70)
	o = msgp.AppendUint64(o, z.Step)
	// string "Quantity"
	o = append(o, 0xa8, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79)
	o = msgp.AppendBytes(o, z.Quantity)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Credential) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			z.Round, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Step":
			z.Step, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Quantity":
			z.Quantity, bts, err = msgp.ReadBytesBytes(bts, z.Quantity)
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
func (z *Credential) Msgsize() (s int) {
	s = 1 + 6 + msgp.Uint64Size + 5 + msgp.Uint64Size + 9 + msgp.BytesPrefixSize + len(z.Quantity)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *EphemeralSig) DecodeMsg(dc *msgp.Reader) (err error) {
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
			z.Round, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Step":
			z.Step, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Val":
			z.Val, err = dc.ReadBytes(z.Val)
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
func (z *EphemeralSig) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Round"
	err = en.Append(0x83, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Round)
	if err != nil {
		return
	}
	// write "Step"
	err = en.Append(0xa4, 0x53, 0x74, 0x65, 0x70)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Step)
	if err != nil {
		return
	}
	// write "Val"
	err = en.Append(0xa3, 0x56, 0x61, 0x6c)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Val)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *EphemeralSig) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Round"
	o = append(o, 0x83, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendUint64(o, z.Round)
	// string "Step"
	o = append(o, 0xa4, 0x53, 0x74, 0x65, 0x70)
	o = msgp.AppendUint64(o, z.Step)
	// string "Val"
	o = append(o, 0xa3, 0x56, 0x61, 0x6c)
	o = msgp.AppendBytes(o, z.Val)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EphemeralSig) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			z.Round, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Step":
			z.Step, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Val":
			z.Val, bts, err = msgp.ReadBytesBytes(bts, z.Val)
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
func (z *EphemeralSig) Msgsize() (s int) {
	s = 1 + 6 + msgp.Uint64Size + 5 + msgp.Uint64Size + 4 + msgp.BytesPrefixSize + len(z.Val)
	return
}
