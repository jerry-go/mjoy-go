package algorand

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ConsensusData) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Step":
			z.Step, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "Type":
			z.Type, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "Para":
			z.Para, err = dc.ReadBytes(z.Para)
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
func (z *ConsensusData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Step"
	err = en.Append(0x83, 0xa4, 0x53, 0x74, 0x65, 0x70)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Step)
	if err != nil {
		return
	}
	// write "Type"
	err = en.Append(0xa4, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Type)
	if err != nil {
		return
	}
	// write "Para"
	err = en.Append(0xa4, 0x50, 0x61, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Para)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ConsensusData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Step"
	o = append(o, 0x83, 0xa4, 0x53, 0x74, 0x65, 0x70)
	o = msgp.AppendInt(o, z.Step)
	// string "Type"
	o = append(o, 0xa4, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.Type)
	// string "Para"
	o = append(o, 0xa4, 0x50, 0x61, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Para)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ConsensusData) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Step":
			z.Step, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "Type":
			z.Type, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "Para":
			z.Para, bts, err = msgp.ReadBytesBytes(bts, z.Para)
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
func (z *ConsensusData) Msgsize() (s int) {
	s = 1 + 5 + msgp.IntSize + 5 + msgp.IntSize + 5 + msgp.BytesPrefixSize + len(z.Para)
	return
}
