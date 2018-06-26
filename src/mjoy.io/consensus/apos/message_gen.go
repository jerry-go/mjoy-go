package apos

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"mjoy.io/core/blockchain/block"
)

// DecodeMsg implements msgp.Decodable
func (z *BinaryByzantineAgreement) DecodeMsg(dc *msgp.Reader) (err error) {
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
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.EsigB = nil
			} else {
				if z.EsigB == nil {
					z.EsigB = new(EphemeralSign)
				}
				err = z.EsigB.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Hash":
			err = z.Hash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "EsigV":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.EsigV = nil
			} else {
				if z.EsigV == nil {
					z.EsigV = new(EphemeralSign)
				}
				err = z.EsigV.DecodeMsg(dc)
				if err != nil {
					return
				}
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
					z.Credential = new(CredentialSign)
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
func (z *BinaryByzantineAgreement) EncodeMsg(en *msgp.Writer) (err error) {
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
	if z.EsigB == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.EsigB.EncodeMsg(en)
		if err != nil {
			return
		}
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
	if z.EsigV == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.EsigV.EncodeMsg(en)
		if err != nil {
			return
		}
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
func (z *BinaryByzantineAgreement) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "B"
	o = append(o, 0x85, 0xa1, 0x42)
	o = msgp.AppendUint(o, z.B)
	// string "EsigB"
	o = append(o, 0xa5, 0x45, 0x73, 0x69, 0x67, 0x42)
	if z.EsigB == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.EsigB.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Hash"
	o = append(o, 0xa4, 0x48, 0x61, 0x73, 0x68)
	o, err = z.Hash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "EsigV"
	o = append(o, 0xa5, 0x45, 0x73, 0x69, 0x67, 0x56)
	if z.EsigV == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.EsigV.MarshalMsg(o)
		if err != nil {
			return
		}
	}
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
func (z *BinaryByzantineAgreement) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.EsigB = nil
			} else {
				if z.EsigB == nil {
					z.EsigB = new(EphemeralSign)
				}
				bts, err = z.EsigB.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Hash":
			bts, err = z.Hash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "EsigV":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.EsigV = nil
			} else {
				if z.EsigV == nil {
					z.EsigV = new(EphemeralSign)
				}
				bts, err = z.EsigV.UnmarshalMsg(bts)
				if err != nil {
					return
				}
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
					z.Credential = new(CredentialSign)
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
func (z *BinaryByzantineAgreement) Msgsize() (s int) {
	s = 1 + 2 + msgp.UintSize + 6
	if z.EsigB == nil {
		s += msgp.NilSize
	} else {
		s += z.EsigB.Msgsize()
	}
	s += 5 + z.Hash.Msgsize() + 6
	if z.EsigV == nil {
		s += msgp.NilSize
	} else {
		s += z.EsigV.Msgsize()
	}
	s += 11
	if z.Credential == nil {
		s += msgp.NilSize
	} else {
		s += z.Credential.Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BlockProposal) DecodeMsg(dc *msgp.Reader) (err error) {
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
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Esig = nil
			} else {
				if z.Esig == nil {
					z.Esig = new(EphemeralSign)
				}
				err = z.Esig.DecodeMsg(dc)
				if err != nil {
					return
				}
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
					z.Credential = new(CredentialSign)
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
func (z *BlockProposal) EncodeMsg(en *msgp.Writer) (err error) {
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
	if z.Esig == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Esig.EncodeMsg(en)
		if err != nil {
			return
		}
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
func (z *BlockProposal) MarshalMsg(b []byte) (o []byte, err error) {
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
	if z.Esig == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Esig.MarshalMsg(o)
		if err != nil {
			return
		}
	}
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
func (z *BlockProposal) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Esig = nil
			} else {
				if z.Esig == nil {
					z.Esig = new(EphemeralSign)
				}
				bts, err = z.Esig.UnmarshalMsg(bts)
				if err != nil {
					return
				}
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
					z.Credential = new(CredentialSign)
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
func (z *BlockProposal) Msgsize() (s int) {
	s = 1 + 6
	if z.Block == nil {
		s += msgp.NilSize
	} else {
		s += z.Block.Msgsize()
	}
	s += 5
	if z.Esig == nil {
		s += msgp.NilSize
	} else {
		s += z.Esig.Msgsize()
	}
	s += 11
	if z.Credential == nil {
		s += msgp.NilSize
	} else {
		s += z.Credential.Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *GradedConsensus) DecodeMsg(dc *msgp.Reader) (err error) {
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
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Esig = nil
			} else {
				if z.Esig == nil {
					z.Esig = new(EphemeralSign)
				}
				err = z.Esig.DecodeMsg(dc)
				if err != nil {
					return
				}
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
					z.Credential = new(CredentialSign)
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
func (z *GradedConsensus) EncodeMsg(en *msgp.Writer) (err error) {
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
	if z.Esig == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Esig.EncodeMsg(en)
		if err != nil {
			return
		}
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
func (z *GradedConsensus) MarshalMsg(b []byte) (o []byte, err error) {
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
	if z.Esig == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Esig.MarshalMsg(o)
		if err != nil {
			return
		}
	}
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
func (z *GradedConsensus) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Esig = nil
			} else {
				if z.Esig == nil {
					z.Esig = new(EphemeralSign)
				}
				bts, err = z.Esig.UnmarshalMsg(bts)
				if err != nil {
					return
				}
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
					z.Credential = new(CredentialSign)
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
func (z *GradedConsensus) Msgsize() (s int) {
	s = 1 + 5 + z.Hash.Msgsize() + 5
	if z.Esig == nil {
		s += msgp.NilSize
	} else {
		s += z.Esig.Msgsize()
	}
	s += 11
	if z.Credential == nil {
		s += msgp.NilSize
	} else {
		s += z.Credential.Msgsize()
	}
	return
}
