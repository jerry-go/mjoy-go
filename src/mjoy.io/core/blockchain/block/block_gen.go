package block

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
	"mjoy.io/common/types"
	"mjoy.io/core/transaction"
)

// DecodeMsg implements msgp.Decodable
func (z *Block) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "B_header":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.B_header = nil
			} else {
				if z.B_header == nil {
					z.B_header = new(Header)
				}
				err = z.B_header.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "B_body":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Transactions":
					var zb0003 uint32
					zb0003, err = dc.ReadArrayHeader()
					if err != nil {
						return
					}
					if cap(z.B_body.Transactions) >= int(zb0003) {
						z.B_body.Transactions = (z.B_body.Transactions)[:zb0003]
					} else {
						z.B_body.Transactions = make([]*transaction.Transaction, zb0003)
					}
					for za0001 := range z.B_body.Transactions {
						if dc.IsNil() {
							err = dc.ReadNil()
							if err != nil {
								return
							}
							z.B_body.Transactions[za0001] = nil
						} else {
							if z.B_body.Transactions[za0001] == nil {
								z.B_body.Transactions[za0001] = new(transaction.Transaction)
							}
							err = z.B_body.Transactions[za0001].DecodeMsg(dc)
							if err != nil {
								return
							}
						}
					}
				default:
					err = dc.Skip()
					if err != nil {
						return
					}
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
func (z *Block) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "B_header"
	err = en.Append(0x82, 0xa8, 0x42, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72)
	if err != nil {
		return
	}
	if z.B_header == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.B_header.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "B_body"
	// map header, size 1
	// write "Transactions"
	err = en.Append(0xa6, 0x42, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.B_body.Transactions)))
	if err != nil {
		return
	}
	for za0001 := range z.B_body.Transactions {
		if z.B_body.Transactions[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.B_body.Transactions[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Block) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "B_header"
	o = append(o, 0x82, 0xa8, 0x42, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72)
	if z.B_header == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.B_header.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "B_body"
	// map header, size 1
	// string "Transactions"
	o = append(o, 0xa6, 0x42, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.B_body.Transactions)))
	for za0001 := range z.B_body.Transactions {
		if z.B_body.Transactions[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.B_body.Transactions[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Block) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "B_header":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.B_header = nil
			} else {
				if z.B_header == nil {
					z.B_header = new(Header)
				}
				bts, err = z.B_header.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "B_body":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Transactions":
					var zb0003 uint32
					zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
					if err != nil {
						return
					}
					if cap(z.B_body.Transactions) >= int(zb0003) {
						z.B_body.Transactions = (z.B_body.Transactions)[:zb0003]
					} else {
						z.B_body.Transactions = make([]*transaction.Transaction, zb0003)
					}
					for za0001 := range z.B_body.Transactions {
						if msgp.IsNil(bts) {
							bts, err = msgp.ReadNilBytes(bts)
							if err != nil {
								return
							}
							z.B_body.Transactions[za0001] = nil
						} else {
							if z.B_body.Transactions[za0001] == nil {
								z.B_body.Transactions[za0001] = new(transaction.Transaction)
							}
							bts, err = z.B_body.Transactions[za0001].UnmarshalMsg(bts)
							if err != nil {
								return
							}
						}
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						return
					}
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
func (z *Block) Msgsize() (s int) {
	s = 1 + 9
	if z.B_header == nil {
		s += msgp.NilSize
	} else {
		s += z.B_header.Msgsize()
	}
	s += 7 + 1 + 13 + msgp.ArrayHeaderSize
	for za0001 := range z.B_body.Transactions {
		if z.B_body.Transactions[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.B_body.Transactions[za0001].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BlockNonce) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 interface{}
		zb0001, err = dc.ReadIntf()
		if err != nil {
			return
		}
		(*z), err = fromBytes(zb0001)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z BlockNonce) EncodeMsg(en *msgp.Writer) (err error) {
	var zb0001 interface{}
	zb0001, err = toBytes(z)
	if err != nil {
		return
	}
	err = en.WriteIntf(zb0001)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z BlockNonce) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	var zb0001 interface{}
	zb0001, err = toBytes(z)
	if err != nil {
		return
	}
	o, err = msgp.AppendIntf(o, zb0001)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlockNonce) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 interface{}
		zb0001, bts, err = msgp.ReadIntfBytes(bts)
		if err != nil {
			return
		}
		(*z), err = fromBytes(zb0001)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z BlockNonce) Msgsize() (s int) {
	var zb0001 interface{}
	_ = z
	s += msgp.GuessSize(zb0001)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Blocks) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0003 uint32
	zb0003, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0003) {
		(*z) = (*z)[:zb0003]
	} else {
		(*z) = make(Blocks, zb0003)
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
				(*z)[zb0001] = new(Block)
			}
			var field []byte
			_ = field
			var zb0004 uint32
			zb0004, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			for zb0004 > 0 {
				zb0004--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "B_header":
					if dc.IsNil() {
						err = dc.ReadNil()
						if err != nil {
							return
						}
						(*z)[zb0001].B_header = nil
					} else {
						if (*z)[zb0001].B_header == nil {
							(*z)[zb0001].B_header = new(Header)
						}
						err = (*z)[zb0001].B_header.DecodeMsg(dc)
						if err != nil {
							return
						}
					}
				case "B_body":
					var zb0005 uint32
					zb0005, err = dc.ReadMapHeader()
					if err != nil {
						return
					}
					for zb0005 > 0 {
						zb0005--
						field, err = dc.ReadMapKeyPtr()
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "Transactions":
							var zb0006 uint32
							zb0006, err = dc.ReadArrayHeader()
							if err != nil {
								return
							}
							if cap((*z)[zb0001].B_body.Transactions) >= int(zb0006) {
								(*z)[zb0001].B_body.Transactions = ((*z)[zb0001].B_body.Transactions)[:zb0006]
							} else {
								(*z)[zb0001].B_body.Transactions = make([]*transaction.Transaction, zb0006)
							}
							for zb0002 := range (*z)[zb0001].B_body.Transactions {
								if dc.IsNil() {
									err = dc.ReadNil()
									if err != nil {
										return
									}
									(*z)[zb0001].B_body.Transactions[zb0002] = nil
								} else {
									if (*z)[zb0001].B_body.Transactions[zb0002] == nil {
										(*z)[zb0001].B_body.Transactions[zb0002] = new(transaction.Transaction)
									}
									err = (*z)[zb0001].B_body.Transactions[zb0002].DecodeMsg(dc)
									if err != nil {
										return
									}
								}
							}
						default:
							err = dc.Skip()
							if err != nil {
								return
							}
						}
					}
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
func (z Blocks) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0007 := range z {
		if z[zb0007] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// map header, size 2
			// write "B_header"
			err = en.Append(0x82, 0xa8, 0x42, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72)
			if err != nil {
				return
			}
			if z[zb0007].B_header == nil {
				err = en.WriteNil()
				if err != nil {
					return
				}
			} else {
				err = z[zb0007].B_header.EncodeMsg(en)
				if err != nil {
					return
				}
			}
			// write "B_body"
			// map header, size 1
			// write "Transactions"
			err = en.Append(0xa6, 0x42, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
			if err != nil {
				return
			}
			err = en.WriteArrayHeader(uint32(len(z[zb0007].B_body.Transactions)))
			if err != nil {
				return
			}
			for zb0008 := range z[zb0007].B_body.Transactions {
				if z[zb0007].B_body.Transactions[zb0008] == nil {
					err = en.WriteNil()
					if err != nil {
						return
					}
				} else {
					err = z[zb0007].B_body.Transactions[zb0008].EncodeMsg(en)
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Blocks) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0007 := range z {
		if z[zb0007] == nil {
			o = msgp.AppendNil(o)
		} else {
			// map header, size 2
			// string "B_header"
			o = append(o, 0x82, 0xa8, 0x42, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72)
			if z[zb0007].B_header == nil {
				o = msgp.AppendNil(o)
			} else {
				o, err = z[zb0007].B_header.MarshalMsg(o)
				if err != nil {
					return
				}
			}
			// string "B_body"
			// map header, size 1
			// string "Transactions"
			o = append(o, 0xa6, 0x42, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
			o = msgp.AppendArrayHeader(o, uint32(len(z[zb0007].B_body.Transactions)))
			for zb0008 := range z[zb0007].B_body.Transactions {
				if z[zb0007].B_body.Transactions[zb0008] == nil {
					o = msgp.AppendNil(o)
				} else {
					o, err = z[zb0007].B_body.Transactions[zb0008].MarshalMsg(o)
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Blocks) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0003 uint32
	zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0003) {
		(*z) = (*z)[:zb0003]
	} else {
		(*z) = make(Blocks, zb0003)
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
				(*z)[zb0001] = new(Block)
			}
			var field []byte
			_ = field
			var zb0004 uint32
			zb0004, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			for zb0004 > 0 {
				zb0004--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "B_header":
					if msgp.IsNil(bts) {
						bts, err = msgp.ReadNilBytes(bts)
						if err != nil {
							return
						}
						(*z)[zb0001].B_header = nil
					} else {
						if (*z)[zb0001].B_header == nil {
							(*z)[zb0001].B_header = new(Header)
						}
						bts, err = (*z)[zb0001].B_header.UnmarshalMsg(bts)
						if err != nil {
							return
						}
					}
				case "B_body":
					var zb0005 uint32
					zb0005, bts, err = msgp.ReadMapHeaderBytes(bts)
					if err != nil {
						return
					}
					for zb0005 > 0 {
						zb0005--
						field, bts, err = msgp.ReadMapKeyZC(bts)
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "Transactions":
							var zb0006 uint32
							zb0006, bts, err = msgp.ReadArrayHeaderBytes(bts)
							if err != nil {
								return
							}
							if cap((*z)[zb0001].B_body.Transactions) >= int(zb0006) {
								(*z)[zb0001].B_body.Transactions = ((*z)[zb0001].B_body.Transactions)[:zb0006]
							} else {
								(*z)[zb0001].B_body.Transactions = make([]*transaction.Transaction, zb0006)
							}
							for zb0002 := range (*z)[zb0001].B_body.Transactions {
								if msgp.IsNil(bts) {
									bts, err = msgp.ReadNilBytes(bts)
									if err != nil {
										return
									}
									(*z)[zb0001].B_body.Transactions[zb0002] = nil
								} else {
									if (*z)[zb0001].B_body.Transactions[zb0002] == nil {
										(*z)[zb0001].B_body.Transactions[zb0002] = new(transaction.Transaction)
									}
									bts, err = (*z)[zb0001].B_body.Transactions[zb0002].UnmarshalMsg(bts)
									if err != nil {
										return
									}
								}
							}
						default:
							bts, err = msgp.Skip(bts)
							if err != nil {
								return
							}
						}
					}
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
func (z Blocks) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0007 := range z {
		if z[zb0007] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + 9
			if z[zb0007].B_header == nil {
				s += msgp.NilSize
			} else {
				s += z[zb0007].B_header.Msgsize()
			}
			s += 7 + 1 + 13 + msgp.ArrayHeaderSize
			for zb0008 := range z[zb0007].B_body.Transactions {
				if z[zb0007].B_body.Transactions[zb0008] == nil {
					s += msgp.NilSize
				} else {
					s += z[zb0007].B_body.Transactions[zb0008].Msgsize()
				}
			}
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Body) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Transactions":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Transactions) >= int(zb0002) {
				z.Transactions = (z.Transactions)[:zb0002]
			} else {
				z.Transactions = make([]*transaction.Transaction, zb0002)
			}
			for za0001 := range z.Transactions {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Transactions[za0001] = nil
				} else {
					if z.Transactions[za0001] == nil {
						z.Transactions[za0001] = new(transaction.Transaction)
					}
					err = z.Transactions[za0001].DecodeMsg(dc)
					if err != nil {
						return
					}
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
func (z *Body) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Transactions"
	err = en.Append(0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Transactions)))
	if err != nil {
		return
	}
	for za0001 := range z.Transactions {
		if z.Transactions[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Transactions[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Body) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Transactions"
	o = append(o, 0x81, 0xac, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Transactions)))
	for za0001 := range z.Transactions {
		if z.Transactions[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Transactions[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Body) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Transactions":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Transactions) >= int(zb0002) {
				z.Transactions = (z.Transactions)[:zb0002]
			} else {
				z.Transactions = make([]*transaction.Transaction, zb0002)
			}
			for za0001 := range z.Transactions {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Transactions[za0001] = nil
				} else {
					if z.Transactions[za0001] == nil {
						z.Transactions[za0001] = new(transaction.Transaction)
					}
					bts, err = z.Transactions[za0001].UnmarshalMsg(bts)
					if err != nil {
						return
					}
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
func (z *Body) Msgsize() (s int) {
	s = 1 + 13 + msgp.ArrayHeaderSize
	for za0001 := range z.Transactions {
		if z.Transactions[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Transactions[za0001].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Header) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ParentHash":
			err = z.ParentHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Coinbase":
			err = z.Coinbase.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "StateHash":
			err = z.StateHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "TxHash":
			err = z.TxHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "ReceiptHash":
			err = z.ReceiptHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Bloom":
			err = z.Bloom.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Number":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				err = z.Number.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Time":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				err = z.Time.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Extra":
			z.Extra, err = dc.ReadBytes(z.Extra)
			if err != nil {
				return
			}
		case "MixHash":
			err = z.MixHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Nonce":
			{
				var zb0002 interface{}
				zb0002, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Nonce, err = fromBytes(zb0002)
			}
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
func (z *Header) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 11
	// write "ParentHash"
	err = en.Append(0x8b, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ParentHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Coinbase"
	err = en.Append(0xa8, 0x43, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65)
	if err != nil {
		return
	}
	err = z.Coinbase.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "StateHash"
	err = en.Append(0xa9, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.StateHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "TxHash"
	err = en.Append(0xa6, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.TxHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "ReceiptHash"
	err = en.Append(0xab, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ReceiptHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Bloom"
	err = en.Append(0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	if err != nil {
		return
	}
	err = z.Bloom.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Number"
	err = en.Append(0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	if z.Number == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Number.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Time"
	err = en.Append(0xa4, 0x54, 0x69, 0x6d, 0x65)
	if err != nil {
		return
	}
	if z.Time == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Time.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Extra"
	err = en.Append(0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Extra)
	if err != nil {
		return
	}
	// write "MixHash"
	err = en.Append(0xa7, 0x4d, 0x69, 0x78, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.MixHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Nonce"
	err = en.Append(0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	var zb0001 interface{}
	zb0001, err = toBytes(z.Nonce)
	if err != nil {
		return
	}
	err = en.WriteIntf(zb0001)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Header) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 11
	// string "ParentHash"
	o = append(o, 0x8b, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ParentHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Coinbase"
	o = append(o, 0xa8, 0x43, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65)
	o, err = z.Coinbase.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "StateHash"
	o = append(o, 0xa9, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68)
	o, err = z.StateHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "TxHash"
	o = append(o, 0xa6, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68)
	o, err = z.TxHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "ReceiptHash"
	o = append(o, 0xab, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ReceiptHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Bloom"
	o = append(o, 0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	o, err = z.Bloom.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Number"
	o = append(o, 0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if z.Number == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Number.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Time"
	o = append(o, 0xa4, 0x54, 0x69, 0x6d, 0x65)
	if z.Time == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Time.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Extra"
	o = append(o, 0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Extra)
	// string "MixHash"
	o = append(o, 0xa7, 0x4d, 0x69, 0x78, 0x48, 0x61, 0x73, 0x68)
	o, err = z.MixHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Nonce"
	o = append(o, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	var zb0001 interface{}
	zb0001, err = toBytes(z.Nonce)
	if err != nil {
		return
	}
	o, err = msgp.AppendIntf(o, zb0001)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Header) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ParentHash":
			bts, err = z.ParentHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Coinbase":
			bts, err = z.Coinbase.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "StateHash":
			bts, err = z.StateHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "TxHash":
			bts, err = z.TxHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "ReceiptHash":
			bts, err = z.ReceiptHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Bloom":
			bts, err = z.Bloom.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Number":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				bts, err = z.Number.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Time":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				bts, err = z.Time.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Extra":
			z.Extra, bts, err = msgp.ReadBytesBytes(bts, z.Extra)
			if err != nil {
				return
			}
		case "MixHash":
			bts, err = z.MixHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Nonce":
			{
				var zb0002 interface{}
				zb0002, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Nonce, err = fromBytes(zb0002)
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
func (z *Header) Msgsize() (s int) {
	s = 1 + 11 + z.ParentHash.Msgsize() + 9 + z.Coinbase.Msgsize() + 10 + z.StateHash.Msgsize() + 7 + z.TxHash.Msgsize() + 12 + z.ReceiptHash.Msgsize() + 6 + z.Bloom.Msgsize() + 7
	if z.Number == nil {
		s += msgp.NilSize
	} else {
		s += z.Number.Msgsize()
	}
	s += 5
	if z.Time == nil {
		s += msgp.NilSize
	} else {
		s += z.Time.Msgsize()
	}
	s += 6 + msgp.BytesPrefixSize + len(z.Extra) + 8 + z.MixHash.Msgsize() + 6
	var zb0001 interface{}
	_ = z.Nonce
	s += msgp.GuessSize(zb0001)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *HeaderNoNonce) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ParentHash":
			err = z.ParentHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Coinbase":
			err = z.Coinbase.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "StateHash":
			err = z.StateHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "TxHash":
			err = z.TxHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "ReceiptHash":
			err = z.ReceiptHash.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Bloom":
			err = z.Bloom.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "Number":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				err = z.Number.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Time":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				err = z.Time.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Extra":
			z.Extra, err = dc.ReadBytes(z.Extra)
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
func (z *HeaderNoNonce) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 9
	// write "ParentHash"
	err = en.Append(0x89, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ParentHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Coinbase"
	err = en.Append(0xa8, 0x43, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65)
	if err != nil {
		return
	}
	err = z.Coinbase.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "StateHash"
	err = en.Append(0xa9, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.StateHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "TxHash"
	err = en.Append(0xa6, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.TxHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "ReceiptHash"
	err = en.Append(0xab, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x48, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = z.ReceiptHash.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Bloom"
	err = en.Append(0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	if err != nil {
		return
	}
	err = z.Bloom.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Number"
	err = en.Append(0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	if z.Number == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Number.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Time"
	err = en.Append(0xa4, 0x54, 0x69, 0x6d, 0x65)
	if err != nil {
		return
	}
	if z.Time == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Time.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Extra"
	err = en.Append(0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Extra)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HeaderNoNonce) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "ParentHash"
	o = append(o, 0x89, 0xaa, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ParentHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Coinbase"
	o = append(o, 0xa8, 0x43, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65)
	o, err = z.Coinbase.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "StateHash"
	o = append(o, 0xa9, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68)
	o, err = z.StateHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "TxHash"
	o = append(o, 0xa6, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68)
	o, err = z.TxHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "ReceiptHash"
	o = append(o, 0xab, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x48, 0x61, 0x73, 0x68)
	o, err = z.ReceiptHash.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Bloom"
	o = append(o, 0xa5, 0x42, 0x6c, 0x6f, 0x6f, 0x6d)
	o, err = z.Bloom.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Number"
	o = append(o, 0xa6, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if z.Number == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Number.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Time"
	o = append(o, 0xa4, 0x54, 0x69, 0x6d, 0x65)
	if z.Time == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Time.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Extra"
	o = append(o, 0xa5, 0x45, 0x78, 0x74, 0x72, 0x61)
	o = msgp.AppendBytes(o, z.Extra)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HeaderNoNonce) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ParentHash":
			bts, err = z.ParentHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Coinbase":
			bts, err = z.Coinbase.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "StateHash":
			bts, err = z.StateHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "TxHash":
			bts, err = z.TxHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "ReceiptHash":
			bts, err = z.ReceiptHash.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Bloom":
			bts, err = z.Bloom.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "Number":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Number = nil
			} else {
				if z.Number == nil {
					z.Number = new(types.BigInt)
				}
				bts, err = z.Number.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Time":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(types.BigInt)
				}
				bts, err = z.Time.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Extra":
			z.Extra, bts, err = msgp.ReadBytesBytes(bts, z.Extra)
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
func (z *HeaderNoNonce) Msgsize() (s int) {
	s = 1 + 11 + z.ParentHash.Msgsize() + 9 + z.Coinbase.Msgsize() + 10 + z.StateHash.Msgsize() + 7 + z.TxHash.Msgsize() + 12 + z.ReceiptHash.Msgsize() + 6 + z.Bloom.Msgsize() + 7
	if z.Number == nil {
		s += msgp.NilSize
	} else {
		s += z.Number.Msgsize()
	}
	s += 5
	if z.Time == nil {
		s += msgp.NilSize
	} else {
		s += z.Time.Msgsize()
	}
	s += 6 + msgp.BytesPrefixSize + len(z.Extra)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Headers) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Headers":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Headers) >= int(zb0002) {
				z.Headers = (z.Headers)[:zb0002]
			} else {
				z.Headers = make([]*Header, zb0002)
			}
			for za0001 := range z.Headers {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						return
					}
					z.Headers[za0001] = nil
				} else {
					if z.Headers[za0001] == nil {
						z.Headers[za0001] = new(Header)
					}
					err = z.Headers[za0001].DecodeMsg(dc)
					if err != nil {
						return
					}
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
func (z *Headers) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Headers"
	err = en.Append(0x81, 0xa7, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Headers)))
	if err != nil {
		return
	}
	for za0001 := range z.Headers {
		if z.Headers[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Headers[za0001].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Headers) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Headers"
	o = append(o, 0x81, 0xa7, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Headers)))
	for za0001 := range z.Headers {
		if z.Headers[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Headers[za0001].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Headers) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Headers":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Headers) >= int(zb0002) {
				z.Headers = (z.Headers)[:zb0002]
			} else {
				z.Headers = make([]*Header, zb0002)
			}
			for za0001 := range z.Headers {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.Headers[za0001] = nil
				} else {
					if z.Headers[za0001] == nil {
						z.Headers[za0001] = new(Header)
					}
					bts, err = z.Headers[za0001].UnmarshalMsg(bts)
					if err != nil {
						return
					}
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
func (z *Headers) Msgsize() (s int) {
	s = 1 + 8 + msgp.ArrayHeaderSize
	for za0001 := range z.Headers {
		if z.Headers[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Headers[za0001].Msgsize()
		}
	}
	return
}
