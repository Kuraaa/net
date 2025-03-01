// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"net"

	"github.com/Kuraaa/net/internal/socket"
)

// BUG(mikio): On Windows, the ControlMessage for ReadFrom and WriteTo
// methods of PacketConn is not implemented.

// A payloadHandler represents the IPv4 datagram payload handler.
type payloadHandler struct {
	net.PacketConn
	*socket.Conn
	rawOpt
}

func (c *payloadHandler) ok() bool { return c != nil && c.PacketConn != nil && c.Conn != nil }
