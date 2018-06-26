package apos

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"mjoy.io/common/types"
)

// DecodeMsg implements msgp.Decodable
func (z *CredentialSigForHash) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z *CredentialSigForHash) EncodeMsg(en *msgp.Writer) (err error) {
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
func (z *CredentialSigForHash) MarshalMsg(b []byte) (o []byte, err error) {
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
func (z *CredentialSigForHash) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z *CredentialSigForHash) Msgsize() (s int) {
	s = 1 + 6 + msgp.Uint64Size + 5 + msgp.Uint64Size + 9 + msgp.BytesPrefixSize + len(z.Quantity)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *CredentialSign) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Signature":
			err = z.Signature.DecodeMsg(dc)
			if err != nil {
				return
			}
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
func (z *CredentialSign) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Signature"
	err = en.Append(0x83, 0xa9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65)
	if err != nil {
		return
	}
	err = z.Signature.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Round"
	err = en.Append(0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
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
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *CredentialSign) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Signature"
	o = append(o, 0x83, 0xa9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65)
	o, err = z.Signature.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Round"
	o = append(o, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendUint64(o, z.Round)
	// string "Step"
	o = append(o, 0xa4, 0x53, 0x74, 0x65, 0x70)
	o = msgp.AppendUint64(o, z.Step)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CredentialSign) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Signature":
			bts, err = z.Signature.UnmarshalMsg(bts)
			if err != nil {
				return
			}
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
func (z *CredentialSign) Msgsize() (s int) {
	s = 1 + 10 + z.Signature.Msgsize() + 6 + msgp.Uint64Size + 5 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *EphemeralSign) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Signature":
			err = z.Signature.DecodeMsg(dc)
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
func (z *EphemeralSign) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Signature"
	err = en.Append(0x81, 0xa9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65)
	if err != nil {
		return
	}
	err = z.Signature.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *EphemeralSign) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Signature"
	o = append(o, 0x81, 0xa9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65)
	o, err = z.Signature.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EphemeralSign) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Signature":
			bts, err = z.Signature.UnmarshalMsg(bts)
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
func (z *EphemeralSign) Msgsize() (s int) {
	s = 1 + 10 + z.Signature.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Signature) DecodeMsg(dc *msgp.Reader) (err error) {
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
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.R = nil
			} else {
				if z.R == nil {
					z.R = new(types.BigInt)
				}
				err = z.R.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "S":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.S = nil
			} else {
				if z.S == nil {
					z.S = new(types.BigInt)
				}
				err = z.S.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "V":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.V = nil
			} else {
				if z.V == nil {
					z.V = new(types.BigInt)
				}
				err = z.V.DecodeMsg(dc)
				if err != nil {
					return
				}
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
func (z *Signature) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "R"
	err = en.Append(0x83, 0xa1, 0x52)
	if err != nil {
		return
	}
	if z.R == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.R.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "S"
	err = en.Append(0xa1, 0x53)
	if err != nil {
		return
	}
	if z.S == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.S.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "V"
	err = en.Append(0xa1, 0x56)
	if err != nil {
		return
	}
	if z.V == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.V.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Signature) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "R"
	o = append(o, 0x83, 0xa1, 0x52)
	if z.R == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.R.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "S"
	o = append(o, 0xa1, 0x53)
	if z.S == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.S.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "V"
	o = append(o, 0xa1, 0x56)
	if z.V == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.V.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Signature) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.R = nil
			} else {
				if z.R == nil {
					z.R = new(types.BigInt)
				}
				bts, err = z.R.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "S":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.S = nil
			} else {
				if z.S == nil {
					z.S = new(types.BigInt)
				}
				bts, err = z.S.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "V":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.V = nil
			} else {
				if z.V == nil {
					z.V = new(types.BigInt)
				}
				bts, err = z.V.UnmarshalMsg(bts)
				if err != nil {
					return
				}
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
func (z *Signature) Msgsize() (s int) {
	s = 1 + 2
	if z.R == nil {
		s += msgp.NilSize
	} else {
		s += z.R.Msgsize()
	}
	s += 2
	if z.S == nil {
		s += msgp.NilSize
	} else {
		s += z.S.Msgsize()
	}
	s += 2
	if z.V == nil {
		s += msgp.NilSize
	} else {
		s += z.V.Msgsize()
	}
	return
}
