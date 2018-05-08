package state

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"mjoy.io/common/types"
)

// DecodeMsg implements msgp.Decodable
func (z *Account) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Nonce":
			z.Nonce, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Balance":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Balance = nil
			} else {
				if z.Balance == nil {
					z.Balance = new(types.BigInt)
				}
				err = z.Balance.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Root":
			err = z.Root.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "CodeHash":
			z.CodeHash, err = dc.ReadBytes(z.CodeHash)
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
func (z *Account) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Nonce"
	err = en.Append(0x84, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Nonce)
	if err != nil {
		return
	}
	// write "Balance"
	err = en.Append(0xa7, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	if z.Balance == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Balance.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Root"
	err = en.Append(0xa4, 0x52, 0x6f, 0x6f, 0x74)
	if err != nil {
		return
	}
	err = z.Root.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "CodeHash"
	err = en.Append(0xa8, 0x43, 0x6f, 0x64, 0x65, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.CodeHash)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Account) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Nonce"
	o = append(o, 0x84, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.Nonce)
	// string "Balance"
	o = append(o, 0xa7, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65)
	if z.Balance == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Balance.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Root"
	o = append(o, 0xa4, 0x52, 0x6f, 0x6f, 0x74)
	o, err = z.Root.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "CodeHash"
	o = append(o, 0xa8, 0x43, 0x6f, 0x64, 0x65, 0x48, 0x61, 0x73, 0x68)
	o = msgp.AppendBytes(o, z.CodeHash)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Account) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Nonce":
			z.Nonce, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Balance":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Balance = nil
			} else {
				if z.Balance == nil {
					z.Balance = new(types.BigInt)
				}
				bts, err = z.Balance.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Root":
			bts, err = z.Root.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "CodeHash":
			z.CodeHash, bts, err = msgp.ReadBytesBytes(bts, z.CodeHash)
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
func (z *Account) Msgsize() (s int) {
	s = 1 + 6 + msgp.Uint64Size + 8
	if z.Balance == nil {
		s += msgp.NilSize
	} else {
		s += z.Balance.Msgsize()
	}
	s += 5 + z.Root.Msgsize() + 9 + msgp.BytesPrefixSize + len(z.CodeHash)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Code) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 []byte
		zb0001, err = dc.ReadBytes([]byte((*z)))
		if err != nil {
			return
		}
		(*z) = Code(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Code) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes([]byte(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Code) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, []byte(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Code) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 []byte
		zb0001, bts, err = msgp.ReadBytesBytes(bts, []byte((*z)))
		if err != nil {
			return
		}
		(*z) = Code(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Code) Msgsize() (s int) {
	s = msgp.BytesPrefixSize + len([]byte(z))
	return
}
