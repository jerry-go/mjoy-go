package types

import (
	"net"
)

//go:generate msgp
//msgp:shim net.IP as:[]byte using:toBytes/fromBytes

func toBytes(ip net.IP) []byte {
	return []byte(ip)
}

func fromBytes(b []byte) net.IP {
	return net.IP(b)
}

var (
	ipType int8
)

type IP struct {
	Ip net.IP `msg:"ip"`
}

func (ip IP) Get() net.IP {
	return ip.Ip
}

func (ip *IP) Put(in net.IP) *IP {
	ip.Ip = in
	return ip
}

func NewIP(in net.IP) *IP {
	ip := new(IP)
	ip.Ip = in
	return ip
}

// Here, we'll pick an arbitrary number between
// 0 and 127 that isn't already in use
func (*IP) ExtensionType() int8 {
	return ipType
}

// We'll always use 16 bytes to encode the data
func (*IP) Len() int {
	return net.IPv6len
}

// MarshalBinaryTo simply copies the value
// of the bytes into 'b'
func (ip *IP) MarshalBinaryTo(b []byte) error {
	copy(b, ip.Ip)
	return nil
}

// UnmarshalBinary copies the value of 'b'
// into the Hash object. (We might want to add
// a sanity check here later that len(b) <= HashLength.)
func (ip *IP) UnmarshalBinary(b []byte) error {
	// TODO: check b, only hex, len <= HashLength
	if len(b) <= net.IPv6len {
		copy(ip.Ip, b)
		return nil
	}

	return ErrBytesTooLong
}
