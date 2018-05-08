package types

import (
	"fmt"
	"mjoy.io/log"
	"github.com/tinylib/msgp/msgp"
	"os"
	"sync"
)

var (
	globalTypeIndex = int8(10)
	mu              sync.Mutex

	logTag = "common.types"
	logger log.Logger
)

func init() {
	// get a logger
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}

	// Registering an extension is as simple as matching the
	// appropriate type number with a function that initializes
	// a freshly-allocated object of that type
	RegisterExtension(&hashType, func() msgp.Extension { return new(Hash) })
	RegisterExtension(&addressType, func() msgp.Extension { return new(Address) })
	RegisterExtension(&ipType, func() msgp.Extension { return new(IP) })
	RegisterExtension(&bigIntType, func() msgp.Extension { return new(BigInt) })
	RegisterExtension(&bloomType, func() msgp.Extension { return new(Bloom) })
}

type decError struct{ msg string }

func (err decError) Error() string { return err.msg }

var (
	ErrBytesTooLong = &decError{"bytes too long"}

	ErrRegisterFull    = &decError{"Can't register more type"}
	ErrRegisterFailure = &decError{"Register is failure"}
)

func registerExtension(typ *int8, f func() msgp.Extension) error {
	mu.Lock()
	defer func() (err error) {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %v", p)
			err = ErrRegisterFailure
		}

		mu.Unlock()
		return err
	}()

	if globalTypeIndex == -128 {
		return ErrRegisterFull
	}
	msgp.RegisterExtension(globalTypeIndex, f)
	*typ = globalTypeIndex
	globalTypeIndex++

	return nil
}

func RegisterExtension(typ *int8, f func() msgp.Extension) error {
	for {
		err := registerExtension(typ, f)
		if err == ErrRegisterFailure {
			continue
		} else {
			return err
		}
	}
}
