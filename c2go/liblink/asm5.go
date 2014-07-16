package main

import (
	"math"
	"sort"
)

func span5(ctxt *Link, cursym *LSym) {
	var p *Prog
	var op *Prog
	var o *Optab_asm5
	var m int
	var bflag int
	var i int
	var v int
	var times int
	var c int
	var opc int
	var out [9]uint32
	var bp []uint8
	p = cursym.text
	if p == nil || p.link == nil { // handle external functions and ELF section symbols
		return
	}
	if oprange_asm5[AAND_5].start == nil {
		buildop_asm5(ctxt)
	}
	ctxt.cursym = cursym
	ctxt.autosize = int32(p.to.offset + 4)
	c = 0
	op = p
	p = p.link
	for ; p != nil || ctxt.blitrl != nil; (func() { op = p; p = p.link })() {
		if p == nil {
			if checkpool_asm5(ctxt, op, 0) != 0 {
				p = op
				continue
			}
			// can't happen: blitrl is not nil, but checkpool didn't flushpool
			ctxt.diag("internal inconsistency")
			break
		}
		ctxt.curp = p
		p.pc = c
		o = oplook_asm5(ctxt, p)
		if ctxt.headtype != int(Hnacl) {
			m = o.size
		} else {
			m = asmoutnacl_asm5(ctxt, int32(c), p, o, []uint32(nil))
			c = p.pc                 // asmoutnacl might change pc for alignment
			o = oplook_asm5(ctxt, p) // asmoutnacl might change p in rare cases
		}
		if m%4 != 0 || p.pc%4 != 0 {
			ctxt.diag("!pc invalid: %P size=%d", p, m)
		}
		// must check literal pool here in case p generates many instructions
		if ctxt.blitrl != nil {
			var tmp int32
			if p.as == int(ACASE_5) {
				tmp = casesz_asm5(ctxt, p)
			} else {
				tmp = int32(m)
			}
			if checkpool_asm5(ctxt, op, int(tmp)) != 0 {
				p = op
				continue
			}
		}
		if m == 0 && (p.as != int(AFUNCDATA_5) && p.as != int(APCDATA_5) && p.as != int(ADATABUNDLEEND_5)) {
			ctxt.diag("zero-width instruction\n%P", p)
			continue
		}
		switch o.flag & int(LFROM_asm5|LTO_asm5|LPOOL_asm5) {
		case LFROM_asm5:
			addpool_asm5(ctxt, p, &p.from)
			break
		case LTO_asm5:
			addpool_asm5(ctxt, p, &p.to)
			break
		case LPOOL_asm5:
			if (p.scond & int(C_SCOND_5)) == int(C_SCOND_NONE_5) {
				flushpool_asm5(ctxt, p, 0, 0)
			}
			break
		}
		if p.as == int(AMOVW_5) && p.to.typ == int(D_REG_5) && p.to.reg == int(REGPC_5) && (p.scond&int(C_SCOND_5)) == int(C_SCOND_NONE_5) {
			flushpool_asm5(ctxt, p, 0, 0)
		}
		c += m
	}
	cursym.size = c
	/*
	 * if any procedure is large enough to
	 * generate a large SBRA branch, then
	 * generate extra passes putting branches
	 * around jmps to fix. this is rare.
	 */
	times = 0
	for {
		if ctxt.debugvlog != 0 {
			Bprint(ctxt.bso, "%5.2f span1\n", cputime())
		}
		bflag = 0
		c = 0
		times++
		cursym.text.pc = 0 // force re-layout the code.
		for p = cursym.text; p != nil; p = p.link {
			ctxt.curp = p
			o = oplook_asm5(ctxt, p)
			if c > p.pc {
				p.pc = c
			}
			/* very large branches
			if(o->type == 6 && p->pcond) {
				otxt = p->pcond->pc - c;
				if(otxt < 0)
					otxt = -otxt;
				if(otxt >= (1L<<17) - 10) {
					q = ctxt->arch->prg();
					q->link = p->link;
					p->link = q;
					q->as = AB;
					q->to.type = D_BRANCH;
					q->pcond = p->pcond;
					p->pcond = q;
					q = ctxt->arch->prg();
					q->link = p->link;
					p->link = q;
					q->as = AB;
					q->to.type = D_BRANCH;
					q->pcond = q->link->link;
					bflag = 1;
				}
			}
			*/
			opc = p.pc
			if ctxt.headtype != int(Hnacl) {
				m = o.size
			} else {
				m = asmoutnacl_asm5(ctxt, int32(c), p, o, []uint32(nil))
			}
			if p.pc != opc {
				bflag = 1
			}
			//print("%P pc changed %d to %d in iter. %d\n", p, opc, (int32)p->pc, times);
			c = p.pc + m
			if m%4 != 0 || p.pc%4 != 0 {
				ctxt.diag("pc invalid: %P size=%d", p, m)
			}
			if m > len(out)*4 {
				ctxt.diag("instruction size too large: %d > %d", m, len(out)*4)
			}
			if m == 0 && (p.as != int(AFUNCDATA_5) && p.as != int(APCDATA_5) && p.as != int(ADATABUNDLEEND_5)) {
				if p.as == int(ATEXT_5) {
					ctxt.autosize = int32(p.to.offset + 4)
					continue
				}
				ctxt.diag("zero-width instruction\n%P", p)
				continue
			}
		}
		cursym.size = c
		if !(bflag != 0) {
			break
		}
	}
	if c%4 != 0 {
		ctxt.diag("sym->size=%d, invalid", c)
	}
	/*
	 * lay out the code.  all the pc-relative code references,
	 * even cross-function, are resolved now;
	 * only data references need to be relocated.
	 * with more work we could leave cross-function
	 * code references to be relocated too, and then
	 * perhaps we'd be able to parallelize the span loop above.
	 */
	if ctxt.tlsg == nil {
		ctxt.tlsg = linklookup(ctxt, "runtime.tlsg", 0)
	}
	p = cursym.text
	ctxt.autosize = int32(p.to.offset + 4)
	symgrow(ctxt, cursym, int64(cursym.size))
	bp = cursym.p
	c = p.pc // even p->link might need extra padding
	for p = p.link; p != nil; p = p.link {
		ctxt.pc = p.pc
		ctxt.curp = p
		o = oplook_asm5(ctxt, p)
		opc = p.pc
		if ctxt.headtype != int(Hnacl) {
			asmout_asm5(ctxt, p, o, out[:])
			m = o.size
		} else {
			m = asmoutnacl_asm5(ctxt, int32(c), p, o, out[:])
			if opc != p.pc {
				ctxt.diag("asmoutnacl broken: pc changed (%d->%d) in last stage: %P", opc, int32(p.pc), p)
			}
		}
		if m%4 != 0 || p.pc%4 != 0 {
			ctxt.diag("final stage: pc invalid: %P size=%d", p, m)
		}
		if c > p.pc {
			ctxt.diag("PC padding invalid: want %#lld, has %#d: %P", p.pc, c, p)
		}
		for c != p.pc {
			// emit 0xe1a00000 (MOVW R0, R0)
			bp[0] = 0x00
			bp = bp[1:]
			bp[0] = 0x00
			bp = bp[1:]
			bp[0] = 0xa0
			bp = bp[1:]
			bp[0] = 0xe1
			bp = bp[1:]
			c += 4
		}
		for i = 0; i < m/4; i++ {
			v = int(out[i])
			bp[0] = uint8(v)
			bp = bp[1:]
			bp[0] = uint8(v >> 8)
			bp = bp[1:]
			bp[0] = uint8(v >> 16)
			bp = bp[1:]
			bp[0] = uint8(v >> 24)
			bp = bp[1:]
		}
		c += m
	}
}

func chipfloat5(ctxt *Link, e float64) int {
	var n int
	var h1 uint32
	var l int32
	var h uint32
	var ei uint64
	// We use GOARM=7 to gate the use of VFPv3 vmov (imm) instructions.
	if ctxt.goarm < 7 {
		goto no
	}
	ei = math.Float64bits(e)
	l = int32(ei)
	h = uint32(ei >> 32)
	if l != 0 || (h&0xffff) != 0 {
		goto no
	}
	h1 = h & 0x7fc00000
	if h1 != 0x40000000 && h1 != 0x3fc00000 {
		goto no
	}
	n = 0
	// sign bit (a)
	if h&0x80000000 != 0 /*untyped*/ {
		n |= 1 << 7
	}
	// exp sign bit (b)
	if h1 == 0x3fc00000 {
		n |= 1 << 6
	}
	// rest of exp and mantissa (cd-efgh)
	n |= int((h >> 16) & 0x3f)
	//print("match %.8lux %.8lux %d\n", l, h, n);
	return n
no:
	return -1
}

func chipzero5(ctxt *Link, e float64) int {
	// We use GOARM=7 to gate the use of VFPv3 vmov (imm) instructions.
	if ctxt.goarm < 7 || e != 0 {
		return -1
	}
	return 0
}

type Optab_asm5 struct {
	as       int
	a1       uint8
	a2       int
	a3       uint8
	typ      uint8
	size     int
	param    int
	flag     int
	pcrelsiz uint8
}

type Oprang_asm5 struct {
	start []Optab_asm5
	stop  []Optab_asm5
}

type Opcross_asm5 [32][2][32]uint8

const (
	LFROM_asm5  = 1 << 0
	LTO_asm5    = 1 << 1
	LPOOL_asm5  = 1 << 2
	LPCREL_asm5 = 1 << 3
	C_NONE_asm5 = 0 + iota - 4
	C_REG_asm5
	C_REGREG_asm5
	C_REGREG2_asm5
	C_SHIFT_asm5
	C_FREG_asm5
	C_PSR_asm5
	C_FCR_asm5
	C_RCON_asm5
	C_NCON_asm5
	C_SCON_asm5
	C_LCON_asm5
	C_LCONADDR_asm5
	C_ZFCON_asm5
	C_SFCON_asm5
	C_LFCON_asm5
	C_RACON_asm5
	C_LACON_asm5
	C_SBRA_asm5
	C_LBRA_asm5
	C_HAUTO_asm5
	C_FAUTO_asm5
	C_HFAUTO_asm5
	C_SAUTO_asm5
	C_LAUTO_asm5
	C_HOREG_asm5
	C_FOREG_asm5
	C_HFOREG_asm5
	C_SOREG_asm5
	C_ROREG_asm5
	C_SROREG_asm5
	C_LOREG_asm5
	C_PC_asm5
	C_SP_asm5
	C_HREG_asm5
	C_ADDR_asm5
	C_GOK_asm5
)

var optab_asm5 = []Optab_asm5{
	/* struct Optab:
	OPCODE,	from, prog->reg, to,		 type,size,param,flag */
	{ATEXT_5, C_ADDR_asm5, C_NONE_asm5, C_LCON_asm5, 0, 0, 0, 0, 0},
	{ATEXT_5, C_ADDR_asm5, C_REG_asm5, C_LCON_asm5, 0, 0, 0, 0, 0},
	{AADD_5, C_REG_asm5, C_REG_asm5, C_REG_asm5, 1, 4, 0, 0, 0},
	{AADD_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 1, 4, 0, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 1, 4, 0, 0, 0},
	{AMVN_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 1, 4, 0, 0, 0},
	{ACMP_5, C_REG_asm5, C_REG_asm5, C_NONE_asm5, 1, 4, 0, 0, 0},
	{AADD_5, C_RCON_asm5, C_REG_asm5, C_REG_asm5, 2, 4, 0, 0, 0},
	{AADD_5, C_RCON_asm5, C_NONE_asm5, C_REG_asm5, 2, 4, 0, 0, 0},
	{AMOVW_5, C_RCON_asm5, C_NONE_asm5, C_REG_asm5, 2, 4, 0, 0, 0},
	{AMVN_5, C_RCON_asm5, C_NONE_asm5, C_REG_asm5, 2, 4, 0, 0, 0},
	{ACMP_5, C_RCON_asm5, C_REG_asm5, C_NONE_asm5, 2, 4, 0, 0, 0},
	{AADD_5, C_SHIFT_asm5, C_REG_asm5, C_REG_asm5, 3, 4, 0, 0, 0},
	{AADD_5, C_SHIFT_asm5, C_NONE_asm5, C_REG_asm5, 3, 4, 0, 0, 0},
	{AMVN_5, C_SHIFT_asm5, C_NONE_asm5, C_REG_asm5, 3, 4, 0, 0, 0},
	{ACMP_5, C_SHIFT_asm5, C_REG_asm5, C_NONE_asm5, 3, 4, 0, 0, 0},
	{AMOVW_5, C_RACON_asm5, C_NONE_asm5, C_REG_asm5, 4, 4, REGSP_5, 0, 0},
	{AB_5, C_NONE_asm5, C_NONE_asm5, C_SBRA_asm5, 5, 4, 0, LPOOL_asm5, 0},
	{ABL_5, C_NONE_asm5, C_NONE_asm5, C_SBRA_asm5, 5, 4, 0, 0, 0},
	{ABX_5, C_NONE_asm5, C_NONE_asm5, C_SBRA_asm5, 74, 20, 0, 0, 0},
	{ABEQ_5, C_NONE_asm5, C_NONE_asm5, C_SBRA_asm5, 5, 4, 0, 0, 0},
	{AB_5, C_NONE_asm5, C_NONE_asm5, C_ROREG_asm5, 6, 4, 0, LPOOL_asm5, 0},
	{ABL_5, C_NONE_asm5, C_NONE_asm5, C_ROREG_asm5, 7, 4, 0, 0, 0},
	{ABL_5, C_REG_asm5, C_NONE_asm5, C_ROREG_asm5, 7, 4, 0, 0, 0},
	{ABX_5, C_NONE_asm5, C_NONE_asm5, C_ROREG_asm5, 75, 12, 0, 0, 0},
	{ABXRET_5, C_NONE_asm5, C_NONE_asm5, C_ROREG_asm5, 76, 4, 0, 0, 0},
	{ASLL_5, C_RCON_asm5, C_REG_asm5, C_REG_asm5, 8, 4, 0, 0, 0},
	{ASLL_5, C_RCON_asm5, C_NONE_asm5, C_REG_asm5, 8, 4, 0, 0, 0},
	{ASLL_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 9, 4, 0, 0, 0},
	{ASLL_5, C_REG_asm5, C_REG_asm5, C_REG_asm5, 9, 4, 0, 0, 0},
	{ASWI_5, C_NONE_asm5, C_NONE_asm5, C_NONE_asm5, 10, 4, 0, 0, 0},
	{ASWI_5, C_NONE_asm5, C_NONE_asm5, C_LOREG_asm5, 10, 4, 0, 0, 0},
	{ASWI_5, C_NONE_asm5, C_NONE_asm5, C_LCON_asm5, 10, 4, 0, 0, 0},
	{AWORD_5, C_NONE_asm5, C_NONE_asm5, C_LCON_asm5, 11, 4, 0, 0, 0},
	{AWORD_5, C_NONE_asm5, C_NONE_asm5, C_LCONADDR_asm5, 11, 4, 0, 0, 0},
	{AWORD_5, C_NONE_asm5, C_NONE_asm5, C_ADDR_asm5, 11, 4, 0, 0, 0},
	{AMOVW_5, C_NCON_asm5, C_NONE_asm5, C_REG_asm5, 12, 4, 0, 0, 0},
	{AMOVW_5, C_LCON_asm5, C_NONE_asm5, C_REG_asm5, 12, 4, 0, LFROM_asm5, 0},
	{AMOVW_5, C_LCONADDR_asm5, C_NONE_asm5, C_REG_asm5, 12, 4, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AADD_5, C_NCON_asm5, C_REG_asm5, C_REG_asm5, 13, 8, 0, 0, 0},
	{AADD_5, C_NCON_asm5, C_NONE_asm5, C_REG_asm5, 13, 8, 0, 0, 0},
	{AMVN_5, C_NCON_asm5, C_NONE_asm5, C_REG_asm5, 13, 8, 0, 0, 0},
	{ACMP_5, C_NCON_asm5, C_REG_asm5, C_NONE_asm5, 13, 8, 0, 0, 0},
	{AADD_5, C_LCON_asm5, C_REG_asm5, C_REG_asm5, 13, 8, 0, LFROM_asm5, 0},
	{AADD_5, C_LCON_asm5, C_NONE_asm5, C_REG_asm5, 13, 8, 0, LFROM_asm5, 0},
	{AMVN_5, C_LCON_asm5, C_NONE_asm5, C_REG_asm5, 13, 8, 0, LFROM_asm5, 0},
	{ACMP_5, C_LCON_asm5, C_REG_asm5, C_NONE_asm5, 13, 8, 0, LFROM_asm5, 0},
	{AMOVB_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 1, 4, 0, 0, 0},
	{AMOVBS_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 14, 8, 0, 0, 0},
	{AMOVBU_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 58, 4, 0, 0, 0},
	{AMOVH_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 1, 4, 0, 0, 0},
	{AMOVHS_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 14, 8, 0, 0, 0},
	{AMOVHU_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 14, 8, 0, 0, 0},
	{AMUL_5, C_REG_asm5, C_REG_asm5, C_REG_asm5, 15, 4, 0, 0, 0},
	{AMUL_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 15, 4, 0, 0, 0},
	{ADIV_5, C_REG_asm5, C_REG_asm5, C_REG_asm5, 16, 4, 0, 0, 0},
	{ADIV_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 16, 4, 0, 0, 0},
	{AMULL_5, C_REG_asm5, C_REG_asm5, C_REGREG_asm5, 17, 4, 0, 0, 0},
	{AMULA_5, C_REG_asm5, C_REG_asm5, C_REGREG2_asm5, 17, 4, 0, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_SAUTO_asm5, 20, 4, REGSP_5, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_SOREG_asm5, 20, 4, 0, 0, 0},
	{AMOVB_5, C_REG_asm5, C_NONE_asm5, C_SAUTO_asm5, 20, 4, REGSP_5, 0, 0},
	{AMOVB_5, C_REG_asm5, C_NONE_asm5, C_SOREG_asm5, 20, 4, 0, 0, 0},
	{AMOVBS_5, C_REG_asm5, C_NONE_asm5, C_SAUTO_asm5, 20, 4, REGSP_5, 0, 0},
	{AMOVBS_5, C_REG_asm5, C_NONE_asm5, C_SOREG_asm5, 20, 4, 0, 0, 0},
	{AMOVBU_5, C_REG_asm5, C_NONE_asm5, C_SAUTO_asm5, 20, 4, REGSP_5, 0, 0},
	{AMOVBU_5, C_REG_asm5, C_NONE_asm5, C_SOREG_asm5, 20, 4, 0, 0, 0},
	{AMOVW_5, C_SAUTO_asm5, C_NONE_asm5, C_REG_asm5, 21, 4, REGSP_5, 0, 0},
	{AMOVW_5, C_SOREG_asm5, C_NONE_asm5, C_REG_asm5, 21, 4, 0, 0, 0},
	{AMOVBU_5, C_SAUTO_asm5, C_NONE_asm5, C_REG_asm5, 21, 4, REGSP_5, 0, 0},
	{AMOVBU_5, C_SOREG_asm5, C_NONE_asm5, C_REG_asm5, 21, 4, 0, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_LAUTO_asm5, 30, 8, REGSP_5, LTO_asm5, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_LOREG_asm5, 30, 8, 0, LTO_asm5, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_ADDR_asm5, 64, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVB_5, C_REG_asm5, C_NONE_asm5, C_LAUTO_asm5, 30, 8, REGSP_5, LTO_asm5, 0},
	{AMOVB_5, C_REG_asm5, C_NONE_asm5, C_LOREG_asm5, 30, 8, 0, LTO_asm5, 0},
	{AMOVB_5, C_REG_asm5, C_NONE_asm5, C_ADDR_asm5, 64, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVBS_5, C_REG_asm5, C_NONE_asm5, C_LAUTO_asm5, 30, 8, REGSP_5, LTO_asm5, 0},
	{AMOVBS_5, C_REG_asm5, C_NONE_asm5, C_LOREG_asm5, 30, 8, 0, LTO_asm5, 0},
	{AMOVBS_5, C_REG_asm5, C_NONE_asm5, C_ADDR_asm5, 64, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVBU_5, C_REG_asm5, C_NONE_asm5, C_LAUTO_asm5, 30, 8, REGSP_5, LTO_asm5, 0},
	{AMOVBU_5, C_REG_asm5, C_NONE_asm5, C_LOREG_asm5, 30, 8, 0, LTO_asm5, 0},
	{AMOVBU_5, C_REG_asm5, C_NONE_asm5, C_ADDR_asm5, 64, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVW_5, C_LAUTO_asm5, C_NONE_asm5, C_REG_asm5, 31, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVW_5, C_LOREG_asm5, C_NONE_asm5, C_REG_asm5, 31, 8, 0, LFROM_asm5, 0},
	{AMOVW_5, C_ADDR_asm5, C_NONE_asm5, C_REG_asm5, 65, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AMOVBU_5, C_LAUTO_asm5, C_NONE_asm5, C_REG_asm5, 31, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVBU_5, C_LOREG_asm5, C_NONE_asm5, C_REG_asm5, 31, 8, 0, LFROM_asm5, 0},
	{AMOVBU_5, C_ADDR_asm5, C_NONE_asm5, C_REG_asm5, 65, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AMOVW_5, C_LACON_asm5, C_NONE_asm5, C_REG_asm5, 34, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVW_5, C_PSR_asm5, C_NONE_asm5, C_REG_asm5, 35, 4, 0, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_PSR_asm5, 36, 4, 0, 0, 0},
	{AMOVW_5, C_RCON_asm5, C_NONE_asm5, C_PSR_asm5, 37, 4, 0, 0, 0},
	{AMOVM_5, C_LCON_asm5, C_NONE_asm5, C_SOREG_asm5, 38, 4, 0, 0, 0},
	{AMOVM_5, C_SOREG_asm5, C_NONE_asm5, C_LCON_asm5, 39, 4, 0, 0, 0},
	{ASWPW_5, C_SOREG_asm5, C_REG_asm5, C_REG_asm5, 40, 4, 0, 0, 0},
	{ARFE_5, C_NONE_asm5, C_NONE_asm5, C_NONE_asm5, 41, 4, 0, 0, 0},
	{AMOVF_5, C_FREG_asm5, C_NONE_asm5, C_FAUTO_asm5, 50, 4, REGSP_5, 0, 0},
	{AMOVF_5, C_FREG_asm5, C_NONE_asm5, C_FOREG_asm5, 50, 4, 0, 0, 0},
	{AMOVF_5, C_FAUTO_asm5, C_NONE_asm5, C_FREG_asm5, 51, 4, REGSP_5, 0, 0},
	{AMOVF_5, C_FOREG_asm5, C_NONE_asm5, C_FREG_asm5, 51, 4, 0, 0, 0},
	{AMOVF_5, C_FREG_asm5, C_NONE_asm5, C_LAUTO_asm5, 52, 12, REGSP_5, LTO_asm5, 0},
	{AMOVF_5, C_FREG_asm5, C_NONE_asm5, C_LOREG_asm5, 52, 12, 0, LTO_asm5, 0},
	{AMOVF_5, C_LAUTO_asm5, C_NONE_asm5, C_FREG_asm5, 53, 12, REGSP_5, LFROM_asm5, 0},
	{AMOVF_5, C_LOREG_asm5, C_NONE_asm5, C_FREG_asm5, 53, 12, 0, LFROM_asm5, 0},
	{AMOVF_5, C_FREG_asm5, C_NONE_asm5, C_ADDR_asm5, 68, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVF_5, C_ADDR_asm5, C_NONE_asm5, C_FREG_asm5, 69, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AADDF_5, C_FREG_asm5, C_NONE_asm5, C_FREG_asm5, 54, 4, 0, 0, 0},
	{AADDF_5, C_FREG_asm5, C_REG_asm5, C_FREG_asm5, 54, 4, 0, 0, 0},
	{AMOVF_5, C_FREG_asm5, C_NONE_asm5, C_FREG_asm5, 54, 4, 0, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_FCR_asm5, 56, 4, 0, 0, 0},
	{AMOVW_5, C_FCR_asm5, C_NONE_asm5, C_REG_asm5, 57, 4, 0, 0, 0},
	{AMOVW_5, C_SHIFT_asm5, C_NONE_asm5, C_REG_asm5, 59, 4, 0, 0, 0},
	{AMOVBU_5, C_SHIFT_asm5, C_NONE_asm5, C_REG_asm5, 59, 4, 0, 0, 0},
	{AMOVB_5, C_SHIFT_asm5, C_NONE_asm5, C_REG_asm5, 60, 4, 0, 0, 0},
	{AMOVBS_5, C_SHIFT_asm5, C_NONE_asm5, C_REG_asm5, 60, 4, 0, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_SHIFT_asm5, 61, 4, 0, 0, 0},
	{AMOVB_5, C_REG_asm5, C_NONE_asm5, C_SHIFT_asm5, 61, 4, 0, 0, 0},
	{AMOVBS_5, C_REG_asm5, C_NONE_asm5, C_SHIFT_asm5, 61, 4, 0, 0, 0},
	{AMOVBU_5, C_REG_asm5, C_NONE_asm5, C_SHIFT_asm5, 61, 4, 0, 0, 0},
	{ACASE_5, C_REG_asm5, C_NONE_asm5, C_NONE_asm5, 62, 4, 0, LPCREL_asm5, 8},
	{ABCASE_5, C_NONE_asm5, C_NONE_asm5, C_SBRA_asm5, 63, 4, 0, LPCREL_asm5, 0},
	{AMOVH_5, C_REG_asm5, C_NONE_asm5, C_HAUTO_asm5, 70, 4, REGSP_5, 0, 0},
	{AMOVH_5, C_REG_asm5, C_NONE_asm5, C_HOREG_asm5, 70, 4, 0, 0, 0},
	{AMOVHS_5, C_REG_asm5, C_NONE_asm5, C_HAUTO_asm5, 70, 4, REGSP_5, 0, 0},
	{AMOVHS_5, C_REG_asm5, C_NONE_asm5, C_HOREG_asm5, 70, 4, 0, 0, 0},
	{AMOVHU_5, C_REG_asm5, C_NONE_asm5, C_HAUTO_asm5, 70, 4, REGSP_5, 0, 0},
	{AMOVHU_5, C_REG_asm5, C_NONE_asm5, C_HOREG_asm5, 70, 4, 0, 0, 0},
	{AMOVB_5, C_HAUTO_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, REGSP_5, 0, 0},
	{AMOVB_5, C_HOREG_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, 0, 0, 0},
	{AMOVBS_5, C_HAUTO_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, REGSP_5, 0, 0},
	{AMOVBS_5, C_HOREG_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, 0, 0, 0},
	{AMOVH_5, C_HAUTO_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, REGSP_5, 0, 0},
	{AMOVH_5, C_HOREG_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, 0, 0, 0},
	{AMOVHS_5, C_HAUTO_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, REGSP_5, 0, 0},
	{AMOVHS_5, C_HOREG_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, 0, 0, 0},
	{AMOVHU_5, C_HAUTO_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, REGSP_5, 0, 0},
	{AMOVHU_5, C_HOREG_asm5, C_NONE_asm5, C_REG_asm5, 71, 4, 0, 0, 0},
	{AMOVH_5, C_REG_asm5, C_NONE_asm5, C_LAUTO_asm5, 72, 8, REGSP_5, LTO_asm5, 0},
	{AMOVH_5, C_REG_asm5, C_NONE_asm5, C_LOREG_asm5, 72, 8, 0, LTO_asm5, 0},
	{AMOVH_5, C_REG_asm5, C_NONE_asm5, C_ADDR_asm5, 94, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVHS_5, C_REG_asm5, C_NONE_asm5, C_LAUTO_asm5, 72, 8, REGSP_5, LTO_asm5, 0},
	{AMOVHS_5, C_REG_asm5, C_NONE_asm5, C_LOREG_asm5, 72, 8, 0, LTO_asm5, 0},
	{AMOVHS_5, C_REG_asm5, C_NONE_asm5, C_ADDR_asm5, 94, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVHU_5, C_REG_asm5, C_NONE_asm5, C_LAUTO_asm5, 72, 8, REGSP_5, LTO_asm5, 0},
	{AMOVHU_5, C_REG_asm5, C_NONE_asm5, C_LOREG_asm5, 72, 8, 0, LTO_asm5, 0},
	{AMOVHU_5, C_REG_asm5, C_NONE_asm5, C_ADDR_asm5, 94, 8, 0, LTO_asm5 | LPCREL_asm5, 4},
	{AMOVB_5, C_LAUTO_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVB_5, C_LOREG_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, 0, LFROM_asm5, 0},
	{AMOVB_5, C_ADDR_asm5, C_NONE_asm5, C_REG_asm5, 93, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AMOVBS_5, C_LAUTO_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVBS_5, C_LOREG_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, 0, LFROM_asm5, 0},
	{AMOVBS_5, C_ADDR_asm5, C_NONE_asm5, C_REG_asm5, 93, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AMOVH_5, C_LAUTO_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVH_5, C_LOREG_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, 0, LFROM_asm5, 0},
	{AMOVH_5, C_ADDR_asm5, C_NONE_asm5, C_REG_asm5, 93, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AMOVHS_5, C_LAUTO_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVHS_5, C_LOREG_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, 0, LFROM_asm5, 0},
	{AMOVHS_5, C_ADDR_asm5, C_NONE_asm5, C_REG_asm5, 93, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{AMOVHU_5, C_LAUTO_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, REGSP_5, LFROM_asm5, 0},
	{AMOVHU_5, C_LOREG_asm5, C_NONE_asm5, C_REG_asm5, 73, 8, 0, LFROM_asm5, 0},
	{AMOVHU_5, C_ADDR_asm5, C_NONE_asm5, C_REG_asm5, 93, 8, 0, LFROM_asm5 | LPCREL_asm5, 4},
	{ALDREX_5, C_SOREG_asm5, C_NONE_asm5, C_REG_asm5, 77, 4, 0, 0, 0},
	{ASTREX_5, C_SOREG_asm5, C_REG_asm5, C_REG_asm5, 78, 4, 0, 0, 0},
	{AMOVF_5, C_ZFCON_asm5, C_NONE_asm5, C_FREG_asm5, 80, 8, 0, 0, 0},
	{AMOVF_5, C_SFCON_asm5, C_NONE_asm5, C_FREG_asm5, 81, 4, 0, 0, 0},
	{ACMPF_5, C_FREG_asm5, C_REG_asm5, C_NONE_asm5, 82, 8, 0, 0, 0},
	{ACMPF_5, C_FREG_asm5, C_NONE_asm5, C_NONE_asm5, 83, 8, 0, 0, 0},
	{AMOVFW_5, C_FREG_asm5, C_NONE_asm5, C_FREG_asm5, 84, 4, 0, 0, 0},
	{AMOVWF_5, C_FREG_asm5, C_NONE_asm5, C_FREG_asm5, 85, 4, 0, 0, 0},
	{AMOVFW_5, C_FREG_asm5, C_NONE_asm5, C_REG_asm5, 86, 8, 0, 0, 0},
	{AMOVWF_5, C_REG_asm5, C_NONE_asm5, C_FREG_asm5, 87, 8, 0, 0, 0},
	{AMOVW_5, C_REG_asm5, C_NONE_asm5, C_FREG_asm5, 88, 4, 0, 0, 0},
	{AMOVW_5, C_FREG_asm5, C_NONE_asm5, C_REG_asm5, 89, 4, 0, 0, 0},
	{ATST_5, C_REG_asm5, C_NONE_asm5, C_NONE_asm5, 90, 4, 0, 0, 0},
	{ALDREXD_5, C_SOREG_asm5, C_NONE_asm5, C_REG_asm5, 91, 4, 0, 0, 0},
	{ASTREXD_5, C_SOREG_asm5, C_REG_asm5, C_REG_asm5, 92, 4, 0, 0, 0},
	{APLD_5, C_SOREG_asm5, C_NONE_asm5, C_NONE_asm5, 95, 4, 0, 0, 0},
	{AUNDEF_5, C_NONE_asm5, C_NONE_asm5, C_NONE_asm5, 96, 4, 0, 0, 0},
	{ACLZ_5, C_REG_asm5, C_NONE_asm5, C_REG_asm5, 97, 4, 0, 0, 0},
	{AMULWT_5, C_REG_asm5, C_REG_asm5, C_REG_asm5, 98, 4, 0, 0, 0},
	{AMULAWT_5, C_REG_asm5, C_REG_asm5, C_REGREG2_asm5, 99, 4, 0, 0, 0},
	{AUSEFIELD_5, C_ADDR_asm5, C_NONE_asm5, C_NONE_asm5, 0, 0, 0, 0, 0},
	{APCDATA_5, C_LCON_asm5, C_NONE_asm5, C_LCON_asm5, 0, 0, 0, 0, 0},
	{AFUNCDATA_5, C_LCON_asm5, C_NONE_asm5, C_ADDR_asm5, 0, 0, 0, 0, 0},
	{ADUFFZERO_5, C_NONE_asm5, C_NONE_asm5, C_SBRA_asm5, 5, 4, 0, 0, 0}, // same as ABL
	{ADUFFCOPY_5, C_NONE_asm5, C_NONE_asm5, C_SBRA_asm5, 5, 4, 0, 0, 0}, // same as ABL
	{ADATABUNDLE_5, C_NONE_asm5, C_NONE_asm5, C_NONE_asm5, 100, 4, 0, 0, 0},
	{ADATABUNDLEEND_5, C_NONE_asm5, C_NONE_asm5, C_NONE_asm5, 100, 0, 0, 0, 0},
	{AXXX_5, C_NONE_asm5, C_NONE_asm5, C_NONE_asm5, 0, 4, 0, 0, 0},
}

var pool_asm5 struct {
	start int
	size  int
	extra uint32
}

/*
 * when the first reference to the literal pool threatens
 * to go out of range of a 12-bit PC-relative offset,
 * drop the pool now, and branch round it.
 * this happens only in extended basic blocks that exceed 4k.
 */
func checkpool_asm5(ctxt *Link, p *Prog, sz int) int {
	if pool_asm5.size >= 0xff0 || immaddr_asm5((int32(p.pc)+int32(sz)+4)+4+(12+int32(pool_asm5.size))-(int32(pool_asm5.start)+8)) == 0 {
		return flushpool_asm5(ctxt, p, 1, 0)
	} else {
		if p.link == nil {
			return flushpool_asm5(ctxt, p, 2, 0)
		}
	}
	return 0
}

func flushpool_asm5(ctxt *Link, p *Prog, skip int, force int) int {
	var q *Prog
	if ctxt.blitrl != nil {
		if skip != 0 {
			if false && skip == 1 {
				print("note: flush literal pool at %llux: len=%ud ref=%ux\n", p.pc+4, pool_asm5.size, pool_asm5.start)
			}
			q = ctxt.arch.prg()
			q.as = int(AB_5)
			q.to.typ = int(D_BRANCH_5)
			q.pcond = p.link
			q.link = ctxt.blitrl
			q.lineno = p.lineno
			ctxt.blitrl = q
		} else {
			if !(force != 0) && (p.pc+(12+pool_asm5.size)-pool_asm5.start < 2048) { // 12 take into account the maximum nacl literal pool alignment padding size
				return 0
			}
		}
		if ctxt.headtype == int(Hnacl) && pool_asm5.size%16 != 0 {
			// if pool is not multiple of 16 bytes, add an alignment marker
			q = ctxt.arch.prg()
			q.as = int(ADATABUNDLEEND_5)
			ctxt.elitrl.link = q
			ctxt.elitrl = q
		}
		ctxt.elitrl.link = p.link
		p.link = ctxt.blitrl
		// BUG(minux): how to correctly handle line number for constant pool entries?
		// for now, we set line number to the last instruction preceding them at least
		// this won't bloat the .debug_line tables
		for ctxt.blitrl != nil {
			ctxt.blitrl.lineno = p.lineno
			ctxt.blitrl = ctxt.blitrl.link
		}
		ctxt.blitrl = nil /* BUG: should refer back to values until out-of-range */
		ctxt.elitrl = nil
		pool_asm5.size = 0
		pool_asm5.start = 0
		pool_asm5.extra = 0
		return 1
	}
	return 0
}

func addpool_asm5(ctxt *Link, p *Prog, a *Addr) {
	var q *Prog
	var t Prog
	var c int
	c = aclass_asm5(ctxt, a)
	t = zprg_asm5
	t.as = int(AWORD_5)
	switch c {
	default:
		t.to = *a
		if ctxt.flag_shared != 0 && t.to.sym != nil {
			t.pcrel = p
		}
		break
	case C_SROREG_asm5:
	case C_LOREG_asm5:
	case C_ROREG_asm5:
	case C_FOREG_asm5:
	case C_SOREG_asm5:
	case C_HOREG_asm5:
	case C_FAUTO_asm5:
	case C_SAUTO_asm5:
	case C_LAUTO_asm5:
	case C_LACON_asm5:
		t.to.typ = int(D_CONST_5)
		t.to.offset = int64(ctxt.instoffset)
		break
	}
	if t.pcrel == nil {
		for q = ctxt.blitrl; q != nil; q = q.link { /* could hash on t.t0.offset */
			if q.pcrel == nil && q.to == t.to {
				p.pcond = q
				return
			}
		}
	}
	if ctxt.headtype == int(Hnacl) && pool_asm5.size%16 == 0 {
		// start a new data bundle
		q = ctxt.arch.prg()
		*q = zprg_asm5
		q.as = int(ADATABUNDLE_5)
		q.pc = pool_asm5.size
		pool_asm5.size += 4
		if ctxt.blitrl == nil {
			ctxt.blitrl = q
			pool_asm5.start = p.pc
		} else {
			ctxt.elitrl.link = q
		}
		ctxt.elitrl = q
	}
	q = ctxt.arch.prg()
	*q = t
	q.pc = pool_asm5.size
	if ctxt.blitrl == nil {
		ctxt.blitrl = q
		pool_asm5.start = p.pc
	} else {
		ctxt.elitrl.link = q
	}
	ctxt.elitrl = q
	pool_asm5.size += 4
	p.pcond = q
}

func asmout_asm5(ctxt *Link, p *Prog, o *Optab_asm5, out []uint32) {
	var o1 uint32
	var o2 uint32
	var o3 uint32
	var o4 uint32
	var o5 uint32
	var o6 uint32
	var v int
	var r int
	var rf int
	var rt int
	var rt2 int
	var rel *Reloc
	ctxt.printp = p
	o1 = 0
	o2 = 0
	o3 = 0
	o4 = 0
	o5 = 0
	o6 = 0
	ctxt.armsize += int32(o.size)
	if false { /*debug['P']*/
		print("%ux: %P	type %d\n", uint32(p.pc), p, o.typ)
	}
	switch o.typ {
	default:
		ctxt.diag("unknown asm %d", o.typ)
		prasm_asm5(p)
		break
	case 0: /* pseudo ops */
		if false { /*debug['G']*/
			print("%ux: %s: arm %d\n", uint32(p.pc), p.from.sym.name, p.from.sym.fnptr)
		}
		break
	case 1: /* op R,[R],R */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		rf = p.from.reg
		rt = p.to.reg
		r = p.reg
		if p.to.typ == int(D_NONE_5) {
			rt = 0
		}
		if p.as == int(AMOVB_5) || p.as == int(AMOVH_5) || p.as == int(AMOVW_5) || p.as == int(AMVN_5) {
			r = 0
		} else {
			if r == int(NREG_5) {
				r = rt
			}
		}
		o1 |= uint32(rf) | (uint32(r) << 16) | (uint32(rt) << 12)
		break
	case 2: /* movbu $I,[R],R */
		aclass_asm5(ctxt, &p.from)
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= uint32(immrot_asm5(uint32(ctxt.instoffset)))
		rt = p.to.reg
		r = p.reg
		if p.to.typ == int(D_NONE_5) {
			rt = 0
		}
		if p.as == int(AMOVW_5) || p.as == int(AMVN_5) {
			r = 0
		} else {
			if r == int(NREG_5) {
				r = rt
			}
		}
		o1 |= (uint32(r) << 16) | (uint32(rt) << 12)
		break
	case 3: /* add R<<[IR],[R],R */
		o1 = mov_asm5(ctxt, p)
		break
	case 4: /* add $I,[R],R */
		aclass_asm5(ctxt, &p.from)
		o1 = oprrr_asm5(ctxt, int(AADD_5), p.scond)
		o1 |= uint32(immrot_asm5(uint32(ctxt.instoffset)))
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o1 |= uint32(r) << 16
		o1 |= uint32(p.to.reg) << 12
		break
	case 5: /* bra s */
		o1 = opbra_asm5(ctxt, p.as, p.scond)
		v = -8
		if p.to.sym != nil {
			rel = addrel(ctxt.cursym)
			rel.off = ctxt.pc
			rel.siz = 4
			rel.sym = p.to.sym
			v += int(p.to.offset)
			rel.add = int64(o1) | ((int64(v) >> 2) & 0xffffff)
			rel.typ = int(R_CALLARM)
			break
		}
		if p.pcond != nil {
			v = (p.pcond.pc - ctxt.pc) - 8
		}
		o1 |= (uint32(v) >> 2) & 0xffffff
		break
	case 6: /* b ,O(R) -> add $O,R,PC */
		aclass_asm5(ctxt, &p.to)
		o1 = oprrr_asm5(ctxt, int(AADD_5), p.scond)
		o1 |= uint32(immrot_asm5(uint32(ctxt.instoffset)))
		o1 |= uint32(p.to.reg) << 16
		o1 |= uint32(REGPC_5) << 12
		break
	case 7: /* bl (R) -> blx R */
		aclass_asm5(ctxt, &p.to)
		if ctxt.instoffset != 0 {
			ctxt.diag("%P: doesn't support BL offset(REG) where offset != 0", p)
		}
		o1 = oprrr_asm5(ctxt, int(ABL_5), p.scond)
		o1 |= uint32(p.to.reg)
		rel = addrel(ctxt.cursym)
		rel.off = ctxt.pc
		rel.siz = 0
		rel.typ = int(R_CALLIND)
		break
	case 8: /* sll $c,[R],R -> mov (R<<$c),R */
		aclass_asm5(ctxt, &p.from)
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		r = p.reg
		if r == int(NREG_5) {
			r = p.to.reg
		}
		o1 |= uint32(r)
		o1 |= (uint32(ctxt.instoffset) & 31) << 7
		o1 |= uint32(p.to.reg) << 12
		break
	case 9: /* sll R,[R],R -> mov (R<<R),R */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		r = p.reg
		if r == int(NREG_5) {
			r = p.to.reg
		}
		o1 |= uint32(r)
		o1 |= (uint32(p.from.reg) << 8) | (1 << 4)
		o1 |= uint32(p.to.reg) << 12
		break
	case 10: /* swi [$con] */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		if p.to.typ != int(D_NONE_5) {
			aclass_asm5(ctxt, &p.to)
			o1 |= uint32(ctxt.instoffset) & 0xffffff
		}
		break
	case 11: /* word */
		aclass_asm5(ctxt, &p.to)
		o1 = uint32(ctxt.instoffset)
		if p.to.sym != nil {
			// This case happens with words generated
			// in the PC stream as part of the literal pool.
			rel = addrel(ctxt.cursym)
			rel.off = ctxt.pc
			rel.siz = 4
			rel.sym = p.to.sym
			rel.add = p.to.offset
			// runtime.tlsg is special.
			// Its "address" is the offset from the TLS thread pointer
			// to the thread-local g and m pointers.
			// Emit a TLS relocation instead of a standard one.
			if rel.sym == ctxt.tlsg {
				rel.typ = int(R_TLS)
				if ctxt.flag_shared != 0 {
					rel.add += int64(ctxt.pc) - int64(p.pcrel.pc) - 8 - int64(rel.siz)
				}
				rel.xadd = rel.add
				rel.xsym = rel.sym
			} else {
				if ctxt.flag_shared != 0 {
					rel.typ = int(R_PCREL)
					rel.add += int64(ctxt.pc) - int64(p.pcrel.pc) - 8
				} else {
					rel.typ = int(R_ADDR)
				}
			}
			o1 = 0
		}
		break
	case 12: /* movw $lcon, reg */
		o1 = omvl_asm5(ctxt, p, &p.from, p.to.reg)
		if o.flag&int(LPCREL_asm5) != 0 {
			o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(p.to.reg) | uint32(REGPC_5)<<16 | uint32(p.to.reg)<<12
		}
		break
	case 13: /* op $lcon, [R], R */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = oprrr_asm5(ctxt, p.as, p.scond)
		o2 |= uint32(REGTMP_5)
		r = p.reg
		if p.as == int(AMOVW_5) || p.as == int(AMVN_5) {
			r = 0
		} else {
			if r == int(NREG_5) {
				r = p.to.reg
			}
		}
		o2 |= uint32(r) << 16
		if p.to.typ != int(D_NONE_5) {
			o2 |= uint32(p.to.reg) << 12
		}
		break
	case 14: /* movb/movbu/movh/movhu R,R */
		o1 = oprrr_asm5(ctxt, int(ASLL_5), p.scond)
		if p.as == int(AMOVBU_5) || p.as == int(AMOVHU_5) {
			o2 = oprrr_asm5(ctxt, int(ASRL_5), p.scond)
		} else {
			o2 = oprrr_asm5(ctxt, int(ASRA_5), p.scond)
		}
		r = p.to.reg
		o1 |= uint32(p.from.reg) | (uint32(r) << 12)
		o2 |= uint32(r) | (uint32(r) << 12)
		if p.as == int(AMOVB_5) || p.as == int(AMOVBS_5) || p.as == int(AMOVBU_5) {
			o1 |= (24 << 7)
			o2 |= (24 << 7)
		} else {
			o1 |= (16 << 7)
			o2 |= (16 << 7)
		}
		break
	case 15: /* mul r,[r,]r */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		rf = p.from.reg
		rt = p.to.reg
		r = p.reg
		if r == int(NREG_5) {
			r = rt
		}
		if rt == r {
			r = rf
			rf = rt
		}
		if false {
			if rt == r || rf == int(REGPC_5) || r == int(REGPC_5) || rt == int(REGPC_5) {
				ctxt.diag("bad registers in MUL")
				prasm_asm5(p)
			}
		}
		o1 |= (uint32(rf) << 8) | uint32(r) | (uint32(rt) << 16)
		break
	case 16: /* div r,[r,]r */
		o1 = 0xf << 28
		o2 = 0
		break
	case 17:
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		rf = p.from.reg
		rt = p.to.reg
		rt2 = int(p.to.offset)
		r = p.reg
		o1 |= (uint32(rf) << 8) | uint32(r) | (uint32(rt) << 16) | (uint32(rt2) << 12)
		break
	case 20: /* mov/movb/movbu R,O(R) */
		aclass_asm5(ctxt, &p.to)
		r = p.to.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o1 = osr_asm5(ctxt, p.as, p.from.reg, int(ctxt.instoffset), r, p.scond)
		break
	case 21: /* mov/movbu O(R),R -> lr */
		aclass_asm5(ctxt, &p.from)
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o1 = olr_asm5(ctxt, int(ctxt.instoffset), r, p.to.reg, p.scond)
		if p.as != int(AMOVW_5) {
			o1 |= 1 << 22
		}
		break
	case 30: /* mov/movb/movbu R,L(R) */
		o1 = omvl_asm5(ctxt, p, &p.to, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		r = p.to.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o2 = osrr_asm5(ctxt, p.from.reg, int(REGTMP_5), r, p.scond)
		if p.as != int(AMOVW_5) {
			o2 |= 1 << 22
		}
		break
	case 31: /* mov/movbu L(R),R -> lr[b] */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o2 = olrr_asm5(ctxt, int(REGTMP_5), r, p.to.reg, p.scond)
		if p.as == int(AMOVBU_5) || p.as == int(AMOVBS_5) || p.as == int(AMOVB_5) {
			o2 |= 1 << 22
		}
		break
	case 34: /* mov $lacon,R */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond)
		o2 |= uint32(REGTMP_5)
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o2 |= uint32(r) << 16
		if p.to.typ != int(D_NONE_5) {
			o2 |= uint32(p.to.reg) << 12
		}
		break
	case 35: /* mov PSR,R */
		o1 = (2 << 23) | (0xf << 16) | (0 << 0)
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		o1 |= (uint32(p.from.reg) & 1) << 22
		o1 |= uint32(p.to.reg) << 12
		break
	case 36: /* mov R,PSR */
		o1 = (2 << 23) | (0x29f << 12) | (0 << 4)
		if p.scond&int(C_FBIT_5) != 0 {
			o1 ^= 0x010 << 12
		}
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		o1 |= (uint32(p.to.reg) & 1) << 22
		o1 |= uint32(p.from.reg) << 0
		break
	case 37: /* mov $con,PSR */
		aclass_asm5(ctxt, &p.from)
		o1 = (2 << 23) | (0x29f << 12) | (0 << 4)
		if p.scond&int(C_FBIT_5) != 0 {
			o1 ^= 0x010 << 12
		}
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		o1 |= uint32(immrot_asm5(uint32(ctxt.instoffset)))
		o1 |= (uint32(p.to.reg) & 1) << 22
		o1 |= uint32(p.from.reg) << 0
		break
	case 38, 39:
		switch o.typ {
		case 38: /* movm $con,oreg -> stm */
			o1 = (0x4 << 25)
			o1 |= uint32(p.from.offset & 0xffff)
			o1 |= uint32(p.to.reg) << 16
			aclass_asm5(ctxt, &p.to)
		case 39: /* movm oreg,$con -> ldm */
			o1 = (0x4 << 25) | (1 << 20)
			o1 |= uint32(p.to.offset & 0xffff)
			o1 |= uint32(p.from.reg) << 16
			aclass_asm5(ctxt, &p.from)
		}
		if ctxt.instoffset != 0 {
			ctxt.diag("offset must be zero in MOVM; %P", p)
		}
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		if p.scond&int(C_PBIT_5) != 0 {
			o1 |= 1 << 24
		}
		if p.scond&int(C_UBIT_5) != 0 {
			o1 |= 1 << 23
		}
		if p.scond&int(C_SBIT_5) != 0 {
			o1 |= 1 << 22
		}
		if p.scond&int(C_WBIT_5) != 0 {
			o1 |= 1 << 21
		}
		break
	case 40: /* swp oreg,reg,reg */
		aclass_asm5(ctxt, &p.from)
		if ctxt.instoffset != 0 {
			ctxt.diag("offset must be zero in SWP")
		}
		o1 = (0x2 << 23) | (0x9 << 4)
		if p.as != int(ASWPW_5) {
			o1 |= 1 << 22
		}
		o1 |= uint32(p.from.reg) << 16
		o1 |= uint32(p.reg) << 0
		o1 |= uint32(p.to.reg) << 12
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		break
	case 41: /* rfe -> movm.s.w.u 0(r13),[r15] */
		o1 = 0xe8fd8000
		break
	case 50: /* floating point store */
		v = int(regoff_asm5(ctxt, &p.to))
		r = p.to.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o1 = ofsr_asm5(ctxt, p.as, p.from.reg, int32(v), r, p.scond, p)
		break
	case 51: /* floating point load */
		v = int(regoff_asm5(ctxt, &p.from))
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o1 = ofsr_asm5(ctxt, p.as, p.to.reg, int32(v), r, p.scond, p) | (1 << 20)
		break
	case 52: /* floating point store, int32 offset UGLY */
		o1 = omvl_asm5(ctxt, p, &p.to, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		r = p.to.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | (uint32(REGTMP_5) << 12) | (uint32(REGTMP_5) << 16) | uint32(r)
		o3 = ofsr_asm5(ctxt, p.as, p.from.reg, 0, int(REGTMP_5), p.scond, p)
		break
	case 53: /* floating point load, int32 offset UGLY */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | (uint32(REGTMP_5) << 12) | (uint32(REGTMP_5) << 16) | uint32(r)
		o3 = ofsr_asm5(ctxt, p.as, p.to.reg, 0, int(REGTMP_5), p.scond, p) | (1 << 20)
		break
	case 54: /* floating point arith */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		rf = p.from.reg
		rt = p.to.reg
		r = p.reg
		if r == int(NREG_5) {
			r = rt
			if p.as == int(AMOVF_5) || p.as == int(AMOVD_5) || p.as == int(AMOVFD_5) || p.as == int(AMOVDF_5) || p.as == int(ASQRTF_5) || p.as == int(ASQRTD_5) || p.as == int(AABSF_5) || p.as == int(AABSD_5) {
				r = 0
			}
		}
		o1 |= uint32(rf) | (uint32(r) << 16) | (uint32(rt) << 12)
		break
	case 56: /* move to FP[CS]R */
		o1 = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | (0xe << 24) | (1 << 8) | (1 << 4)
		o1 |= ((uint32(p.to.reg) + 1) << 21) | (uint32(p.from.reg) << 12)
		break
	case 57: /* move from FP[CS]R */
		o1 = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | (0xe << 24) | (1 << 8) | (1 << 4)
		o1 |= ((uint32(p.from.reg) + 1) << 21) | (uint32(p.to.reg) << 12) | (1 << 20)
		break
	case 58: /* movbu R,R */
		o1 = oprrr_asm5(ctxt, int(AAND_5), p.scond)
		o1 |= uint32(immrot_asm5(0xff))
		rt = p.to.reg
		r = p.from.reg
		if p.to.typ == int(D_NONE_5) {
			rt = 0
		}
		if r == int(NREG_5) {
			r = rt
		}
		o1 |= (uint32(r) << 16) | (uint32(rt) << 12)
		break
	case 59: /* movw/bu R<<I(R),R -> ldr indexed */
		if p.from.reg == int(NREG_5) {
			if p.as != int(AMOVW_5) {
				ctxt.diag("byte MOV from shifter operand")
			}
			o1 = mov_asm5(ctxt, p)
			break
		}
		if p.from.offset&(1<<4) != 0 /*untyped*/ {
			ctxt.diag("bad shift in LDR")
		}
		o1 = olrr_asm5(ctxt, int(p.from.offset), p.from.reg, p.to.reg, p.scond)
		if p.as == int(AMOVBU_5) {
			o1 |= 1 << 22
		}
		break
	case 60: /* movb R(R),R -> ldrsb indexed */
		if p.from.reg == int(NREG_5) {
			ctxt.diag("byte MOV from shifter operand")
			o1 = mov_asm5(ctxt, p)
			break
		}
		if p.from.offset&(^0xf) != 0 /*untyped*/ {
			ctxt.diag("bad shift in LDRSB")
		}
		o1 = olhrr_asm5(ctxt, int(p.from.offset), p.from.reg, p.to.reg, p.scond)
		o1 ^= (1 << 5) | (1 << 6)
		break
	case 61: /* movw/b/bu R,R<<[IR](R) -> str indexed */
		if p.to.reg == int(NREG_5) {
			ctxt.diag("MOV to shifter operand")
		}
		o1 = osrr_asm5(ctxt, p.from.reg, int(p.to.offset), p.to.reg, p.scond)
		if p.as == int(AMOVB_5) || p.as == int(AMOVBS_5) || p.as == int(AMOVBU_5) {
			o1 |= 1 << 22
		}
		break
	case 62: /* case R -> movw	R<<2(PC),PC */
		if o.flag&int(LPCREL_asm5) != 0 {
			o1 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(immrot_asm5(1)) | uint32(p.from.reg)<<16 | uint32(REGTMP_5)<<12
			o2 = olrr_asm5(ctxt, int(REGTMP_5), int(REGPC_5), int(REGTMP_5), p.scond)
			o2 |= 2 << 7
			o3 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(REGTMP_5) | uint32(REGPC_5)<<16 | uint32(REGPC_5)<<12
		} else {
			o1 = olrr_asm5(ctxt, p.from.reg, int(REGPC_5), int(REGPC_5), p.scond)
			o1 |= 2 << 7
		}
		break
	case 63: /* bcase */
		if p.pcond != nil {
			rel = addrel(ctxt.cursym)
			rel.off = ctxt.pc
			rel.siz = 4
			if p.to.sym != nil && p.to.sym.typ != 0 {
				rel.sym = p.to.sym
				rel.add = p.to.offset
			} else {
				rel.sym = ctxt.cursym
				rel.add = int64(p.pcond.pc)
			}
			if o.flag&int(LPCREL_asm5) != 0 {
				rel.typ = int(R_PCREL)
				rel.add += int64(ctxt.pc) - int64(p.pcrel.pc) - 16 + int64(rel.siz)
			} else {
				rel.typ = int(R_ADDR)
			}
			o1 = 0
		}
		break
	/* reloc ops */
	case 64: /* mov/movb/movbu R,addr */
		o1 = omvl_asm5(ctxt, p, &p.to, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = osr_asm5(ctxt, p.as, p.from.reg, 0, int(REGTMP_5), p.scond)
		if o.flag&int(LPCREL_asm5) != 0 {
			o3 = o2
			o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(REGTMP_5) | uint32(REGPC_5)<<16 | uint32(REGTMP_5)<<12
		}
		break
	case 65: /* mov/movbu addr,R */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = olr_asm5(ctxt, 0, int(REGTMP_5), p.to.reg, p.scond)
		if p.as == int(AMOVBU_5) || p.as == int(AMOVBS_5) || p.as == int(AMOVB_5) {
			o2 |= 1 << 22
		}
		if o.flag&int(LPCREL_asm5) != 0 {
			o3 = o2
			o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(REGTMP_5) | uint32(REGPC_5)<<16 | uint32(REGTMP_5)<<12
		}
		break
	case 68: /* floating point store -> ADDR */
		o1 = omvl_asm5(ctxt, p, &p.to, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = ofsr_asm5(ctxt, p.as, p.from.reg, 0, int(REGTMP_5), p.scond, p)
		if o.flag&int(LPCREL_asm5) != 0 {
			o3 = o2
			o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(REGTMP_5) | uint32(REGPC_5)<<16 | uint32(REGTMP_5)<<12
		}
		break
	case 69: /* floating point load <- ADDR */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = ofsr_asm5(ctxt, p.as, p.to.reg, 0, int(REGTMP_5), p.scond, p) | (1 << 20)
		if o.flag&int(LPCREL_asm5) != 0 {
			o3 = o2
			o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(REGTMP_5) | uint32(REGPC_5)<<16 | uint32(REGTMP_5)<<12
		}
		break
	/* ArmV4 ops: */
	case 70: /* movh/movhu R,O(R) -> strh */
		aclass_asm5(ctxt, &p.to)
		r = p.to.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o1 = oshr_asm5(ctxt, p.from.reg, int(ctxt.instoffset), r, p.scond)
		break
	case 71: /* movb/movh/movhu O(R),R -> ldrsb/ldrsh/ldrh */
		aclass_asm5(ctxt, &p.from)
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o1 = olhr_asm5(ctxt, int(ctxt.instoffset), r, p.to.reg, p.scond)
		if p.as == int(AMOVB_5) || p.as == int(AMOVBS_5) {
			o1 ^= (1 << 5) | (1 << 6)
		} else {
			if p.as == int(AMOVH_5) || p.as == int(AMOVHS_5) {
				o1 ^= (1 << 6)
			}
		}
		break
	case 72: /* movh/movhu R,L(R) -> strh */
		o1 = omvl_asm5(ctxt, p, &p.to, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		r = p.to.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o2 = oshrr_asm5(ctxt, p.from.reg, int(REGTMP_5), r, p.scond)
		break
	case 73: /* movb/movh/movhu L(R),R -> ldrsb/ldrsh/ldrh */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		r = p.from.reg
		if r == int(NREG_5) {
			r = o.param
		}
		o2 = olhrr_asm5(ctxt, int(REGTMP_5), r, p.to.reg, p.scond)
		if p.as == int(AMOVB_5) || p.as == int(AMOVBS_5) {
			o2 ^= (1 << 5) | (1 << 6)
		} else {
			if p.as == int(AMOVH_5) || p.as == int(AMOVHS_5) {
				o2 ^= (1 << 6)
			}
		}
		break
	case 74: /* bx $I */
		ctxt.diag("ABX $I")
		break
	case 75: /* bx O(R) */
		aclass_asm5(ctxt, &p.to)
		if ctxt.instoffset != 0 {
			ctxt.diag("non-zero offset in ABX")
		}
		/*
			o1 = 	oprrr(ctxt, AADD, p->scond) | immrot(0) | (REGPC<<16) | (REGLINK<<12);	// mov PC, LR
			o2 = ((p->scond&C_SCOND)<<28) | (0x12fff<<8) | (1<<4) | p->to.reg;		// BX R
		*/
		// p->to.reg may be REGLINK
		o1 = oprrr_asm5(ctxt, int(AADD_5), p.scond)
		o1 |= uint32(immrot_asm5(uint32(ctxt.instoffset)))
		o1 |= uint32(p.to.reg) << 16
		o1 |= uint32(REGTMP_5) << 12
		o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(immrot_asm5(0)) | (uint32(REGPC_5) << 16) | (uint32(REGLINK_5) << 12) // mov PC, LR
		o3 = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | (0x12fff << 8) | (1 << 4) | uint32(REGTMP_5)                          // BX Rtmp
		break
	case 76: /* bx O(R) when returning from fn*/
		ctxt.diag("ABXRET")
		break
	case 77: /* ldrex oreg,reg */
		aclass_asm5(ctxt, &p.from)
		if ctxt.instoffset != 0 {
			ctxt.diag("offset must be zero in LDREX")
		}
		o1 = (0x19 << 20) | uint32(0xf9f)
		o1 |= uint32(p.from.reg) << 16
		o1 |= uint32(p.to.reg) << 12
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		break
	case 78: /* strex reg,oreg,reg */
		aclass_asm5(ctxt, &p.from)
		if ctxt.instoffset != 0 {
			ctxt.diag("offset must be zero in STREX")
		}
		o1 = (0x18 << 20) | uint32(0xf90)
		o1 |= uint32(p.from.reg) << 16
		o1 |= uint32(p.reg) << 0
		o1 |= uint32(p.to.reg) << 12
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		break
	case 80: /* fmov zfcon,freg */
		if p.as == int(AMOVD_5) {
			o1 = 0xeeb00b00 // VMOV imm 64
			o2 = oprrr_asm5(ctxt, int(ASUBD_5), p.scond)
		} else {
			o1 = 0x0eb00a00 // VMOV imm 32
			o2 = oprrr_asm5(ctxt, int(ASUBF_5), p.scond)
		}
		v = 0x70 // 1.0
		r = p.to.reg
		// movf $1.0, r
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		o1 |= uint32(r) << 12
		o1 |= (uint32(v) & 0xf) << 0
		o1 |= (uint32(v) & 0xf0) << 12
		// subf r,r,r
		o2 |= uint32(r) | (uint32(r) << 16) | (uint32(r) << 12)
		break
	case 81: /* fmov sfcon,freg */
		o1 = 0x0eb00a00 // VMOV imm 32
		if p.as == int(AMOVD_5) {
			o1 = 0xeeb00b00 // VMOV imm 64
		}
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		o1 |= uint32(p.to.reg) << 12
		v = chipfloat5(ctxt, p.from.u.dval)
		o1 |= (uint32(v) & 0xf) << 0
		o1 |= (uint32(v) & 0xf0) << 12
		break
	case 82: /* fcmp freg,freg, */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= (uint32(p.reg) << 12) | (uint32(p.from.reg) << 0)
		o2 = 0x0ef1fa10 // VMRS R15
		o2 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		break
	case 83: /* fcmp freg,, */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= (uint32(p.from.reg) << 12) | (1 << 16)
		o2 = 0x0ef1fa10 // VMRS R15
		o2 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		break
	case 84: /* movfw freg,freg - truncate float-to-fix */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= (uint32(p.from.reg) << 0)
		o1 |= (uint32(p.to.reg) << 12)
		break
	case 85: /* movwf freg,freg - fix-to-float */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= (uint32(p.from.reg) << 0)
		o1 |= (uint32(p.to.reg) << 12)
		break
	// macro for movfw freg,FTMP; movw FTMP,reg
	case 86: /* movfw freg,reg - truncate float-to-fix */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= (uint32(p.from.reg) << 0)
		o1 |= (uint32(FREGTMP_5) << 12)
		o2 = oprrr_asm5(ctxt, int(AMOVFW_5+AEND_5), p.scond)
		o2 |= (uint32(FREGTMP_5) << 16)
		o2 |= (uint32(p.to.reg) << 12)
		break
	// macro for movw reg,FTMP; movwf FTMP,freg
	case 87: /* movwf reg,freg - fix-to-float */
		o1 = oprrr_asm5(ctxt, int(AMOVWF_5+AEND_5), p.scond)
		o1 |= (uint32(p.from.reg) << 12)
		o1 |= (uint32(FREGTMP_5) << 16)
		o2 = oprrr_asm5(ctxt, p.as, p.scond)
		o2 |= (uint32(FREGTMP_5) << 0)
		o2 |= (uint32(p.to.reg) << 12)
		break
	case 88: /* movw reg,freg  */
		o1 = oprrr_asm5(ctxt, int(AMOVWF_5+AEND_5), p.scond)
		o1 |= (uint32(p.from.reg) << 12)
		o1 |= (uint32(p.to.reg) << 16)
		break
	case 89: /* movw freg,reg  */
		o1 = oprrr_asm5(ctxt, int(AMOVFW_5+AEND_5), p.scond)
		o1 |= (uint32(p.from.reg) << 16)
		o1 |= (uint32(p.to.reg) << 12)
		break
	case 90: /* tst reg  */
		o1 = oprrr_asm5(ctxt, int(ACMP_5+AEND_5), p.scond)
		o1 |= uint32(p.from.reg) << 16
		break
	case 91: /* ldrexd oreg,reg */
		aclass_asm5(ctxt, &p.from)
		if ctxt.instoffset != 0 {
			ctxt.diag("offset must be zero in LDREX")
		}
		o1 = (0x1b << 20) | uint32(0xf9f)
		o1 |= uint32(p.from.reg) << 16
		o1 |= uint32(p.to.reg) << 12
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		break
	case 92: /* strexd reg,oreg,reg */
		aclass_asm5(ctxt, &p.from)
		if ctxt.instoffset != 0 {
			ctxt.diag("offset must be zero in STREX")
		}
		o1 = (0x1a << 20) | uint32(0xf90)
		o1 |= uint32(p.from.reg) << 16
		o1 |= uint32(p.reg) << 0
		o1 |= uint32(p.to.reg) << 12
		o1 |= (uint32(p.scond) & uint32(C_SCOND_5)) << 28
		break
	case 93: /* movb/movh/movhu addr,R -> ldrsb/ldrsh/ldrh */
		o1 = omvl_asm5(ctxt, p, &p.from, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = olhr_asm5(ctxt, 0, int(REGTMP_5), p.to.reg, p.scond)
		if p.as == int(AMOVB_5) || p.as == int(AMOVBS_5) {
			o2 ^= (1 << 5) | (1 << 6)
		} else {
			if p.as == int(AMOVH_5) || p.as == int(AMOVHS_5) {
				o2 ^= (1 << 6)
			}
		}
		if o.flag&int(LPCREL_asm5) != 0 {
			o3 = o2
			o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(REGTMP_5) | uint32(REGPC_5)<<16 | uint32(REGTMP_5)<<12
		}
		break
	case 94: /* movh/movhu R,addr -> strh */
		o1 = omvl_asm5(ctxt, p, &p.to, int(REGTMP_5))
		if !(o1 != 0) {
			break
		}
		o2 = oshr_asm5(ctxt, p.from.reg, 0, int(REGTMP_5), p.scond)
		if o.flag&int(LPCREL_asm5) != 0 {
			o3 = o2
			o2 = oprrr_asm5(ctxt, int(AADD_5), p.scond) | uint32(REGTMP_5) | uint32(REGPC_5)<<16 | uint32(REGTMP_5)<<12
		}
		break
	case 95: /* PLD off(reg) */
		o1 = 0xf5d0f000
		o1 |= uint32(p.from.reg) << 16
		if p.from.offset < 0 {
			o1 &^= (1 << 23)
			o1 |= uint32((-p.from.offset) & 0xfff)
		} else {
			o1 |= uint32(p.from.offset & 0xfff)
		}
		break
	// This is supposed to be something that stops execution.
	// It's not supposed to be reached, ever, but if it is, we'd
	// like to be able to tell how we got there.  Assemble as
	// 0xf7fabcfd which is guaranteed to raise undefined instruction
	// exception.
	case 96: /* UNDEF */
		o1 = 0xf7fabcfd
		break
	case 97: /* CLZ Rm, Rd */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= uint32(p.to.reg) << 12
		o1 |= uint32(p.from.reg)
		break
	case 98: /* MULW{T,B} Rs, Rm, Rd */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= uint32(p.to.reg) << 16
		o1 |= uint32(p.from.reg) << 8
		o1 |= uint32(p.reg)
		break
	case 99: /* MULAW{T,B} Rs, Rm, Rn, Rd */
		o1 = oprrr_asm5(ctxt, p.as, p.scond)
		o1 |= uint32(p.to.reg) << 12
		o1 |= uint32(p.from.reg) << 8
		o1 |= uint32(p.reg)
		o1 |= uint32(p.to.offset << 16)
		break
	// DATABUNDLE: BKPT $0x5be0, signify the start of NaCl data bundle;
	// DATABUNDLEEND: zero width alignment marker
	case 100:
		if p.as == int(ADATABUNDLE_5) {
			o1 = 0xe125be70
		}
		break
	}
	out[0] = o1
	out[1] = o2
	out[2] = o3
	out[3] = o4
	out[4] = o5
	out[5] = o6
	return
}

func mov_asm5(ctxt *Link, p *Prog) uint32 {
	aclass_asm5(ctxt, &p.from)
	o1 := oprrr_asm5(ctxt, p.as, p.scond)
	o1 |= uint32(p.from.offset)
	rt := p.to.reg
	r := p.reg
	if p.to.typ == int(D_NONE_5) {
		rt = 0
	}
	if p.as == int(AMOVW_5) || p.as == int(AMVN_5) {
		r = 0
	} else {
		if r == int(NREG_5) {
			r = rt
		}
	}
	o1 |= (uint32(r) << 16) | (uint32(rt) << 12)
	return o1
}

// asmoutnacl assembles the instruction p. It replaces asmout for NaCl.
// It returns the total number of bytes put in out, and it can change
// p->pc if extra padding is necessary.
// In rare cases, asmoutnacl might split p into two instructions.
// origPC is the PC for this Prog (no padding is taken into account).
func asmoutnacl_asm5(ctxt *Link, origPC int32, p *Prog, o *Optab_asm5, out []uint32) int {
	var size int
	var reg int
	var q *Prog
	var a *Addr
	var a2 *Addr
	size = o.size
	// instruction specific
	switch p.as {
	default:
		if out != nil {
			asmout_asm5(ctxt, p, o, out)
		}
		break
	case ADATABUNDLE_5: // align to 16-byte boundary
	case ADATABUNDLEEND_5: // zero width instruction, just to align next instruction to 16-byte boundary
		p.pc = (p.pc + 15) &^ 15
		if out != nil {
			asmout_asm5(ctxt, p, o, out)
		}
		break
	case AUNDEF_5:
	case APLD_5:
		size = 4
		if out != nil {
			switch p.as {
			case AUNDEF_5:
				out[0] = 0xe7fedef0 // NACL_INSTR_ARM_ABORT_NOW (UDF #0xEDE0)
				break
			case APLD_5:
				out[0] = 0xe1a01001 // (MOVW R1, R1)
				break
			}
		}
		break
	case AB_5:
	case ABL_5:
		if p.to.typ != int(D_OREG_5) {
			if out != nil {
				asmout_asm5(ctxt, p, o, out)
			}
		} else {
			if p.to.offset != 0 || size != 4 || p.to.reg >= 16 || p.to.reg < 0 {
				ctxt.diag("unsupported instruction: %P", p)
			}
			if (p.pc & 15) == 12 {
				p.pc += 4
			}
			if out != nil {
				out[0] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x03c0013f | (uint32(p.to.reg) << 12) | (uint32(p.to.reg) << 16) // BIC $0xc000000f, Rx
				if p.as == int(AB_5) {
					out[1] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x012fff10 | uint32(p.to.reg) // BX Rx // ABL
				} else {
					out[1] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x012fff30 | uint32(p.to.reg) // BLX Rx
				}
			}
			size = 8
		}
		// align the last instruction (the actual BL) to the last instruction in a bundle
		if p.as == int(ABL_5) {
			if deferreturn_asm5 == nil {
				deferreturn_asm5 = linklookup(ctxt, "runtime.deferreturn", 0)
			}
			if p.to.sym == deferreturn_asm5 {
				p.pc = int(((origPC + 15) &^ 15) + 16 - int32(size))
			} else {
				p.pc += (16 - ((p.pc + size) & 15)) & 15
			}
		}
		break
	case ALDREX_5:
	case ALDREXD_5:
	case AMOVB_5:
	case AMOVBS_5:
	case AMOVBU_5:
	case AMOVD_5:
	case AMOVF_5:
	case AMOVH_5:
	case AMOVHS_5:
	case AMOVHU_5:
	case AMOVM_5:
	case AMOVW_5:
	case ASTREX_5:
	case ASTREXD_5:
		if p.to.typ == int(D_REG_5) && p.to.reg == 15 && p.from.reg == 13 { // MOVW.W x(R13), PC
			if out != nil {
				asmout_asm5(ctxt, p, o, out)
			}
			if size == 4 {
				if out != nil {
					// Note: 5c and 5g reg.c know that DIV/MOD smashes R12
					// so that this return instruction expansion is valid.
					out[0] = out[0] &^ 0x3000                                           // change PC to R12
					out[1] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x03ccc13f // BIC $0xc000000f, R12
					out[2] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x012fff1c // BX R12
				}
				size += 8
				if ((p.pc + size) & 15) == 4 {
					p.pc += 4
				}
				break
			} else {
				// if the instruction used more than 4 bytes, then it must have used a very large
				// offset to update R13, so we need to additionally mask R13.
				if out != nil {
					out[size/4-1] &^= 0x3000                                                   // change PC to R12
					out[size/4] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x03cdd103   // BIC $0xc0000000, R13
					out[size/4+1] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x03ccc13f // BIC $0xc000000f, R12
					out[size/4+2] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x012fff1c // BX R12
				}
				// p->pc+size is only ok at 4 or 12 mod 16.
				if (p.pc+size)%8 == 0 {
					p.pc += 4
				}
				size += 12
				break
			}
		}
		if p.to.typ == int(D_REG_5) && p.to.reg == 15 {
			ctxt.diag("unsupported instruction (move to another register and use indirect jump instead): %P", p)
		}
		if p.to.typ == int(D_OREG_5) && p.to.reg == 13 && (p.scond&int(C_WBIT_5) != 0) && size > 4 {
			// function prolog with very large frame size: MOVW.W R14,-100004(R13)
			// split it into two instructions:
			// 	ADD $-100004, R13
			// 	MOVW R14, 0(R13)
			q = ctxt.arch.prg()
			p.scond &^= int(C_WBIT_5)
			*q = *p
			a = &p.to
			if p.to.typ == int(D_OREG_5) {
				a2 = &q.to
			} else {
				a2 = &q.from
			}
			nocache_asm5(q)
			nocache_asm5(p)
			// insert q after p
			q.link = p.link
			p.link = q
			q.pcond = (*Prog)(nil)
			// make p into ADD $X, R13
			p.as = int(AADD_5)
			p.from = *a
			p.from.reg = int(NREG_5)
			p.from.typ = int(D_CONST_5)
			p.to = zprg_asm5.to
			p.to.typ = int(D_REG_5)
			p.to.reg = 13
			// make q into p but load/store from 0(R13)
			q.spadj = 0
			*a2 = zprg_asm5.from
			a2.typ = int(D_OREG_5)
			a2.reg = 13
			a2.sym = (*LSym)(nil)
			a2.offset = 0
			size = oplook_asm5(ctxt, p).size
			break
		}
		if (p.to.typ == int(D_OREG_5) && p.to.reg != 13 && p.to.reg != 9) || (p.from.typ == int(D_OREG_5) && p.from.reg != 13 && p.from.reg != 9) { // MOVW Rx, X(Ry), y != 13 && y != 9 // MOVW X(Rx), Ry, x != 13 && x != 9
			if p.to.typ == int(D_OREG_5) {
				a = &p.to
			} else {
				a = &p.from
			}
			reg = a.reg
			if size == 4 {
				// if addr.reg == NREG, then it is probably load from x(FP) with small x, no need to modify.
				if reg == int(NREG_5) {
					if out != nil {
						asmout_asm5(ctxt, p, o, out)
					}
				} else {
					if out != nil {
						out[0] = ((uint32(p.scond) & uint32(C_SCOND_5)) << 28) | 0x03c00103 | (uint32(reg) << 16) | (uint32(reg) << 12) // BIC $0xc0000000, Rx
					}
					if (p.pc & 15) == 12 {
						p.pc += 4
					}
					size += 4
					if out != nil {
						asmout_asm5(ctxt, p, o, out[1:])
					}
				}
				break
			} else {
				// if a load/store instruction takes more than 1 word to implement, then
				// we need to seperate the instruction into two:
				// 1. explicitly load the address into R11.
				// 2. load/store from R11.
				// This won't handle .W/.P, so we should reject such code.
				if p.scond&int(C_PBIT_5|C_WBIT_5) != 0 {
					ctxt.diag("unsupported instruction (.P/.W): %P", p)
				}
				q = ctxt.arch.prg()
				*q = *p
				if p.to.typ == int(D_OREG_5) {
					a2 = &q.to
				} else {
					a2 = &q.from
				}
				nocache_asm5(q)
				nocache_asm5(p)
				// insert q after p
				q.link = p.link
				p.link = q
				q.pcond = (*Prog)(nil)
				// make p into MOVW $X(R), R11
				p.as = int(AMOVW_5)
				p.from = *a
				p.from.typ = int(D_CONST_5)
				p.to = zprg_asm5.to
				p.to.typ = int(D_REG_5)
				p.to.reg = 11
				// make q into p but load/store from 0(R11)
				*a2 = zprg_asm5.from
				a2.typ = int(D_OREG_5)
				a2.reg = 11
				a2.sym = (*LSym)(nil)
				a2.offset = 0
				size = oplook_asm5(ctxt, p).size
				break
			}
		} else {
			if out != nil {
				asmout_asm5(ctxt, p, o, out)
			}
		}
		break
	}
	// destination register specific
	if p.to.typ == int(D_REG_5) {
		switch p.to.reg {
		case 9:
			ctxt.diag("invalid instruction, cannot write to R9: %P", p)
			break
		case 13:
			if out != nil {
				out[size/4] = 0xe3cdd103 // BIC $0xc0000000, R13
			}
			if ((p.pc + size) & 15) == 0 {
				p.pc += 4
			}
			size += 4
			break
		}
	}
	return size
}

func oplook_asm5(ctxt *Link, p *Prog) *Optab_asm5 {
	var a1 int
	var a2 int
	var a3 int
	var r int
	var c1 []int8
	var c3 []int8
	var o []Optab_asm5
	var e []Optab_asm5
	a1 = p.optab
	if a1 != 0 {
		return &optab_asm5[a1-1:][0]
	}
	a1 = p.from.class
	if a1 == 0 {
		a1 = aclass_asm5(ctxt, &p.from) + 1
		p.from.class = a1
	}
	a1--
	a3 = p.to.class
	if a3 == 0 {
		a3 = aclass_asm5(ctxt, &p.to) + 1
		p.to.class = a3
	}
	a3--
	a2 = int(C_NONE_asm5)
	if p.reg != int(NREG_5) {
		a2 = int(C_REG_asm5)
	}
	r = p.as
	o = oprange_asm5[r].start
	if o == nil {
		a1 = int(opcross_asm5[repop_asm5[r]][a1][a2][a3])
		if a1 != 0 {
			p.optab = a1 + 1
			return &optab_asm5[a1:][0]
		}
		o = oprange_asm5[r].stop /* just generate an error */
	}
	if false { /*debug['O']*/
		print("oplook %A %d %d %d\n", int(p.as), a1, a2, a3)
		print("		%d %d\n", p.from.typ, p.to.typ)
	}
	e = oprange_asm5[r].stop
	c1 = xcmp_asm5[a1][:]
	c3 = xcmp_asm5[a3][:]
	for ; -cap(o) < -cap(e); o = o[1:] {
		if o[0].a2 == a2 {
			if c1[o[0].a1] != 0 {
				if c3[o[0].a3] != 0 {
					p.optab = (-cap(o) + cap(optab_asm5)) + 1
					return &o[0]
				}
			}
		}
	}
	ctxt.diag("illegal combination %P; %d %d %d, %d %d", p, a1, a2, a3, p.from.typ, p.to.typ)
	ctxt.diag("from %d %d to %d %d\n", p.from.typ, p.from.name, p.to.typ, p.to.name)
	prasm_asm5(p)
	if o == nil {
		o = optab_asm5[0:]
	}
	return &o[0]
}

func oprrr_asm5(ctxt *Link, a int, sc int) uint32 {
	var o int32
	o = (int32(sc) & int32(C_SCOND_5)) << 28
	if sc&int(C_SBIT_5) != 0 {
		o |= 1 << 20
	}
	if sc&int(C_PBIT_5|C_WBIT_5) != 0 {
		ctxt.diag(".nil/.W on dp instruction")
	}
	switch a {
	case AMULU_5:
	case AMUL_5:
		return uint32(o) | (0x0 << 21) | (0x9 << 4)
	case AMULA_5:
		return uint32(o) | (0x1 << 21) | (0x9 << 4)
	case AMULLU_5:
		return uint32(o) | (0x4 << 21) | (0x9 << 4)
	case AMULL_5:
		return uint32(o) | (0x6 << 21) | (0x9 << 4)
	case AMULALU_5:
		return uint32(o) | (0x5 << 21) | (0x9 << 4)
	case AMULAL_5:
		return uint32(o) | (0x7 << 21) | (0x9 << 4)
	case AAND_5:
		return uint32(o) | (0x0 << 21)
	case AEOR_5:
		return uint32(o) | (0x1 << 21)
	case ASUB_5:
		return uint32(o) | (0x2 << 21)
	case ARSB_5:
		return uint32(o) | (0x3 << 21)
	case AADD_5:
		return uint32(o) | (0x4 << 21)
	case AADC_5:
		return uint32(o) | (0x5 << 21)
	case ASBC_5:
		return uint32(o) | (0x6 << 21)
	case ARSC_5:
		return uint32(o) | (0x7 << 21)
	case ATST_5:
		return uint32(o) | (0x8 << 21) | (1 << 20)
	case ATEQ_5:
		return uint32(o) | (0x9 << 21) | (1 << 20)
	case ACMP_5:
		return uint32(o) | (0xa << 21) | (1 << 20)
	case ACMN_5:
		return uint32(o) | (0xb << 21) | (1 << 20)
	case AORR_5:
		return uint32(o) | (0xc << 21)
	case AMOVB_5:
	case AMOVH_5:
	case AMOVW_5:
		return uint32(o) | (0xd << 21)
	case ABIC_5:
		return uint32(o) | (0xe << 21)
	case AMVN_5:
		return uint32(o) | (0xf << 21)
	case ASLL_5:
		return uint32(o) | (0xd << 21) | (0 << 5)
	case ASRL_5:
		return uint32(o) | (0xd << 21) | (1 << 5)
	case ASRA_5:
		return uint32(o) | (0xd << 21) | (2 << 5)
	case ASWI_5:
		return uint32(o) | (0xf << 24)
	case AADDD_5:
		return uint32(o) | (0xe << 24) | (0x3 << 20) | (0xb << 8) | (0 << 4)
	case AADDF_5:
		return uint32(o) | (0xe << 24) | (0x3 << 20) | (0xa << 8) | (0 << 4)
	case ASUBD_5:
		return uint32(o) | (0xe << 24) | (0x3 << 20) | (0xb << 8) | (4 << 4)
	case ASUBF_5:
		return uint32(o) | (0xe << 24) | (0x3 << 20) | (0xa << 8) | (4 << 4)
	case AMULD_5:
		return uint32(o) | (0xe << 24) | (0x2 << 20) | (0xb << 8) | (0 << 4)
	case AMULF_5:
		return uint32(o) | (0xe << 24) | (0x2 << 20) | (0xa << 8) | (0 << 4)
	case ADIVD_5:
		return uint32(o) | (0xe << 24) | (0x8 << 20) | (0xb << 8) | (0 << 4)
	case ADIVF_5:
		return uint32(o) | (0xe << 24) | (0x8 << 20) | (0xa << 8) | (0 << 4)
	case ASQRTD_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (1 << 16) | (0xb << 8) | (0xc << 4)
	case ASQRTF_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (1 << 16) | (0xa << 8) | (0xc << 4)
	case AABSD_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (0 << 16) | (0xb << 8) | (0xc << 4)
	case AABSF_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (0 << 16) | (0xa << 8) | (0xc << 4)
	case ACMPD_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (4 << 16) | (0xb << 8) | (0xc << 4)
	case ACMPF_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (4 << 16) | (0xa << 8) | (0xc << 4)
	case AMOVF_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (0 << 16) | (0xa << 8) | (4 << 4)
	case AMOVD_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (0 << 16) | (0xb << 8) | (4 << 4)
	case AMOVDF_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (7 << 16) | (0xa << 8) | (0xc << 4) | (1 << 8) // dtof
	case AMOVFD_5:
		return uint32(o) | (0xe << 24) | (0xb << 20) | (7 << 16) | (0xa << 8) | (0xc << 4) | (0 << 8) // dtof
	case AMOVWF_5:
		if (sc & int(C_UBIT_5)) == 0 {
			o |= 1 << 7 /* signed */
		}
		return uint32(o) | (0xe << 24) | (0xb << 20) | (8 << 16) | (0xa << 8) | (4 << 4) | (0 << 18) | (0 << 8) // toint, double
	case AMOVWD_5:
		if (sc & int(C_UBIT_5)) == 0 {
			o |= 1 << 7 /* signed */
		}
		return uint32(o) | (0xe << 24) | (0xb << 20) | (8 << 16) | (0xa << 8) | (4 << 4) | (0 << 18) | (1 << 8) // toint, double
	case AMOVFW_5:
		if (sc & int(C_UBIT_5)) == 0 {
			o |= 1 << 16 /* signed */
		}
		return uint32(o) | (0xe << 24) | (0xb << 20) | (8 << 16) | (0xa << 8) | (4 << 4) | (1 << 18) | (0 << 8) | (1 << 7) // toint, double, trunc
	case AMOVDW_5:
		if (sc & int(C_UBIT_5)) == 0 {
			o |= 1 << 16 /* signed */
		}
		return uint32(o) | (0xe << 24) | (0xb << 20) | (8 << 16) | (0xa << 8) | (4 << 4) | (1 << 18) | (1 << 8) | (1 << 7) // toint, double, trunc
	case AMOVWF_5 + AEND_5: // copy WtoF
		return uint32(o) | (0xe << 24) | (0x0 << 20) | (0xb << 8) | (1 << 4)
	case AMOVFW_5 + AEND_5: // copy FtoW
		return uint32(o) | (0xe << 24) | (0x1 << 20) | (0xb << 8) | (1 << 4)
	case ACMP_5 + AEND_5: // cmp imm
		return uint32(o) | (0x3 << 24) | (0x5 << 20)
	// CLZ doesn't support .nil
	case ACLZ_5:
		return (uint32(o) & (0xf << 28)) | (0x16f << 16) | (0xf1 << 4)
	case AMULWT_5:
		return (uint32(o) & (0xf << 28)) | (0x12 << 20) | (0xe << 4)
	case AMULWB_5:
		return (uint32(o) & (0xf << 28)) | (0x12 << 20) | (0xa << 4)
	case AMULAWT_5:
		return (uint32(o) & (0xf << 28)) | (0x12 << 20) | (0xc << 4)
	case AMULAWB_5:
		return (uint32(o) & (0xf << 28)) | (0x12 << 20) | (0x8 << 4)
	case ABL_5: // BLX REG
		return (uint32(o) & (0xf << 28)) | (0x12fff3 << 4)
	}
	ctxt.diag("bad rrr %d", a)
	prasm_asm5(ctxt.curp)
	return 0
}

func olr_asm5(ctxt *Link, v int, b int, r int, sc int) uint32 {
	var o uint32
	if sc&int(C_SBIT_5) != 0 {
		ctxt.diag(".nil on LDR/STR instruction")
	}
	o = (uint32(sc) & uint32(C_SCOND_5)) << 28
	if !(sc&int(C_PBIT_5) != 0) {
		o |= 1 << 24
	}
	if !(sc&int(C_UBIT_5) != 0) {
		o |= 1 << 23
	}
	if sc&int(C_WBIT_5) != 0 {
		o |= 1 << 21
	}
	o |= (1 << 26) | (1 << 20)
	if v < 0 {
		if sc&int(C_UBIT_5) != 0 {
			ctxt.diag(".U on neg offset")
		}
		v = -v
		o ^= 1 << 23
	}
	if v >= (1<<12) || v < 0 {
		ctxt.diag("literal span too large: %d (R%d)\n%P", v, b, ctxt.printp)
	}
	o |= uint32(v)
	o |= uint32(b) << 16
	o |= uint32(r) << 12
	return o
}

func olhr_asm5(ctxt *Link, v int, b int, r int, sc int) uint32 {
	var o uint32
	if sc&int(C_SBIT_5) != 0 {
		ctxt.diag(".nil on LDRH/STRH instruction")
	}
	o = (uint32(sc) & uint32(C_SCOND_5)) << 28
	if !(sc&int(C_PBIT_5) != 0) {
		o |= 1 << 24
	}
	if sc&int(C_WBIT_5) != 0 {
		o |= 1 << 21
	}
	o |= (1 << 23) | (1 << 20) | (0xb << 4)
	if v < 0 {
		v = -v
		o ^= 1 << 23
	}
	if v >= (1<<8) || v < 0 {
		ctxt.diag("literal span too large: %d (R%d)\n%P", v, b, ctxt.printp)
	}
	o |= (uint32(v) & 0xf) | ((uint32(v) >> 4) << 8) | (1 << 22)
	o |= uint32(b) << 16
	o |= uint32(r) << 12
	return o
}

func olrr_asm5(ctxt *Link, i int, b int, r int, sc int) uint32 {
	return olr_asm5(ctxt, i, b, r, sc) ^ (1 << 25)
}

func olhrr_asm5(ctxt *Link, i int, b int, r int, sc int) uint32 {
	return olhr_asm5(ctxt, i, b, r, sc) ^ (1 << 22)
}

func osr_asm5(ctxt *Link, a int, r int, v int, b int, sc int) uint32 {
	var o uint32
	o = olr_asm5(ctxt, v, b, r, sc) ^ (1 << 20)
	if a != int(AMOVW_5) {
		o |= 1 << 22
	}
	return o
}

func oshr_asm5(ctxt *Link, r int, v int, b int, sc int) uint32 {
	var o uint32
	o = olhr_asm5(ctxt, v, b, r, sc) ^ (1 << 20)
	return o
}

func ofsr_asm5(ctxt *Link, a int, r int, v int32, b int, sc int, p *Prog) uint32 {
	var o uint32
	if sc&int(C_SBIT_5) != 0 {
		ctxt.diag(".nil on FLDR/FSTR instruction")
	}
	o = (uint32(sc) & uint32(C_SCOND_5)) << 28
	if !(sc&int(C_PBIT_5) != 0) {
		o |= 1 << 24
	}
	if sc&int(C_WBIT_5) != 0 {
		o |= 1 << 21
	}
	o |= (6 << 25) | (1 << 24) | (1 << 23) | (10 << 8)
	if v < 0 {
		v = -v
		o ^= 1 << 23
	}
	if v&3 != 0 /*untyped*/ {
		ctxt.diag("odd offset for floating point op: %d\n%P", v, p)
	} else {
		if v >= (1<<10) || v < 0 {
			ctxt.diag("literal span too large: %d\n%P", v, p)
		}
	}
	o |= (uint32(v) >> 2) & 0xFF
	o |= uint32(b) << 16
	o |= uint32(r) << 12
	switch a {
	default:
		ctxt.diag("bad fst %A", a)
	case AMOVD_5:
		o |= 1 << 8
	case AMOVF_5:
		break
	}
	return o
}

func osrr_asm5(ctxt *Link, r int, i int, b int, sc int) uint32 {
	return olr_asm5(ctxt, i, b, r, sc) ^ ((1 << 25) | (1 << 20))
}

func oshrr_asm5(ctxt *Link, r int, i int, b int, sc int) uint32 {
	return olhr_asm5(ctxt, i, b, r, sc) ^ ((1 << 22) | (1 << 20))
}

func omvl_asm5(ctxt *Link, p *Prog, a *Addr, dr int) uint32 {
	var v int
	var o1 int32
	if !(p.pcond != nil) {
		aclass_asm5(ctxt, a)
		v = int(immrot_asm5(uint32(^ctxt.instoffset)))
		if v == 0 {
			ctxt.diag("missing literal")
			prasm_asm5(p)
			return 0
		}
		o1 = int32(oprrr_asm5(ctxt, int(AMVN_5), p.scond&int(C_SCOND_5)))
		o1 |= int32(v)
		o1 |= int32(dr) << 12
	} else {
		v = p.pcond.pc - p.pc - 8
		o1 = int32(olr_asm5(ctxt, v, int(REGPC_5), dr, p.scond&int(C_SCOND_5)))
	}
	return uint32(o1)
}

func immaddr_asm5(v int32) int {
	if v >= 0 && v <= 0xfff {
		return int((v & 0xfff) | int32(1<<24) | int32(1<<23)) /* pre indexing */ /* pre indexing, up */
	}
	if v >= -0xfff && v < 0 {
		return int((-v & 0xfff) | int32(1<<24)) /* pre indexing */
	}
	return 0
}

func aclass_asm5(ctxt *Link, a *Addr) int {
	var s *LSym
	var t int
	switch a.typ {
	case D_NONE_5:
		return int(C_NONE_asm5)
	case D_REG_5:
		return int(C_REG_asm5)
	case D_REGREG_5:
		return int(C_REGREG_asm5)
	case D_REGREG2_5:
		return int(C_REGREG2_asm5)
	case D_SHIFT_5:
		return int(C_SHIFT_asm5)
	case D_FREG_5:
		return int(C_FREG_asm5)
	case D_FPCR_5:
		return int(C_FCR_asm5)
	case D_OREG_5:
		switch a.name {
		case D_EXTERN_5:
		case D_STATIC_5:
			if a.sym == nil || a.sym.name == "" {
				print("null sym external\n")
				print("%D\n", a)
				return int(C_GOK_asm5)
			}
			ctxt.instoffset = 0 // s.b. unused but just in case
			return int(C_ADDR_asm5)
		case D_AUTO_5:
			ctxt.instoffset = int32(int64(ctxt.autosize) + a.offset)
			t = immaddr_asm5(ctxt.instoffset)
			if t != 0 {
				if immhalf_asm5(ctxt.instoffset) != 0 {
					var tmp int
					if immfloat_asm5(int32(t)) {
						tmp = int(C_HFAUTO_asm5)
					} else {
						tmp = int(C_HAUTO_asm5)
					}
					return tmp
				}
				if immfloat_asm5(int32(t)) {
					return int(C_FAUTO_asm5)
				}
				return int(C_SAUTO_asm5)
			}
			return int(C_LAUTO_asm5)
		case D_PARAM_5:
			ctxt.instoffset = int32(int64(ctxt.autosize) + a.offset + 4)
			t = immaddr_asm5(ctxt.instoffset)
			if t != 0 {
				if immhalf_asm5(ctxt.instoffset) != 0 {
					var tmp int
					if immfloat_asm5(int32(t)) {
						tmp = int(C_HFAUTO_asm5)
					} else {
						tmp = int(C_HAUTO_asm5)
					}
					return tmp
				}
				if immfloat_asm5(int32(t)) {
					return int(C_FAUTO_asm5)
				}
				return int(C_SAUTO_asm5)
			}
			return int(C_LAUTO_asm5)
		case D_NONE_5:
			ctxt.instoffset = int32(a.offset)
			t = immaddr_asm5(ctxt.instoffset)
			if t != 0 {
				if immhalf_asm5(ctxt.instoffset) != 0 { /* n.b. that it will also satisfy immrot */
					var tmp int
					if immfloat_asm5(int32(t)) {
						tmp = int(C_HFOREG_asm5)
					} else {
						tmp = int(C_HOREG_asm5)
					}
					return tmp
				}
				if immfloat_asm5(int32(t)) {
					return int(C_FOREG_asm5) /* n.b. that it will also satisfy immrot */
				}
				t = int(immrot_asm5(uint32(ctxt.instoffset)))
				if t != 0 {
					return int(C_SROREG_asm5)
				}
				if immhalf_asm5(ctxt.instoffset) != 0 {
					return int(C_HOREG_asm5)
				}
				return int(C_SOREG_asm5)
			}
			t = int(immrot_asm5(uint32(ctxt.instoffset)))
			if t != 0 {
				return int(C_ROREG_asm5)
			}
			return int(C_LOREG_asm5)
		}
		return int(C_GOK_asm5)
	case D_PSR_5:
		return int(C_PSR_asm5)
	case D_OCONST_5:
		switch a.name {
		case D_EXTERN_5:
		case D_STATIC_5:
			ctxt.instoffset = 0 // s.b. unused but just in case
			return int(C_ADDR_asm5)
		}
		return int(C_GOK_asm5)
	case D_FCONST_5:
		if chipzero5(ctxt, a.u.dval) >= 0 {
			return int(C_ZFCON_asm5)
		}
		if chipfloat5(ctxt, a.u.dval) >= 0 {
			return int(C_SFCON_asm5)
		}
		return int(C_LFCON_asm5)
	case D_CONST_5:
	case D_CONST2_5:
		switch a.name {
		case D_NONE_5:
			ctxt.instoffset = int32(a.offset)
			if a.reg != int(NREG_5) {
				return aconsize(ctxt)
			}
			t = int(immrot_asm5(uint32(ctxt.instoffset)))
			if t != 0 {
				return int(C_RCON_asm5)
			}
			t = int(immrot_asm5(uint32(^ctxt.instoffset)))
			if t != 0 {
				return int(C_NCON_asm5)
			}
			return int(C_LCON_asm5)
		case D_EXTERN_5:
		case D_STATIC_5:
			s = a.sym
			if s == nil {
				break
			}
			ctxt.instoffset = 0 // s.b. unused but just in case
			return int(C_LCONADDR_asm5)
		case D_AUTO_5:
			ctxt.instoffset = int32(int64(ctxt.autosize) + a.offset)
			return aconsize(ctxt)
		case D_PARAM_5:
			ctxt.instoffset = int32(int64(ctxt.autosize) + a.offset + 4)
			return aconsize(ctxt)
		}
		return int(C_GOK_asm5)
	case D_BRANCH_5:
		return int(C_SBRA_asm5)
	}
	return int(C_GOK_asm5)
}

func aconsize(ctxt *Link) int {
	t := int(immrot_asm5(uint32(ctxt.instoffset)))
	if t != 0 {
		return int(C_RACON_asm5)
	}
	return int(C_LACON_asm5)

}

func immrot_asm5(v uint32) int32 {
	var i int
	for i = 0; i < 16; i++ {
		if (v &^ 0xff) == 0 {
			return int32(uint32(int32(i)<<8) | v | uint32(1<<25))
		}
		v = (v << 2) | (v >> 30)
	}
	return 0
}

func opbra_asm5(ctxt *Link, a int, sc int) uint32 {
	if sc&int(C_SBIT_5|C_PBIT_5|C_WBIT_5) != 0 {
		ctxt.diag(".nil/.nil/.W on bra instruction")
	}
	sc &= int(C_SCOND_5)
	if a == int(ABL_5) || a == int(ADUFFZERO_5) || a == int(ADUFFCOPY_5) {
		return (uint32(sc) << 28) | (0x5 << 25) | (0x1 << 24)
	}
	if sc != 0xe {
		ctxt.diag(".COND on bcond instruction")
	}
	switch a {
	case ABEQ_5:
		return (0x0 << 28) | (0x5 << 25)
	case ABNE_5:
		return (0x1 << 28) | (0x5 << 25)
	case ABCS_5:
		return (0x2 << 28) | (0x5 << 25)
	case ABHS_5:
		return (0x2 << 28) | (0x5 << 25)
	case ABCC_5:
		return (0x3 << 28) | (0x5 << 25)
	case ABLO_5:
		return (0x3 << 28) | (0x5 << 25)
	case ABMI_5:
		return (0x4 << 28) | (0x5 << 25)
	case ABPL_5:
		return (0x5 << 28) | (0x5 << 25)
	case ABVS_5:
		return (0x6 << 28) | (0x5 << 25)
	case ABVC_5:
		return (0x7 << 28) | (0x5 << 25)
	case ABHI_5:
		return (0x8 << 28) | (0x5 << 25)
	case ABLS_5:
		return (0x9 << 28) | (0x5 << 25)
	case ABGE_5:
		return (0xa << 28) | (0x5 << 25)
	case ABLT_5:
		return (0xb << 28) | (0x5 << 25)
	case ABGT_5:
		return (0xc << 28) | (0x5 << 25)
	case ABLE_5:
		return (0xd << 28) | (0x5 << 25)
	case AB_5:
		return (0xe << 28) | (0x5 << 25)
	}
	ctxt.diag("bad bra %A", a)
	prasm_asm5(ctxt.curp)
	return 0
}

var opcross_asm5 [8]Opcross_asm5

var oprange_asm5 [ALAST_5]Oprang_asm5

var xcmp_asm5 [C_GOK_asm5 + 1][C_GOK_asm5 + 1]int8

var repop_asm5 [ALAST_5]uint8

var zprg_asm5 = Prog{
	as:    AGOK_5,
	scond: C_SCOND_NONE_5,
	reg:   NREG_5,
	from: Addr{
		name: D_NONE_5,
		typ:  D_NONE_5,
		reg:  NREG_5,
	},
	to: Addr{
		name: D_NONE_5,
		typ:  D_NONE_5,
		reg:  NREG_5,
	},
}

var deferreturn_asm5 *LSym

func nocache_asm5(p *Prog) {
	p.optab = 0
	p.from.class = 0
	p.to.class = 0
}

/* size of a case statement including jump table */
func casesz_asm5(ctxt *Link, p *Prog) int32 {
	var jt int = 0
	var n int32 = 0
	var o *Optab_asm5
	for ; p != nil; p = p.link {
		if p.as == int(ABCASE_5) {
			jt = 1
		} else {
			if jt != 0 {
				break
			}
		}
		o = oplook_asm5(ctxt, p)
		n += int32(o.size)
	}
	return n
}

func buildop_asm5(ctxt *Link) {
	var i int
	var n int
	var r int
	for i = 0; i < int(C_GOK_asm5); i++ {
		for n = 0; n < int(C_GOK_asm5); n++ {
			xcmp_asm5[i][n] = int8(bool2int(cmp_asm5(n, i)))
		}
	}
	for n = 0; optab_asm5[n].as != int(AXXX_5); n++ {
		if (optab_asm5[n].flag & int(LPCREL_asm5)) != 0 {
			if ctxt.flag_shared != 0 {
				optab_asm5[n].size += int(optab_asm5[n].pcrelsiz)
			} else {
				optab_asm5[n].flag &^= int(LPCREL_asm5)
			}
		}
	}

	sort.Sort(ocmp_asm5(optab_asm5[:]))

	for i = 0; i < n; i++ {
		r = optab_asm5[i].as
		oprange_asm5[r].start = optab_asm5[i:][0:]
		for optab_asm5[i].as == r {
			i++
		}
		oprange_asm5[r].stop = optab_asm5[i:][0:]
		i--
		switch r {
		default:
			ctxt.diag("unknown op in build: %A", r)
			sysfatal("bad code")
		case AADD_5:
			oprange_asm5[AAND_5] = oprange_asm5[r]
			oprange_asm5[AEOR_5] = oprange_asm5[r]
			oprange_asm5[ASUB_5] = oprange_asm5[r]
			oprange_asm5[ARSB_5] = oprange_asm5[r]
			oprange_asm5[AADC_5] = oprange_asm5[r]
			oprange_asm5[ASBC_5] = oprange_asm5[r]
			oprange_asm5[ARSC_5] = oprange_asm5[r]
			oprange_asm5[AORR_5] = oprange_asm5[r]
			oprange_asm5[ABIC_5] = oprange_asm5[r]
			break
		case ACMP_5:
			oprange_asm5[ATEQ_5] = oprange_asm5[r]
			oprange_asm5[ACMN_5] = oprange_asm5[r]
			break
		case AMVN_5:
			break
		case ABEQ_5:
			oprange_asm5[ABNE_5] = oprange_asm5[r]
			oprange_asm5[ABCS_5] = oprange_asm5[r]
			oprange_asm5[ABHS_5] = oprange_asm5[r]
			oprange_asm5[ABCC_5] = oprange_asm5[r]
			oprange_asm5[ABLO_5] = oprange_asm5[r]
			oprange_asm5[ABMI_5] = oprange_asm5[r]
			oprange_asm5[ABPL_5] = oprange_asm5[r]
			oprange_asm5[ABVS_5] = oprange_asm5[r]
			oprange_asm5[ABVC_5] = oprange_asm5[r]
			oprange_asm5[ABHI_5] = oprange_asm5[r]
			oprange_asm5[ABLS_5] = oprange_asm5[r]
			oprange_asm5[ABGE_5] = oprange_asm5[r]
			oprange_asm5[ABLT_5] = oprange_asm5[r]
			oprange_asm5[ABGT_5] = oprange_asm5[r]
			oprange_asm5[ABLE_5] = oprange_asm5[r]
			break
		case ASLL_5:
			oprange_asm5[ASRL_5] = oprange_asm5[r]
			oprange_asm5[ASRA_5] = oprange_asm5[r]
			break
		case AMUL_5:
			oprange_asm5[AMULU_5] = oprange_asm5[r]
			break
		case ADIV_5:
			oprange_asm5[AMOD_5] = oprange_asm5[r]
			oprange_asm5[AMODU_5] = oprange_asm5[r]
			oprange_asm5[ADIVU_5] = oprange_asm5[r]
			break
		case AMOVW_5:
		case AMOVB_5:
		case AMOVBS_5:
		case AMOVBU_5:
		case AMOVH_5:
		case AMOVHS_5:
		case AMOVHU_5:
			break
		case ASWPW_5:
			oprange_asm5[ASWPBU_5] = oprange_asm5[r]
			break
		case AB_5:
		case ABL_5:
		case ABX_5:
		case ABXRET_5:
		case ADUFFZERO_5:
		case ADUFFCOPY_5:
		case ASWI_5:
		case AWORD_5:
		case AMOVM_5:
		case ARFE_5:
		case ATEXT_5:
		case AUSEFIELD_5:
		case ACASE_5:
		case ABCASE_5:
		case ATYPE_5:
			break
		case AADDF_5:
			oprange_asm5[AADDD_5] = oprange_asm5[r]
			oprange_asm5[ASUBF_5] = oprange_asm5[r]
			oprange_asm5[ASUBD_5] = oprange_asm5[r]
			oprange_asm5[AMULF_5] = oprange_asm5[r]
			oprange_asm5[AMULD_5] = oprange_asm5[r]
			oprange_asm5[ADIVF_5] = oprange_asm5[r]
			oprange_asm5[ADIVD_5] = oprange_asm5[r]
			oprange_asm5[ASQRTF_5] = oprange_asm5[r]
			oprange_asm5[ASQRTD_5] = oprange_asm5[r]
			oprange_asm5[AMOVFD_5] = oprange_asm5[r]
			oprange_asm5[AMOVDF_5] = oprange_asm5[r]
			oprange_asm5[AABSF_5] = oprange_asm5[r]
			oprange_asm5[AABSD_5] = oprange_asm5[r]
			break
		case ACMPF_5:
			oprange_asm5[ACMPD_5] = oprange_asm5[r]
			break
		case AMOVF_5:
			oprange_asm5[AMOVD_5] = oprange_asm5[r]
			break
		case AMOVFW_5:
			oprange_asm5[AMOVDW_5] = oprange_asm5[r]
			break
		case AMOVWF_5:
			oprange_asm5[AMOVWD_5] = oprange_asm5[r]
			break
		case AMULL_5:
			oprange_asm5[AMULAL_5] = oprange_asm5[r]
			oprange_asm5[AMULLU_5] = oprange_asm5[r]
			oprange_asm5[AMULALU_5] = oprange_asm5[r]
			break
		case AMULWT_5:
			oprange_asm5[AMULWB_5] = oprange_asm5[r]
			break
		case AMULAWT_5:
			oprange_asm5[AMULAWB_5] = oprange_asm5[r]
			break
		case AMULA_5:
		case ALDREX_5:
		case ASTREX_5:
		case ALDREXD_5:
		case ASTREXD_5:
		case ATST_5:
		case APLD_5:
		case AUNDEF_5:
		case ACLZ_5:
		case AFUNCDATA_5:
		case APCDATA_5:
		case ADATABUNDLE_5:
		case ADATABUNDLEEND_5:
			break
		}
	}
}

func regoff_asm5(ctxt *Link, a *Addr) int32 {
	ctxt.instoffset = 0
	aclass_asm5(ctxt, a)
	return ctxt.instoffset
}

func immfloat_asm5(v int32) bool {
	return (v & 0xC03) == 0 /* offset will fit in floating-point load/store */
}

func immhalf_asm5(v int32) int {
	if v >= 0 && v <= 0xff {
		return int(v | int32(1<<24) | int32(1<<23)) /* pre indexing */ /* pre indexing, up */
	}
	if v >= -0xff && v < 0 {
		return int((-v & 0xff) | int32(1<<24)) /* pre indexing */
	}
	return 0
}

func prasm_asm5(p *Prog) {
	print("%P\n", p)
}

func cmp_asm5(a int, b int) bool {
	var zz bool
	if a == b {
		return true
	}
	switch a {
	case C_LCON_asm5:
		if b == int(C_RCON_asm5) || b == int(C_NCON_asm5) {
			return true
		}
		break
	case C_LACON_asm5:
		if b == int(C_RACON_asm5) {
			return true
		}
		break
	case C_LFCON_asm5:
		if b == int(C_ZFCON_asm5) || b == int(C_SFCON_asm5) {
			return true
		}
		break
	case C_HFAUTO_asm5:
		return b == int(C_HAUTO_asm5) || b == int(C_FAUTO_asm5)
	case C_FAUTO_asm5:
	case C_HAUTO_asm5:
		zz = (b == int(C_HFAUTO_asm5))
		return zz
	case C_SAUTO_asm5:
		return cmp_asm5(int(C_HFAUTO_asm5), b)
	case C_LAUTO_asm5:
		return cmp_asm5(int(C_SAUTO_asm5), b)
	case C_HFOREG_asm5:
		return b == int(C_HOREG_asm5) || b == int(C_FOREG_asm5)
	case C_FOREG_asm5:
	case C_HOREG_asm5:
		return b == int(C_HFOREG_asm5)
	case C_SROREG_asm5:
		return cmp_asm5(int(C_SOREG_asm5), b) || cmp_asm5(int(C_ROREG_asm5), b)
	case C_SOREG_asm5:
	case C_ROREG_asm5:
		return b == int(C_SROREG_asm5) || cmp_asm5(int(C_HFOREG_asm5), b)
	case C_LOREG_asm5:
		return cmp_asm5(int(C_SROREG_asm5), b)
	case C_LBRA_asm5:
		if b == int(C_SBRA_asm5) {
			return true
		}
		break
	case C_HREG_asm5:
		return cmp_asm5(int(C_SP_asm5), b) || cmp_asm5(int(C_PC_asm5), b)
	}
	return false
}

type ocmp_asm5 []Optab_asm5

func (x ocmp_asm5) Len() int      { return len(x) }
func (x ocmp_asm5) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x ocmp_asm5) Less(i, j int) bool {
	p1 := &x[i]
	p2 := &x[j]

	var n int32
	n = int32(p1.as) - int32(p2.as)
	if n != 0 {
		return n < 0
	}
	n = int32(p1.a1) - int32(p2.a1)
	if n != 0 {
		return n < 0
	}
	n = int32(p1.a2) - int32(p2.a2)
	if n != 0 {
		return n < 0
	}
	n = int32(p1.a3) - int32(p2.a3)
	if n != 0 {
		return n < 0
	}
	return false
}
