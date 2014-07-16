package main

// Inferno utils/6l/span.c
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
// Instruction layout.
// Inferno utils/6c/6.out.h
// http://code.google.com/p/inferno-os/source/browse/utils/6c/6.out.h
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
/*
 *	amd64
 */
const (
	AXXX_6 = iota
	AAAA_6
	AAAD_6
	AAAM_6
	AAAS_6
	AADCB_6
	AADCL_6
	AADCW_6
	AADDB_6
	AADDL_6
	AADDW_6
	AADJSP_6
	AANDB_6
	AANDL_6
	AANDW_6
	AARPL_6
	ABOUNDL_6
	ABOUNDW_6
	ABSFL_6
	ABSFW_6
	ABSRL_6
	ABSRW_6
	ABTL_6
	ABTW_6
	ABTCL_6
	ABTCW_6
	ABTRL_6
	ABTRW_6
	ABTSL_6
	ABTSW_6
	ABYTE_6
	ACALL_6
	ACLC_6
	ACLD_6
	ACLI_6
	ACLTS_6
	ACMC_6
	ACMPB_6
	ACMPL_6
	ACMPW_6
	ACMPSB_6
	ACMPSL_6
	ACMPSW_6
	ADAA_6
	ADAS_6
	ADATA_6
	ADECB_6
	ADECL_6
	ADECQ_6
	ADECW_6
	ADIVB_6
	ADIVL_6
	ADIVW_6
	AENTER_6
	AGLOBL_6
	AGOK_6
	AHISTORY_6
	AHLT_6
	AIDIVB_6
	AIDIVL_6
	AIDIVW_6
	AIMULB_6
	AIMULL_6
	AIMULW_6
	AINB_6
	AINL_6
	AINW_6
	AINCB_6
	AINCL_6
	AINCQ_6
	AINCW_6
	AINSB_6
	AINSL_6
	AINSW_6
	AINT_6
	AINTO_6
	AIRETL_6
	AIRETW_6
	AJCC_6
	AJCS_6
	AJCXZL_6
	AJEQ_6
	AJGE_6
	AJGT_6
	AJHI_6
	AJLE_6
	AJLS_6
	AJLT_6
	AJMI_6
	AJMP_6
	AJNE_6
	AJOC_6
	AJOS_6
	AJPC_6
	AJPL_6
	AJPS_6
	ALAHF_6
	ALARL_6
	ALARW_6
	ALEAL_6
	ALEAW_6
	ALEAVEL_6
	ALEAVEW_6
	ALOCK_6
	ALODSB_6
	ALODSL_6
	ALODSW_6
	ALONG_6
	ALOOP_6
	ALOOPEQ_6
	ALOOPNE_6
	ALSLL_6
	ALSLW_6
	AMOVB_6
	AMOVL_6
	AMOVW_6
	AMOVBLSX_6
	AMOVBLZX_6
	AMOVBQSX_6
	AMOVBQZX_6
	AMOVBWSX_6
	AMOVBWZX_6
	AMOVWLSX_6
	AMOVWLZX_6
	AMOVWQSX_6
	AMOVWQZX_6
	AMOVSB_6
	AMOVSL_6
	AMOVSW_6
	AMULB_6
	AMULL_6
	AMULW_6
	ANAME_6
	ANEGB_6
	ANEGL_6
	ANEGW_6
	ANOP_6
	ANOTB_6
	ANOTL_6
	ANOTW_6
	AORB_6
	AORL_6
	AORW_6
	AOUTB_6
	AOUTL_6
	AOUTW_6
	AOUTSB_6
	AOUTSL_6
	AOUTSW_6
	APAUSE_6
	APOPAL_6
	APOPAW_6
	APOPFL_6
	APOPFW_6
	APOPL_6
	APOPW_6
	APUSHAL_6
	APUSHAW_6
	APUSHFL_6
	APUSHFW_6
	APUSHL_6
	APUSHW_6
	ARCLB_6
	ARCLL_6
	ARCLW_6
	ARCRB_6
	ARCRL_6
	ARCRW_6
	AREP_6
	AREPN_6
	ARET_6
	AROLB_6
	AROLL_6
	AROLW_6
	ARORB_6
	ARORL_6
	ARORW_6
	ASAHF_6
	ASALB_6
	ASALL_6
	ASALW_6
	ASARB_6
	ASARL_6
	ASARW_6
	ASBBB_6
	ASBBL_6
	ASBBW_6
	ASCASB_6
	ASCASL_6
	ASCASW_6
	ASETCC_6
	ASETCS_6
	ASETEQ_6
	ASETGE_6
	ASETGT_6
	ASETHI_6
	ASETLE_6
	ASETLS_6
	ASETLT_6
	ASETMI_6
	ASETNE_6
	ASETOC_6
	ASETOS_6
	ASETPC_6
	ASETPL_6
	ASETPS_6
	ACDQ_6
	ACWD_6
	ASHLB_6
	ASHLL_6
	ASHLW_6
	ASHRB_6
	ASHRL_6
	ASHRW_6
	ASTC_6
	ASTD_6
	ASTI_6
	ASTOSB_6
	ASTOSL_6
	ASTOSW_6
	ASUBB_6
	ASUBL_6
	ASUBW_6
	ASYSCALL_6
	ATESTB_6
	ATESTL_6
	ATESTW_6
	ATEXT_6
	AVERR_6
	AVERW_6
	AWAIT_6
	AWORD_6
	AXCHGB_6
	AXCHGL_6
	AXCHGW_6
	AXLAT_6
	AXORB_6
	AXORL_6
	AXORW_6
	AFMOVB_6
	AFMOVBP_6
	AFMOVD_6
	AFMOVDP_6
	AFMOVF_6
	AFMOVFP_6
	AFMOVL_6
	AFMOVLP_6
	AFMOVV_6
	AFMOVVP_6
	AFMOVW_6
	AFMOVWP_6
	AFMOVX_6
	AFMOVXP_6
	AFCOMB_6
	AFCOMBP_6
	AFCOMD_6
	AFCOMDP_6
	AFCOMDPP_6
	AFCOMF_6
	AFCOMFP_6
	AFCOML_6
	AFCOMLP_6
	AFCOMW_6
	AFCOMWP_6
	AFUCOM_6
	AFUCOMP_6
	AFUCOMPP_6
	AFADDDP_6
	AFADDW_6
	AFADDL_6
	AFADDF_6
	AFADDD_6
	AFMULDP_6
	AFMULW_6
	AFMULL_6
	AFMULF_6
	AFMULD_6
	AFSUBDP_6
	AFSUBW_6
	AFSUBL_6
	AFSUBF_6
	AFSUBD_6
	AFSUBRDP_6
	AFSUBRW_6
	AFSUBRL_6
	AFSUBRF_6
	AFSUBRD_6
	AFDIVDP_6
	AFDIVW_6
	AFDIVL_6
	AFDIVF_6
	AFDIVD_6
	AFDIVRDP_6
	AFDIVRW_6
	AFDIVRL_6
	AFDIVRF_6
	AFDIVRD_6
	AFXCHD_6
	AFFREE_6
	AFLDCW_6
	AFLDENV_6
	AFRSTOR_6
	AFSAVE_6
	AFSTCW_6
	AFSTENV_6
	AFSTSW_6
	AF2XM1_6
	AFABS_6
	AFCHS_6
	AFCLEX_6
	AFCOS_6
	AFDECSTP_6
	AFINCSTP_6
	AFINIT_6
	AFLD1_6
	AFLDL2E_6
	AFLDL2T_6
	AFLDLG2_6
	AFLDLN2_6
	AFLDPI_6
	AFLDZ_6
	AFNOP_6
	AFPATAN_6
	AFPREM_6
	AFPREM1_6
	AFPTAN_6
	AFRNDINT_6
	AFSCALE_6
	AFSIN_6
	AFSINCOS_6
	AFSQRT_6
	AFTST_6
	AFXAM_6
	AFXTRACT_6
	AFYL2X_6
	AFYL2XP1_6
	AEND_6
	ADYNT__6
	AINIT__6
	ASIGNAME_6
	ACMPXCHGB_6
	ACMPXCHGL_6
	ACMPXCHGW_6
	ACMPXCHG8B_6
	ACPUID_6
	AINVD_6
	AINVLPG_6
	ALFENCE_6
	AMFENCE_6
	AMOVNTIL_6
	ARDMSR_6
	ARDPMC_6
	ARDTSC_6
	ARSM_6
	ASFENCE_6
	ASYSRET_6
	AWBINVD_6
	AWRMSR_6
	AXADDB_6
	AXADDL_6
	AXADDW_6
	ACMOVLCC_6
	ACMOVLCS_6
	ACMOVLEQ_6
	ACMOVLGE_6
	ACMOVLGT_6
	ACMOVLHI_6
	ACMOVLLE_6
	ACMOVLLS_6
	ACMOVLLT_6
	ACMOVLMI_6
	ACMOVLNE_6
	ACMOVLOC_6
	ACMOVLOS_6
	ACMOVLPC_6
	ACMOVLPL_6
	ACMOVLPS_6
	ACMOVQCC_6
	ACMOVQCS_6
	ACMOVQEQ_6
	ACMOVQGE_6
	ACMOVQGT_6
	ACMOVQHI_6
	ACMOVQLE_6
	ACMOVQLS_6
	ACMOVQLT_6
	ACMOVQMI_6
	ACMOVQNE_6
	ACMOVQOC_6
	ACMOVQOS_6
	ACMOVQPC_6
	ACMOVQPL_6
	ACMOVQPS_6
	ACMOVWCC_6
	ACMOVWCS_6
	ACMOVWEQ_6
	ACMOVWGE_6
	ACMOVWGT_6
	ACMOVWHI_6
	ACMOVWLE_6
	ACMOVWLS_6
	ACMOVWLT_6
	ACMOVWMI_6
	ACMOVWNE_6
	ACMOVWOC_6
	ACMOVWOS_6
	ACMOVWPC_6
	ACMOVWPL_6
	ACMOVWPS_6
	AADCQ_6
	AADDQ_6
	AANDQ_6
	ABSFQ_6
	ABSRQ_6
	ABTCQ_6
	ABTQ_6
	ABTRQ_6
	ABTSQ_6
	ACMPQ_6
	ACMPSQ_6
	ACMPXCHGQ_6
	ACQO_6
	ADIVQ_6
	AIDIVQ_6
	AIMULQ_6
	AIRETQ_6
	AJCXZQ_6
	ALEAQ_6
	ALEAVEQ_6
	ALODSQ_6
	AMOVQ_6
	AMOVLQSX_6
	AMOVLQZX_6
	AMOVNTIQ_6
	AMOVSQ_6
	AMULQ_6
	ANEGQ_6
	ANOTQ_6
	AORQ_6
	APOPFQ_6
	APOPQ_6
	APUSHFQ_6
	APUSHQ_6
	ARCLQ_6
	ARCRQ_6
	AROLQ_6
	ARORQ_6
	AQUAD_6
	ASALQ_6
	ASARQ_6
	ASBBQ_6
	ASCASQ_6
	ASHLQ_6
	ASHRQ_6
	ASTOSQ_6
	ASUBQ_6
	ATESTQ_6
	AXADDQ_6
	AXCHGQ_6
	AXORQ_6
	AADDPD_6
	AADDPS_6
	AADDSD_6
	AADDSS_6
	AANDNPD_6
	AANDNPS_6
	AANDPD_6
	AANDPS_6
	ACMPPD_6
	ACMPPS_6
	ACMPSD_6
	ACMPSS_6
	ACOMISD_6
	ACOMISS_6
	ACVTPD2PL_6
	ACVTPD2PS_6
	ACVTPL2PD_6
	ACVTPL2PS_6
	ACVTPS2PD_6
	ACVTPS2PL_6
	ACVTSD2SL_6
	ACVTSD2SQ_6
	ACVTSD2SS_6
	ACVTSL2SD_6
	ACVTSL2SS_6
	ACVTSQ2SD_6
	ACVTSQ2SS_6
	ACVTSS2SD_6
	ACVTSS2SL_6
	ACVTSS2SQ_6
	ACVTTPD2PL_6
	ACVTTPS2PL_6
	ACVTTSD2SL_6
	ACVTTSD2SQ_6
	ACVTTSS2SL_6
	ACVTTSS2SQ_6
	ADIVPD_6
	ADIVPS_6
	ADIVSD_6
	ADIVSS_6
	AEMMS_6
	AFXRSTOR_6
	AFXRSTOR64_6
	AFXSAVE_6
	AFXSAVE64_6
	ALDMXCSR_6
	AMASKMOVOU_6
	AMASKMOVQ_6
	AMAXPD_6
	AMAXPS_6
	AMAXSD_6
	AMAXSS_6
	AMINPD_6
	AMINPS_6
	AMINSD_6
	AMINSS_6
	AMOVAPD_6
	AMOVAPS_6
	AMOVOU_6
	AMOVHLPS_6
	AMOVHPD_6
	AMOVHPS_6
	AMOVLHPS_6
	AMOVLPD_6
	AMOVLPS_6
	AMOVMSKPD_6
	AMOVMSKPS_6
	AMOVNTO_6
	AMOVNTPD_6
	AMOVNTPS_6
	AMOVNTQ_6
	AMOVO_6
	AMOVQOZX_6
	AMOVSD_6
	AMOVSS_6
	AMOVUPD_6
	AMOVUPS_6
	AMULPD_6
	AMULPS_6
	AMULSD_6
	AMULSS_6
	AORPD_6
	AORPS_6
	APACKSSLW_6
	APACKSSWB_6
	APACKUSWB_6
	APADDB_6
	APADDL_6
	APADDQ_6
	APADDSB_6
	APADDSW_6
	APADDUSB_6
	APADDUSW_6
	APADDW_6
	APANDB_6
	APANDL_6
	APANDSB_6
	APANDSW_6
	APANDUSB_6
	APANDUSW_6
	APANDW_6
	APAND_6
	APANDN_6
	APAVGB_6
	APAVGW_6
	APCMPEQB_6
	APCMPEQL_6
	APCMPEQW_6
	APCMPGTB_6
	APCMPGTL_6
	APCMPGTW_6
	APEXTRW_6
	APFACC_6
	APFADD_6
	APFCMPEQ_6
	APFCMPGE_6
	APFCMPGT_6
	APFMAX_6
	APFMIN_6
	APFMUL_6
	APFNACC_6
	APFPNACC_6
	APFRCP_6
	APFRCPIT1_6
	APFRCPI2T_6
	APFRSQIT1_6
	APFRSQRT_6
	APFSUB_6
	APFSUBR_6
	APINSRW_6
	APINSRD_6
	APINSRQ_6
	APMADDWL_6
	APMAXSW_6
	APMAXUB_6
	APMINSW_6
	APMINUB_6
	APMOVMSKB_6
	APMULHRW_6
	APMULHUW_6
	APMULHW_6
	APMULLW_6
	APMULULQ_6
	APOR_6
	APSADBW_6
	APSHUFHW_6
	APSHUFL_6
	APSHUFLW_6
	APSHUFW_6
	APSHUFB_6
	APSLLO_6
	APSLLL_6
	APSLLQ_6
	APSLLW_6
	APSRAL_6
	APSRAW_6
	APSRLO_6
	APSRLL_6
	APSRLQ_6
	APSRLW_6
	APSUBB_6
	APSUBL_6
	APSUBQ_6
	APSUBSB_6
	APSUBSW_6
	APSUBUSB_6
	APSUBUSW_6
	APSUBW_6
	APSWAPL_6
	APUNPCKHBW_6
	APUNPCKHLQ_6
	APUNPCKHQDQ_6
	APUNPCKHWL_6
	APUNPCKLBW_6
	APUNPCKLLQ_6
	APUNPCKLQDQ_6
	APUNPCKLWL_6
	APXOR_6
	ARCPPS_6
	ARCPSS_6
	ARSQRTPS_6
	ARSQRTSS_6
	ASHUFPD_6
	ASHUFPS_6
	ASQRTPD_6
	ASQRTPS_6
	ASQRTSD_6
	ASQRTSS_6
	ASTMXCSR_6
	ASUBPD_6
	ASUBPS_6
	ASUBSD_6
	ASUBSS_6
	AUCOMISD_6
	AUCOMISS_6
	AUNPCKHPD_6
	AUNPCKHPS_6
	AUNPCKLPD_6
	AUNPCKLPS_6
	AXORPD_6
	AXORPS_6
	APF2IW_6
	APF2IL_6
	API2FW_6
	API2FL_6
	ARETFW_6
	ARETFL_6
	ARETFQ_6
	ASWAPGS_6
	AMODE_6
	ACRC32B_6
	ACRC32Q_6
	AIMUL3Q_6
	APREFETCHT0_6
	APREFETCHT1_6
	APREFETCHT2_6
	APREFETCHNTA_6
	AMOVQL_6
	ABSWAPL_6
	ABSWAPQ_6
	AUNDEF_6
	AAESENC_6
	AAESENCLAST_6
	AAESDEC_6
	AAESDECLAST_6
	AAESIMC_6
	AAESKEYGENASSIST_6
	APSHUFD_6
	APCLMULQDQ_6
	AUSEFIELD_6
	ATYPE_6
	AFUNCDATA_6
	APCDATA_6
	ACHECKNIL_6
	AVARDEF_6
	AVARKILL_6
	ADUFFCOPY_6
	ADUFFZERO_6
	ALAST_6
)

const (
	D_AL_6 = 0 + iota
	D_CL_6
	D_DL_6
	D_BL_6
	D_SPB_6
	D_BPB_6
	D_SIB_6
	D_DIB_6
	D_R8B_6
	D_R9B_6
	D_R10B_6
	D_R11B_6
	D_R12B_6
	D_R13B_6
	D_R14B_6
	D_R15B_6
	D_AX_6 = 16 + iota - 16
	D_CX_6
	D_DX_6
	D_BX_6
	D_SP_6
	D_BP_6
	D_SI_6
	D_DI_6
	D_R8_6
	D_R9_6
	D_R10_6
	D_R11_6
	D_R12_6
	D_R13_6
	D_R14_6
	D_R15_6
	D_AH_6 = 32 + iota - 32
	D_CH_6
	D_DH_6
	D_BH_6
	D_F0_6 = 36
	D_M0_6 = 44
	D_X0_6 = 52 + iota - 38
	D_X1_6
	D_X2_6
	D_X3_6
	D_X4_6
	D_X5_6
	D_X6_6
	D_X7_6
	D_X8_6
	D_X9_6
	D_X10_6
	D_X11_6
	D_X12_6
	D_X13_6
	D_X14_6
	D_X15_6
	D_CS_6 = 68 + iota - 54
	D_SS_6
	D_DS_6
	D_ES_6
	D_FS_6
	D_GS_6
	D_GDTR_6
	D_IDTR_6
	D_LDTR_6
	D_MSW_6
	D_TASK_6
	D_CR_6     = 79
	D_DR_6     = 95
	D_TR_6     = 103
	D_TLS_6    = 111
	D_NONE_6   = 112
	D_BRANCH_6 = 113
	D_EXTERN_6 = 114
	D_STATIC_6 = 115
	D_AUTO_6   = 116
	D_PARAM_6  = 117
	D_CONST_6  = 118
	D_FCONST_6 = 119
	D_SCONST_6 = 120
	D_ADDR_6   = 121 + iota - 78
	D_INDIR_6
	T_TYPE_6   = 1 << 0
	T_INDEX_6  = 1 << 1
	T_OFFSET_6 = 1 << 2
	T_FCONST_6 = 1 << 3
	T_SYM_6    = 1 << 4
	T_SCONST_6 = 1 << 5
	T_64_6     = 1 << 6
	T_GOTYPE_6 = 1 << 7
	REGARG_6   = -1
	REGRET_6   = D_AX_6
	FREGRET_6  = D_X0_6
	REGSP_6    = D_SP_6
	REGTMP_6   = D_DI_6
	REGEXT_6   = D_R15_6
	FREGMIN_6  = D_X0_6 + 5
	FREGEXT_6  = D_X0_6 + 15
)
