// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !aix && !darwin && !freebsd && !linux && !solaris && !zos
// +build !aix,!darwin,!freebsd,!linux,!solaris,!zos

package ipv6

import (
	"net"

	"github.com/Kuraaa/net/internal/socket"
)

func (so *sockOpt) setGroupReq(c *socket.Conn, ifi *net.Interface, grp net.IP) error {
	return errNotImplemented
}

func (so *sockOpt) setGroupSourceReq(c *socket.Conn, ifi *net.Interface, grp, src net.IP) error {
	return errNotImplemented
}
