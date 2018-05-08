package mjoy

import (
	"fmt"
	"os"
	"mjoy.io/log"
)

var (
	logTag = "node.services.mjoy"
	logger log.Logger
)



func init() {
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
}
