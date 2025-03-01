// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"net"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Kuraaa/net/internal/iana"
	"github.com/Kuraaa/net/internal/socket"

	"golang.org/x/sys/unix"
)

var (
	ctlOpts = [ctlMax]ctlOpt{
		ctlTrafficClass: {unix.IPV6_TCLASS, 4, marshalTrafficClass, parseTrafficClass},
		ctlHopLimit:     {unix.IPV6_HOPLIMIT, 4, marshalHopLimit, parseHopLimit},
		ctlPacketInfo:   {unix.IPV6_PKTINFO, sizeofInet6Pktinfo, marshalPacketInfo, parsePacketInfo},
		ctlNextHop:      {unix.IPV6_NEXTHOP, sizeofSockaddrInet6, marshalNextHop, parseNextHop},
		ctlPathMTU:      {unix.IPV6_PATHMTU, sizeofIPv6Mtuinfo, marshalPathMTU, parsePathMTU},
	}

	sockOpts = map[int]sockOpt{
		ssoTrafficClass:        {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_TCLASS, Len: 4}},
		ssoHopLimit:            {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_UNICAST_HOPS, Len: 4}},
		ssoMulticastInterface:  {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_MULTICAST_IF, Len: 4}},
		ssoMulticastHopLimit:   {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_MULTICAST_HOPS, Len: 4}},
		ssoMulticastLoopback:   {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_MULTICAST_LOOP, Len: 4}},
		ssoReceiveTrafficClass: {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_RECVTCLASS, Len: 4}},
		ssoReceiveHopLimit:     {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_RECVHOPLIMIT, Len: 4}},
		ssoReceivePacketInfo:   {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_RECVPKTINFO, Len: 4}},
		ssoReceivePathMTU:      {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_RECVPATHMTU, Len: 4}},
		ssoPathMTU:             {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_PATHMTU, Len: sizeofIPv6Mtuinfo}},
		ssoChecksum:            {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.IPV6_CHECKSUM, Len: 4}},
		ssoICMPFilter:          {Option: socket.Option{Level: iana.ProtocolIPv6ICMP, Name: unix.ICMP6_FILTER, Len: sizeofICMPv6Filter}},
		ssoJoinGroup:           {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.MCAST_JOIN_GROUP, Len: sizeofGroupReq}, typ: ssoTypeGroupReq},
		ssoLeaveGroup:          {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.MCAST_LEAVE_GROUP, Len: sizeofGroupReq}, typ: ssoTypeGroupReq},
		ssoJoinSourceGroup:     {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.MCAST_JOIN_SOURCE_GROUP, Len: sizeofGroupSourceReq}, typ: ssoTypeGroupSourceReq},
		ssoLeaveSourceGroup:    {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.MCAST_LEAVE_SOURCE_GROUP, Len: sizeofGroupSourceReq}, typ: ssoTypeGroupSourceReq},
		ssoBlockSourceGroup:    {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.MCAST_BLOCK_SOURCE, Len: sizeofGroupSourceReq}, typ: ssoTypeGroupSourceReq},
		ssoUnblockSourceGroup:  {Option: socket.Option{Level: iana.ProtocolIPv6, Name: unix.MCAST_UNBLOCK_SOURCE, Len: sizeofGroupSourceReq}, typ: ssoTypeGroupSourceReq},
	}
)

func init() {
	if runtime.GOOS == "freebsd" && runtime.GOARCH == "386" {
		archs, _ := syscall.Sysctl("kern.supported_archs")
		for _, s := range strings.Fields(archs) {
			if s == "amd64" {
				compatFreeBSD32 = true
				break
			}
		}
	}
}

func (sa *sockaddrInet6) setSockaddr(ip net.IP, i int) {
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], ip)
	sa.Scope_id = uint32(i)
}

func (pi *inet6Pktinfo) setIfindex(i int) {
	pi.Ifindex = uint32(i)
}

func (mreq *ipv6Mreq) setIfindex(i int) {
	mreq.Interface = uint32(i)
}

func (gr *groupReq) setGroup(grp net.IP) {
	sa := (*sockaddrInet6)(unsafe.Pointer(&gr.Group))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], grp)
}

func (gsr *groupSourceReq) setSourceGroup(grp, src net.IP) {
	sa := (*sockaddrInet6)(unsafe.Pointer(&gsr.Group))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], grp)
	sa = (*sockaddrInet6)(unsafe.Pointer(&gsr.Source))
	sa.Len = sizeofSockaddrInet6
	sa.Family = syscall.AF_INET6
	copy(sa.Addr[:], src)
}
