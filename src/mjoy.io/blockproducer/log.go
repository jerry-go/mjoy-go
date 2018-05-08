package blockproducer

import (
	"fmt"
	"os"
	"mjoy.io/log"
)

var (
	LogTag = "blockproducer"
	logger log.Logger
)

func init() {
	logger = log.GetLogger(LogTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", LogTag)
		os.Exit(1)
	}
}
