package genesis

import (
	"fmt"
	"os"
	"mjoy.io/log"
)

var (
	logTag = "core.genesis"
	logger log.Logger
)



func init() {
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
}