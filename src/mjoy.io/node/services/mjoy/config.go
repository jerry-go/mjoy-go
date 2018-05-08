
package mjoy

import (
	"os"
	"os/user"


	"mjoy.io/common/types"
	"mjoy.io/common/types/util/hex"
	"mjoy.io/node/services/mjoy/downloader"
	"mjoy.io/core/txprocessor"
	"mjoy.io/core/genesis"
	"mjoy.io/node"
)

// DefaultConfig contains default settings for use on the Mjoy main net.
var DefaultConfig = Config{
	SyncMode: downloader.FastSync,

	NetworkId:     1,
	LightPeers:    20,
	DatabaseCache: 128,

	TxPool: txprocessor.DefaultTxPoolConfig,

}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}

}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Mjoy main net block is used.
	Genesis *genesis.Genesis `toml:",omitempty"`

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to
	SyncMode  downloader.SyncMode

	// Light client options
	LightServ  int `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightPeers int `toml:",omitempty"` // Maximum number of LES client peers

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int

	// Producing-related options
	Coinbase    types.Address `toml:",omitempty"`
	BlockproducerThreads int  `toml:",omitempty"`
	ExtraData    []byte       `toml:",omitempty"`


	// Transaction pool options
	TxPool txprocessor.TxPoolConfig



	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`

	//should we start blockproducer at first
	StartBlockproducerAtStart bool
}

type configMarshaling struct {
	ExtraData hex.Bytes
}


type MjoydConfig struct{
	Mjoy Config
	Node node.Config
}