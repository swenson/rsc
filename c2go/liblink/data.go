package main

import (
	"log"
	"math"
)

// Derived from Inferno utils/6l/obj.c and utils/6l/span.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/obj.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/span.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
func mangle(file string) {
	log.Fatalf("%s: mangled input file", file)
}

func symgrow(ctxt *Link, s *LSym, lsiz int64) {
	var siz int
	siz = int(lsiz)
	if int64(siz) != lsiz {
		sysfatal("symgrow size %d too long", lsiz)
	}
	if len(s.p) >= siz {
		return
	}
	for cap(s.p) < siz {
		s.p = append(s.p[:cap(s.p)], 0)
	}
	s.p = s.p[:siz]
}

func savedata(ctxt *Link, s *LSym, p *Prog, pn string) {
	var off int
	var siz int32
	var i int32
	var o int64
	var r *Reloc
	off = int(p.from.offset)
	siz = int32(ctxt.arch.datasize(p))
	if off < 0 || siz < 0 || off >= 1<<30 || siz >= 100 {
		mangle(pn)
	}
	symgrow(ctxt, s, int64(off)+int64(siz))
	if p.to.typ == ctxt.arch.D_FCONST {
		switch siz {
		default:
		case 4:
			ctxt.arch.byteOrder.PutUint32(s.p[off:], math.Float32bits(float32(p.to.u.dval)))
		case 8:
			ctxt.arch.byteOrder.PutUint64(s.p[off:], math.Float64bits(p.to.u.dval))
		}
	} else {
		if p.to.typ == ctxt.arch.D_SCONST {
			for i = 0; i < siz; i++ {
				s.p[int32(off)+i] = uint8(p.to.u.sval[i])
			}
		} else {
			if p.to.typ == ctxt.arch.D_CONST {
				if p.to.sym != nil {
					r = addrel(s)
					r.off = int64(off)
					r.siz = int64(siz)
					r.sym = p.to.sym
					r.typ = R_ADDR
					r.add = p.to.offset
					goto out
				}
				o = p.to.offset
				switch siz {
				default:
					ctxt.diag("bad nuxi %d\n%v", siz, p)
					break
				case 1:
					s.p[off] = byte(o)
				case 2:
					ctxt.arch.byteOrder.PutUint16(s.p[off:], uint16(o))
				case 4:
					ctxt.arch.byteOrder.PutUint32(s.p[off:], uint32(o))
				case 8:
					ctxt.arch.byteOrder.PutUint64(s.p[off:], uint64(o))
				}
			} else {
				if p.to.typ == ctxt.arch.D_ADDR {
					r = addrel(s)
					r.off = int64(off)
					r.siz = int64(siz)
					r.sym = p.to.sym
					r.typ = R_ADDR
					r.add = p.to.offset
				} else {
					ctxt.diag("bad data: %v", p)
				}
			}
		out:
		}
	}
}

func addrel(s *LSym) *Reloc {
	s.r = append(s.r, Reloc{})
	return &s.r[len(s.r)-1]
}

func setuintxx(ctxt *Link, s *LSym, off int64, v uint64, wid int64) int64 {
	if s.typ == 0 {
		s.typ = SDATA
	}
	s.reachable = 1
	if int64(s.size) < off+wid {
		s.size = off + wid
		symgrow(ctxt, s, int64(s.size))
	}
	switch wid {
	case 1:
		s.p[off] = uint8(v)
		break
	case 2:
		ctxt.arch.byteOrder.PutUint16(s.p[off:], uint16(v))
	case 4:
		ctxt.arch.byteOrder.PutUint32(s.p[off:], uint32(v))
	case 8:
		ctxt.arch.byteOrder.PutUint64(s.p[off:], uint64(v))
	}
	return off + wid
}

func adduintxx(ctxt *Link, s *LSym, v uint64, wid int64) int64 {
	var off int64
	off = s.size
	setuintxx(ctxt, s, off, v, wid)
	return off
}

func adduint8(ctxt *Link, s *LSym, v uint64) int64 {
	return adduintxx(ctxt, s, v, 1)
}

func adduint16(ctxt *Link, s *LSym, v uint64) int64 {
	return adduintxx(ctxt, s, v, 2)
}

func adduint32(ctxt *Link, s *LSym, v uint64) int64 {
	return adduintxx(ctxt, s, v, 4)
}

func adduint64(ctxt *Link, s *LSym, v uint64) int64 {
	return adduintxx(ctxt, s, v, 8)
}

func setuint8(ctxt *Link, s *LSym, r int64, v uint64) int64 {
	return setuintxx(ctxt, s, r, v, 1)
}

func setuint16(ctxt *Link, s *LSym, r int64, v uint64) int64 {
	return setuintxx(ctxt, s, r, v, 2)
}

func setuint32(ctxt *Link, s *LSym, r int64, v uint64) int64 {
	return setuintxx(ctxt, s, r, v, 4)
}

func setuint64(ctxt *Link, s *LSym, r int64, v uint64) int64 {
	return setuintxx(ctxt, s, r, v, 8)
}

func addaddrplus(ctxt *Link, s *LSym, t *LSym, add int64) int64 {
	var i int64
	var r *Reloc
	if s.typ == 0 {
		s.typ = SDATA
	}
	s.reachable = 1
	i = s.size
	s.size += ctxt.arch.ptrsize
	symgrow(ctxt, s, s.size)
	r = addrel(s)
	r.sym = t
	r.off = i
	r.siz = ctxt.arch.ptrsize
	r.typ = R_ADDR
	r.add = add
	return i + r.siz
}

func addpcrelplus(ctxt *Link, s *LSym, t *LSym, add int64) int64 {
	var i int64
	var r *Reloc
	if s.typ == 0 {
		s.typ = SDATA
	}
	s.reachable = 1
	i = s.size
	s.size += 4
	symgrow(ctxt, s, s.size)
	r = addrel(s)
	r.sym = t
	r.off = i
	r.add = add
	r.typ = R_PCREL
	r.siz = 4
	return i + r.siz
}

func addaddr(ctxt *Link, s *LSym, t *LSym) int64 {
	return addaddrplus(ctxt, s, t, 0)
}

func setaddrplus(ctxt *Link, s *LSym, off int64, t *LSym, add int64) int64 {
	var r *Reloc
	if s.typ == 0 {
		s.typ = SDATA
	}
	s.reachable = 1
	if off+ctxt.arch.ptrsize > s.size {
		s.size = off + ctxt.arch.ptrsize
		symgrow(ctxt, s, s.size)
	}
	r = addrel(s)
	r.sym = t
	r.off = off
	r.siz = ctxt.arch.ptrsize
	r.typ = R_ADDR
	r.add = add
	return off + r.siz
}

func setaddr(ctxt *Link, s *LSym, off int64, t *LSym) int64 {
	return setaddrplus(ctxt, s, off, t, 0)
}

func addsize(ctxt *Link, s *LSym, t *LSym) int64 {
	var i int64
	var r *Reloc
	if s.typ == 0 {
		s.typ = SDATA
	}
	s.reachable = 1
	i = s.size
	s.size += ctxt.arch.ptrsize
	symgrow(ctxt, s, s.size)
	r = addrel(s)
	r.sym = t
	r.off = i
	r.siz = ctxt.arch.ptrsize
	r.typ = R_SIZE
	return i + r.siz
}

func addaddrplus4(ctxt *Link, s *LSym, t *LSym, add int64) int64 {
	var i int64
	var r *Reloc
	if s.typ == 0 {
		s.typ = SDATA
	}
	s.reachable = 1
	i = s.size
	s.size += 4
	symgrow(ctxt, s, s.size)
	r = addrel(s)
	r.sym = t
	r.off = i
	r.siz = 4
	r.typ = R_ADDR
	r.add = add
	return i + r.siz
}
