// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Derived from FUSE's fuse_kernel.h
/*
   This file defines the kernel interface of FUSE
   Copyright (C) 2001-2007  Miklos Szeredi <miklos@szeredi.hu>


   This -- and only this -- header file may also be distributed under
   the terms of the BSD Licence as follows:

   Copyright (C) 2001-2007 Miklos Szeredi. All rights reserved.

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:
   1. Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
   2. Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.

   THIS SOFTWARE IS PROVIDED BY AUTHOR AND CONTRIBUTORS ``AS IS'' AND
   ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
   IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
   ARE DISCLAIMED.  IN NO EVENT SHALL AUTHOR OR CONTRIBUTORS BE LIABLE
   FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
   DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
   OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
   HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
   LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
   OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
   SUCH DAMAGE.
*/

package fuse

import (
	"fmt"
	"unsafe"
)

// Version is the FUSE version implemented by the package.
const Version = "7.8"

const (
	kernelVersion      = 7
	kernelMinorVersion = 8
	rootID             = 1
)

type attr struct {
	Ino        uint64
	Size       uint64
	Blocks     uint64
	Atime      uint64
	Mtime      uint64
	Ctime      uint64
	Crtime     uint64 // OS X only
	AtimeNsec  uint32
	MtimeNsec  uint32
	CtimeNsec  uint32
	CrtimeNsec uint32 // OS X only
	Mode       uint32
	Nlink      uint32
	Uid        uint32
	Gid        uint32
	Rdev       uint32
	Flags      uint32 // OS X only; see chflags(2)
}

type kstatfs struct {
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Bsize   uint32
	Namelen uint32
	Frsize  uint32
	Padding uint32
	Spare   [6]uint32
}

type fileLock struct {
	Start uint64
	End   uint64
	Type  uint32
	Pid   uint32
}

// Bitmasks for fuse_setattr_in.valid
const (
	fattrMode  = 1 << 0
	fattrUid   = 1 << 1
	fattrGid   = 1 << 2
	fattrSize  = 1 << 3
	fattrAtime = 1 << 4
	fattrMtime = 1 << 5
	fattrFh    = 1 << 6

	// OS X only
	fattrCrtime   = 1 << 28
	fattrChgtime  = 1 << 29
	fattrBkuptime = 1 << 30
	fattrFlags    = 1 << 31
)

// The OpenFlags are returned in the OpenResponse.
type OpenFlags uint32

const (
	OpenDirectIO  OpenFlags = 1 << 0 // bypass page cache for this open file
	OpenKeepCache OpenFlags = 1 << 1 // don't invalidate the data cache on open

	OpenPurgeAttr OpenFlags = 1 << 30 // OS X
	OpenPurgeUBC  OpenFlags = 1 << 31 // OS X
)

func (fl OpenFlags) String() string {
	return flagString(uint32(fl), openFlagNames)
}

var openFlagNames = []flagName{
	{uint32(OpenDirectIO), "OpenDirectIO"},
	{uint32(OpenKeepCache), "OpenKeepCache"},
	{uint32(OpenPurgeAttr), "OpenPurgeAttr"},
	{uint32(OpenPurgeUBC), "OpenPurgeUBC"},
}

// The InitFlags are used in the Init exchange.
type InitFlags uint32

const (
	InitAsyncRead  InitFlags = 1 << 0
	InitPosixLocks InitFlags = 1 << 1

	InitCaseSensitive InitFlags = 1 << 29 // OS X only
	InitVolRename     InitFlags = 1 << 30 // OS X only
	InitXtimes        InitFlags = 1 << 31 // OS X only
)

type flagName struct {
	bit  uint32
	name string
}

var initFlagNames = []flagName{
	{uint32(InitAsyncRead), "InitAsyncRead"},
	{uint32(InitPosixLocks), "InitPosixLocks"},
	{uint32(InitCaseSensitive), "InitCaseSensitive"},
	{uint32(InitVolRename), "InitVolRename"},
	{uint32(InitXtimes), "InitXtimes"},
}

func (fl InitFlags) String() string {
	return flagString(uint32(fl), initFlagNames)
}

func flagString(f uint32, names []flagName) string {
	var s string

	if f == 0 {
		return "0"
	}

	for _, n := range names {
		if f&n.bit != 0 {
			s += "+" + n.name
			f &^= n.bit
		}
	}
	if f != 0 {
		s += fmt.Sprintf("%+#x", f)
	}
	return s[1:]
}

// The ReleaseFlags are used in the Release exchange.
type ReleaseFlags uint32

const (
	ReleaseFlush ReleaseFlags = 1 << 0
)

func (fl ReleaseFlags) String() string {
	return flagString(uint32(fl), releaseFlagNames)
}

var releaseFlagNames = []flagName{
	{uint32(ReleaseFlush), "ReleaseFlush"},
}

// Opcodes
const (
	opLookup      = 1
	opForget      = 2 // no reply
	opGetattr     = 3
	opSetattr     = 4
	opReadlink    = 5
	opSymlink     = 6
	opMknod       = 8
	opMkdir       = 9
	opUnlink      = 10
	opRmdir       = 11
	opRename      = 12
	opLink        = 13
	opOpen        = 14
	opRead        = 15
	opWrite       = 16
	opStatfs      = 17
	opRelease     = 18
	opFsync       = 20
	opSetxattr    = 21
	opGetxattr    = 22
	opListxattr   = 23
	opRemovexattr = 24
	opFlush       = 25
	opInit        = 26
	opOpendir     = 27
	opReaddir     = 28
	opReleasedir  = 29
	opFsyncdir    = 30
	opGetlk       = 31
	opSetlk       = 32
	opSetlkw      = 33
	opAccess      = 34
	opCreate      = 35
	opInterrupt   = 36
	opBmap        = 37
	opDestroy     = 38

	// OS X
	opSetvolname = 61
	opGetxtimes  = 62
	opExchange   = 63
)

// The read buffer is required to be at least 8k but may be much larger
const minReadBuffer = 8192

type entryOut struct {
	outHeader
	Nodeid         uint64 // Inode ID
	Generation     uint64 // Inode generation
	EntryValid     uint64 // Cache timeout for the name
	AttrValid      uint64 // Cache timeout for the attributes
	EntryValidNsec uint32
	AttrValidNsec  uint32
	Attr           attr
}

type forgetIn struct {
	Nlookup uint64
}

type attrOut struct {
	outHeader
	AttrValid     uint64 // Cache timeout for the attributes
	AttrValidNsec uint32
	Dummy         uint32
	Attr          attr
}

// OS X
type getxtimesOut struct {
	outHeader
	Bkuptime     uint64
	Crtime       uint64
	BkuptimeNsec uint32
	CrtimeNsec   uint32
}

type mknodIn struct {
	Mode uint32
	Rdev uint32
}

type mkdirIn struct {
	Mode    uint32
	Padding uint32
}

type renameIn struct {
	Newdir uint64
}

// OS X 
type exchangeIn struct {
	Olddir  uint64
	Newdir  uint64
	Options uint64
}

type linkIn struct {
	Oldnodeid uint64
}
type setattrIn struct {
	Valid     uint32
	Padding   uint32
	Fh        uint64
	Size      uint64
	Unused1   uint64
	Atime     uint64
	Mtime     uint64
	Unused2   uint64
	AtimeNsec uint32
	MtimeNsec uint32
	Unused3   uint32
	Mode      uint32
	Unused4   uint32
	Uid       uint32
	Gid       uint32
	Unused5   uint32

	// OS X only
	Bkuptime     uint64
	Chgtime      uint64
	Crtime       uint64
	BkuptimeNsec uint32
	ChgtimeNsec  uint32
	CrtimeNsec   uint32
	Flags        uint32 // see chflags(2)
}

type openIn struct {
	Flags uint32
	Mode  uint32
}

type openOut struct {
	outHeader
	Fh        uint64
	OpenFlags uint32
	Padding   uint32
}

type releaseIn struct {
	Fh           uint64
	Flags        uint32
	ReleaseFlags uint32
	LockOwner    uint32
}

type flushIn struct {
	Fh         uint64
	FlushFlags uint32
	Padding    uint32
	LockOwner  uint64
}

type readIn struct {
	Fh      uint64
	Offset  uint64
	Size    uint32
	Padding uint32
}

type writeIn struct {
	Fh         uint64
	Offset     uint64
	Size       uint32
	WriteFlags uint32
}

type writeOut struct {
	outHeader
	Size    uint32
	Padding uint32
}

const compatStatfsSize = 48

type statfsOut struct {
	outHeader
	St kstatfs
}

type fsyncIn struct {
	Fh         uint64
	FsyncFlags uint32
	Padding    uint32
}

type setxattrIn struct {
	Size  uint32
	Flags uint32
}

type setxattrInOSX struct {
	Size  uint32
	Flags uint32

	// OS X only
	Position uint32
	Padding  uint32
}

type getxattrIn struct {
	Size    uint32
	Padding uint32
}

type getxattrInOSX struct {
	Size    uint32
	Padding uint32

	// OS X only
	Position uint32
	Padding2 uint32
}

type getxattrOut struct {
	outHeader
	Size    uint32
	Padding uint32
}

type lkIn struct {
	Fh    uint64
	Owner uint64
	Lk    fileLock
}

type lkOut struct {
	outHeader
	Lk fileLock
}

type accessIn struct {
	Mask    uint32
	Padding uint32
}

type initIn struct {
	Major        uint32
	Minor        uint32
	MaxReadahead uint32
	Flags        uint32
}

const initInSize = int(unsafe.Sizeof(initIn{}))

type initOut struct {
	outHeader
	Major        uint32
	Minor        uint32
	MaxReadahead uint32
	Flags        uint32
	Unused       uint32
	MaxWrite     uint32
}

type interruptIn struct {
	Unique uint64
}

type bmapIn struct {
	Block     uint64
	BlockSize uint32
	Padding   uint32
}

type bmapOut struct {
	outHeader
	Block uint64
}

type inHeader struct {
	Len     uint32
	Opcode  uint32
	Unique  uint64
	Nodeid  uint64
	Uid     uint32
	Gid     uint32
	Pid     uint32
	Padding uint32
}

const inHeaderSize = int(unsafe.Sizeof(inHeader{}))

type outHeader struct {
	Len    uint32
	Error  int32
	Unique uint64
}

type dirent struct {
	Ino     uint64
	Off     uint64
	Namelen uint32
	Type    uint32
	Name    [0]byte
}

const direntSize = 8 + 8 + 4 + 4
