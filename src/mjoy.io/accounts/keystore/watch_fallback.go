////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software.
//
// @File: watch_fallback.go
// @Date: 2018/05/08 17:13:08
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  <http://www.apache.org/licenses/>
//
////////////////////////////////////////////////////////////////////////////////

// It is used on unsupported platforms.

// +build ios linux,arm64 windows !darwin,!freebsd,!linux,!netbsd,!solaris

package keystore

type watcher struct{ running bool }

func newWatcher(*accountCache) *watcher { return new(watcher) }
func (*watcher) start()                 {}
func (*watcher) close()                 {}
