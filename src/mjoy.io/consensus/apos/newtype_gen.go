package apos

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"mjoy.io/core/blockchain/block"
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
func (z *CredentialSigForKey) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "R":
			z.R, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "S":
			z.S, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "V":
			z.V, err = dc.ReadUint64()
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
func (z *CredentialSigForKey) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "Round"
	err = en.Append(0x85, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
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
	// write "R"
	err = en.Append(0xa1, 0x52)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.R)
	if err != nil {
		return
	}
	// write "S"
	err = en.Append(0xa1, 0x53)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.S)
	if err != nil {
		return
	}
	// write "V"
	err = en.Append(0xa1, 0x56)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.V)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *CredentialSigForKey) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Round"
	o = append(o, 0x85, 0xa5, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendUint64(o, z.Round)
	// string "Step"
	o = append(o, 0xa4, 0x53, 0x74, 0x65, 0x70)
	o = msgp.AppendUint64(o, z.Step)
	// string "R"
	o = append(o, 0xa1, 0x52)
	o = msgp.AppendUint64(o, z.R)
	// string "S"
	o = append(o, 0xa1, 0x53)
	o = msgp.AppendUint64(o, z.S)
	// string "V"
	o = append(o, 0xa1, 0x56)
	o = msgp.AppendUint64(o, z.V)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CredentialSigForKey) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "R":
			z.R, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "S":
			z.S, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "V":
			z.V, bts, err = msgp.ReadUint64Bytes(bts)
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
func (z *CredentialSigForKey) Msgsize() (s int) {
	s = 1 + 6 + msgp.Uint64Size + 5 + msgp.Uint64Size + 2 + msgp.Uint64Size + 2 + msgp.Uint64Size + 2 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *CredentialSigStatus) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z CredentialSigStatus) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 0
	err = en.Append(0x80)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z CredentialSigStatus) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 0
	o = append(o, 0x80)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CredentialSigStatus) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z CredentialSigStatus) Msgsize() (s int) {
	s = 1
	return
}

// DecodeMsg implements msgp.Decodable
func (z *CredentialSigStatusHeap) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(CredentialSigStatusHeap, zb0002)
	}
	for zb0001 := range *z {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(CredentialSigStatus)
			}
			var field []byte
			_ = field
			var zb0003 uint32
			zb0003, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			for zb0003 > 0 {
				zb0003--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				default:
					err = dc.Skip()
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z CredentialSigStatusHeap) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0004 := range z {
		if z[zb0004] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// map header, size 0
			err = en.Append(0x80)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z CredentialSigStatusHeap) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		if z[zb0004] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 0
			o = append(o, 0x80)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CredentialSigStatusHeap) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(CredentialSigStatusHeap, zb0002)
	}
	for zb0001 := range *z {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(CredentialSigStatus)
			}
			var field []byte
			_ = field
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			for zb0003 > 0 {
				zb0003--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						return
					}
				}
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z CredentialSigStatusHeap) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		if z[zb0004] == nil {
			s += msgp.NilSize
		} else {
			s += 1
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *M1) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Block":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Block = nil
			} else {
				if z.Block == nil {
					z.Block = new(block.Block)
				}
				err = z.Block.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Esig":
			z.Esig, err = dc.ReadBytes(z.Esig)
			if err != nil {
				return
			}
		case "Credential":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Credential = nil
			} else {
				if z.Credential == nil {
					z.Credential = new(CredentialSig)
				}
				err = z.Credential.DecodeMsg(dc)
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
func (z *M1) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Block"
	err = en.Append(0x83, 0xa5, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	if z.Block == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Block.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Esig"
	err = en.Append(0xa4, 0x45, 0x73, 0x69, 0x67)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Esig)
	if err != nil {
		return
	}
	// write "Credential"
	err = en.Append(0xaa, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c)
	if err != nil {
		return
	}
	if z.Credential == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Credential.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *M1) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Block"
	o = append(o, 0x83, 0xa5, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if z.Block == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Block.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Esig"
	o = append(o, 0xa4, 0x45, 0x73, 0x69, 0x67)
	o = msgp.AppendBytes(o, z.Esig)
	// string "Credential"
	o = append(o, 0xaa, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c)
	if z.Credential == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Credential.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *M1) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Block":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Block = nil
			} else {
				if z.Block == nil {
					z.Block = new(block.Block)
				}
				bts, err = z.Block.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Esig":
			z.Esig, bts, err = msgp.ReadBytesBytes(bts, z.Esig)
			if err != nil {
				return
			}
		case "Credential":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Credential = nil
			} else {
				if z.Credential == nil {
					z.Credential = new(CredentialSig)
				}
				bts, err = z.Credential.UnmarshalMsg(bts)
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
func (z *M1) Msgsize() (s int) {
	s = 1 + 6
	if z.Block == nil {
		s += msgp.NilSize
	} else {
		s += z.Block.Msgsize()
	}
	s += 5 + msgp.BytesPrefixSize + len(z.Esig) + 11
	if z.Credential == nil {
		s += msgp.NilSize
	} else {
		s += z.Credential.Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *M23) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Hash":
			err = z.Hash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Esig":
			z.Esig, err = dc.ReadBytes(z.Esig)
			if err != nil {
				return
			}
		case "Credential":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Credential = nil
			} else {
				if z.Credential == nil {
					z.Credential = new(CredentialSig)
				}
				err = z.Credential.DecodeMsg(dc)
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
func (z *M23) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Hash"
	err = en.Append(0x83, 0xa4, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.Hash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Esig"
	err = en.Append(0xa4, 0x45, 0x73, 0x69, 0x67)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Esig)
	if err != nil {
		return
	}
	// write "Credential"
	err = en.Append(0xaa, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c)
	if err != nil {
		return
	}
	if z.Credential == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Credential.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *M23) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Hash"
	o = append(o, 0x83, 0xa4, 0x48, 0x61, 0x73, 0x68)
	o, err = z.Hash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Esig"
	o = append(o, 0xa4, 0x45, 0x73, 0x69, 0x67)
	o = msgp.AppendBytes(o, z.Esig)
	// string "Credential"
	o = append(o, 0xaa, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c)
	if z.Credential == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Credential.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *M23) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Hash":
			bts, err = z.Hash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Esig":
			z.Esig, bts, err = msgp.ReadBytesBytes(bts, z.Esig)
			if err != nil {
				return
			}
		case "Credential":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Credential = nil
			} else {
				if z.Credential == nil {
					z.Credential = new(CredentialSig)
				}
				bts, err = z.Credential.UnmarshalMsg(bts)
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
func (z *M23) Msgsize() (s int) {
	s = 1 + 5 + z.Hash.Msgsize() + 5 + msgp.BytesPrefixSize + len(z.Esig) + 11
	if z.Credential == nil {
		s += msgp.NilSize
	} else {
		s += z.Credential.Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *MCommon) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "B":
			z.B, err = dc.ReadUint()
			if err != nil {
				return
			}
		case "EsigB":
			z.EsigB, err = dc.ReadBytes(z.EsigB)
			if err != nil {
				return
			}
		case "Hash":
			err = z.Hash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "EsigV":
			z.EsigV, err = dc.ReadBytes(z.EsigV)
			if err != nil {
				return
			}
		case "Credential":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Credential = nil
			} else {
				if z.Credential == nil {
					z.Credential = new(CredentialSig)
				}
				err = z.Credential.DecodeMsg(dc)
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
func (z *MCommon) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "B"
	err = en.Append(0x85, 0xa1, 0x42)
	if err != nil {
		return
	}
	err = en.WriteUint(z.B)
	if err != nil {
		return
	}
	// write "EsigB"
	err = en.Append(0xa5, 0x45, 0x73, 0x69, 0x67, 0x42)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.EsigB)
	if err != nil {
		return
	}
	// write "Hash"
	err = en.Append(0xa4, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.Hash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "EsigV"
	err = en.Append(0xa5, 0x45, 0x73, 0x69, 0x67, 0x56)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.EsigV)
	if err != nil {
		return
	}
	// write "Credential"
	err = en.Append(0xaa, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c)
	if err != nil {
		return
	}
	if z.Credential == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Credential.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *MCommon) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "B"
	o = append(o, 0x85, 0xa1, 0x42)
	o = msgp.AppendUint(o, z.B)
	// string "EsigB"
	o = append(o, 0xa5, 0x45, 0x73, 0x69, 0x67, 0x42)
	o = msgp.AppendBytes(o, z.EsigB)
	// string "Hash"
	o = append(o, 0xa4, 0x48, 0x61, 0x73, 0x68)
	o, err = z.Hash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "EsigV"
	o = append(o, 0xa5, 0x45, 0x73, 0x69, 0x67, 0x56)
	o = msgp.AppendBytes(o, z.EsigV)
	// string "Credential"
	o = append(o, 0xaa, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c)
	if z.Credential == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Credential.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MCommon) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "B":
			z.B, bts, err = msgp.ReadUintBytes(bts)
			if err != nil {
				return
			}
		case "EsigB":
			z.EsigB, bts, err = msgp.ReadBytesBytes(bts, z.EsigB)
			if err != nil {
				return
			}
		case "Hash":
			bts, err = z.Hash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "EsigV":
			z.EsigV, bts, err = msgp.ReadBytesBytes(bts, z.EsigV)
			if err != nil {
				return
			}
		case "Credential":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Credential = nil
			} else {
				if z.Credential == nil {
					z.Credential = new(CredentialSig)
				}
				bts, err = z.Credential.UnmarshalMsg(bts)
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
func (z *MCommon) Msgsize() (s int) {
	s = 1 + 2 + msgp.UintSize + 6 + msgp.BytesPrefixSize + len(z.EsigB) + 5 + z.Hash.Msgsize() + 6 + msgp.BytesPrefixSize + len(z.EsigV) + 11
	if z.Credential == nil {
		s += msgp.NilSize
	} else {
		s += z.Credential.Msgsize()
	}
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
