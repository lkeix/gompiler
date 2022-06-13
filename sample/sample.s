"".min STEXT nosplit size=1 args=0x0 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (./main.go:3)	TEXT	"".min(SB), NOSPLIT|ABIInternal, $0-0
	0x0000 00000 (./main.go:3)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (./main.go:3)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (./main.go:5)	RET
	0x0000 c3                                               .
"".arg1 STEXT nosplit size=6 args=0x8 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (./main.go:7)	TEXT	"".arg1(SB), NOSPLIT|ABIInternal, $0-8
	0x0000 00000 (./main.go:7)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (./main.go:7)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (./main.go:7)	FUNCDATA	$5, "".arg1.arginfo1(SB)
	0x0000 00000 (./main.go:7)	MOVQ	AX, "".x+8(SP)
	0x0005 00005 (./main.go:9)	RET
	0x0000 48 89 44 24 08 c3                                H.D$..
"".arg1ret1 STEXT nosplit size=46 args=0x8 locals=0x10 funcid=0x0 align=0x0
	0x0000 00000 (./main.go:11)	TEXT	"".arg1ret1(SB), NOSPLIT|ABIInternal, $16-8
	0x0000 00000 (./main.go:11)	SUBQ	$16, SP
	0x0004 00004 (./main.go:11)	MOVQ	BP, 8(SP)
	0x0009 00009 (./main.go:11)	LEAQ	8(SP), BP
	0x000e 00014 (./main.go:11)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000e 00014 (./main.go:11)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000e 00014 (./main.go:11)	FUNCDATA	$5, "".arg1ret1.arginfo1(SB)
	0x000e 00014 (./main.go:11)	MOVQ	AX, "".x+24(SP)
	0x0013 00019 (./main.go:11)	MOVQ	$0, "".~r0(SP)
	0x001b 00027 (./main.go:12)	MOVQ	"".x+24(SP), AX
	0x0020 00032 (./main.go:12)	MOVQ	AX, "".~r0(SP)
	0x0024 00036 (./main.go:12)	MOVQ	8(SP), BP
	0x0029 00041 (./main.go:12)	ADDQ	$16, SP
	0x002d 00045 (./main.go:12)	RET
	0x0000 48 83 ec 10 48 89 6c 24 08 48 8d 6c 24 08 48 89  H...H.l$.H.l$.H.
	0x0010 44 24 18 48 c7 04 24 00 00 00 00 48 8b 44 24 18  D$.H..$....H.D$.
	0x0020 48 89 04 24 48 8b 6c 24 08 48 83 c4 10 c3        H..$H.l$.H....
"".main STEXT nosplit size=56 args=0x0 locals=0x20 funcid=0x0 align=0x0
	0x0000 00000 (./main.go:15)	TEXT	"".main(SB), NOSPLIT|ABIInternal, $32-0
	0x0000 00000 (./main.go:15)	SUBQ	$32, SP
	0x0004 00004 (./main.go:15)	MOVQ	BP, 24(SP)
	0x0009 00009 (./main.go:15)	LEAQ	24(SP), BP
	0x000e 00014 (./main.go:15)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000e 00014 (./main.go:15)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000e 00014 (./main.go:16)	JMP	16
	0x0010 00016 (./main.go:17)	MOVQ	$1, "".x+16(SP)
	0x0019 00025 (./main.go:17)	JMP	27
	0x001b 00027 (./main.go:18)	MOVQ	$2, "".x+8(SP)
	0x0024 00036 (./main.go:18)	MOVQ	$2, "".~R0(SP)
	0x002c 00044 (./main.go:18)	JMP	46
	0x002e 00046 (./main.go:19)	MOVQ	24(SP), BP
	0x0033 00051 (./main.go:19)	ADDQ	$32, SP
	0x0037 00055 (./main.go:19)	RET
	0x0000 48 83 ec 20 48 89 6c 24 18 48 8d 6c 24 18 eb 00  H.. H.l$.H.l$...
	0x0010 48 c7 44 24 10 01 00 00 00 eb 00 48 c7 44 24 08  H.D$.......H.D$.
	0x0020 02 00 00 00 48 c7 04 24 02 00 00 00 eb 00 48 8b  ....H..$......H.
	0x0030 6c 24 18 48 83 c4 20 c3                          l$.H.. .
go.cuinfo.packagename. SDWARFCUINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
go.info."".min$abstract SDWARFABSFCN dupok size=9
	0x0000 05 2e 6d 69 6e 00 01 01 00                       ..min....
go.info."".arg1$abstract SDWARFABSFCN dupok size=18
	0x0000 05 2e 61 72 67 31 00 01 01 13 78 00 00 00 00 00  ..arg1....x.....
	0x0010 00 00                                            ..
	rel 0+0 t=22 type.int+0
	rel 13+4 t=31 go.info.int+0
go.info."".arg1ret1$abstract SDWARFABSFCN dupok size=22
	0x0000 05 2e 61 72 67 31 72 65 74 31 00 01 01 13 78 00  ..arg1ret1....x.
	0x0010 00 00 00 00 00 00                                ......
	rel 0+0 t=22 type.int+0
	rel 17+4 t=31 go.info.int+0
""..inittask SNOPTRDATA size=24
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 00 00 00                          ........
gclocals·33cdeccccebe80329f1fdbee7f5874cb SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
"".arg1.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
"".arg1ret1.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
