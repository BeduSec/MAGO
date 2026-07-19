// Copyright (c) BeduSec. All rights reserved.
#include "textflag.h"

TEXT ·takeToken(SB), NOSPLIT, $0-8
	MOVQ counter+0(FP), AX
	MOVQ (AX), CX
	CMPQ CX, $0
	JE   fail
	DECQ CX
	MOVQ CX, (AX)
	MOVQ CX, ret+8(FP)
	RET
fail:
	MOVQ $0, ret+8(FP)
	RET