  "".min STEXT nosplit size=1 args=0x0 locals=0x0 funcid=0x0 align=0x0
	"".min(SB), NOSPLIT|ABIInternal, $0-0
#	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
#	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
  RET                                              .
"".arg1 STEXT nosplit size=6 args=0x8 locals=0x0 funcid=0x0 align=0x0
	TEXT	"".arg1(SB), NOSPLIT|ABIInternal, $0-8
#	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
#	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
#	FUNCDATA	$5, "".arg1.arginfo1(SB)
	MOVQ	AX, "".x+8(SP)
	RET
"".arg1ret1 STEXT nosplit size=46 args=0x8 locals=0x10 funcid=0x0 align=0x0
	TEXT	"".arg1ret1(SB), NOSPLIT|ABIInternal, $16-8
	SUBQ	$16, SP
	MOVQ	BP, 8(SP)
	LEAQ	8(SP), BP
#	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
#	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
#	FUNCDATA	$5, "".arg1ret1.arginfo1(SB)
	MOVQ	AX, "".x+24(SP)
	MOVQ	$0, "".~r0(SP)
	MOVQ	"".x+24(SP), AX
	MOVQ	AX, "".~r0(SP)
	MOVQ	8(SP), BP
	ADDQ	$16, SP
	RETå
"".main STEXT nosplit size=56 args=0x0 locals=0x20 funcid=0x0 align=0x0
	TEXT	"".main(SB), NOSPLIT|ABIInternal, $32-0
	SUBQ	$32, SP
	MOVQ	BP, 24(SP)
	LEAQ	24(SP), BP
#	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
#	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	JMP	16
	MOVQ	$1, "".x+16(SP)
	JMP	27
	MOVQ	$2, "".x+8(SP)
	MOVQ	$2, "".~R0(SP)
	JMP	46
	MOVQ	24(SP), BP
	ADDQ	$32, SP
	RET