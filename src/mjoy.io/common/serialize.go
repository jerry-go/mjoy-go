package common

import "io"

type Serializer interface {
	Serialize(w io.Writer) error
}

type UnSerializer interface {
	UnSerialize(stream interface{}) error
}
