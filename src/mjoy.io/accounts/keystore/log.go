////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software.
//
// @File: log.go
// @Date: 2018/05/08 17:10:08
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  <http://www.apache.org/licenses/>
//
////////////////////////////////////////////////////////////////////////////////

package keystore

import (
	"fmt"
	"os"
	"mjoy.io/log"
)

var (
	LogTag = "accounts.keystore"
	logger log.Logger
)

func init() {
	logger = log.GetLogger(LogTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", LogTag)
		os.Exit(1)
	}
}

