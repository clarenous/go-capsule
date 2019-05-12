package capsvm

type Op uint8

type Instruction struct {
	Op   Op
	Len  uint32
	Data []byte
}

const (
	OP_FALSE Op = 0x00
	OP_0     Op = 0x00 // synonym

	OP_1    Op = 0x51
	OP_TRUE Op = 0x51 // synonym

	OP_2  Op = 0x52
	OP_3  Op = 0x53
	OP_4  Op = 0x54
	OP_5  Op = 0x55
	OP_6  Op = 0x56
	OP_7  Op = 0x57
	OP_8  Op = 0x58
	OP_9  Op = 0x59
	OP_10 Op = 0x5a
	OP_11 Op = 0x5b
	OP_12 Op = 0x5c
	OP_13 Op = 0x5d
	OP_14 Op = 0x5e
	OP_15 Op = 0x5f
	OP_16 Op = 0x60

	OP_DATA_1  Op = 0x01
	OP_DATA_2  Op = 0x02
	OP_DATA_3  Op = 0x03
	OP_DATA_4  Op = 0x04
	OP_DATA_5  Op = 0x05
	OP_DATA_6  Op = 0x06
	OP_DATA_7  Op = 0x07
	OP_DATA_8  Op = 0x08
	OP_DATA_9  Op = 0x09
	OP_DATA_10 Op = 0x0a
	OP_DATA_11 Op = 0x0b
	OP_DATA_12 Op = 0x0c
	OP_DATA_13 Op = 0x0d
	OP_DATA_14 Op = 0x0e
	OP_DATA_15 Op = 0x0f
	OP_DATA_16 Op = 0x10
	OP_DATA_17 Op = 0x11
	OP_DATA_18 Op = 0x12
	OP_DATA_19 Op = 0x13
	OP_DATA_20 Op = 0x14
	OP_DATA_21 Op = 0x15
	OP_DATA_22 Op = 0x16
	OP_DATA_23 Op = 0x17
	OP_DATA_24 Op = 0x18
	OP_DATA_25 Op = 0x19
	OP_DATA_26 Op = 0x1a
	OP_DATA_27 Op = 0x1b
	OP_DATA_28 Op = 0x1c
	OP_DATA_29 Op = 0x1d
	OP_DATA_30 Op = 0x1e
	OP_DATA_31 Op = 0x1f
	OP_DATA_32 Op = 0x20
	OP_DATA_33 Op = 0x21
	OP_DATA_34 Op = 0x22
	OP_DATA_35 Op = 0x23
	OP_DATA_36 Op = 0x24
	OP_DATA_37 Op = 0x25
	OP_DATA_38 Op = 0x26
	OP_DATA_39 Op = 0x27
	OP_DATA_40 Op = 0x28
	OP_DATA_41 Op = 0x29
	OP_DATA_42 Op = 0x2a
	OP_DATA_43 Op = 0x2b
	OP_DATA_44 Op = 0x2c
	OP_DATA_45 Op = 0x2d
	OP_DATA_46 Op = 0x2e
	OP_DATA_47 Op = 0x2f
	OP_DATA_48 Op = 0x30
	OP_DATA_49 Op = 0x31
	OP_DATA_50 Op = 0x32
	OP_DATA_51 Op = 0x33
	OP_DATA_52 Op = 0x34
	OP_DATA_53 Op = 0x35
	OP_DATA_54 Op = 0x36
	OP_DATA_55 Op = 0x37
	OP_DATA_56 Op = 0x38
	OP_DATA_57 Op = 0x39
	OP_DATA_58 Op = 0x3a
	OP_DATA_59 Op = 0x3b
	OP_DATA_60 Op = 0x3c
	OP_DATA_61 Op = 0x3d
	OP_DATA_62 Op = 0x3e
	OP_DATA_63 Op = 0x3f
	OP_DATA_64 Op = 0x40
	OP_DATA_65 Op = 0x41
	OP_DATA_66 Op = 0x42
	OP_DATA_67 Op = 0x43
	OP_DATA_68 Op = 0x44
	OP_DATA_69 Op = 0x45
	OP_DATA_70 Op = 0x46
	OP_DATA_71 Op = 0x47
	OP_DATA_72 Op = 0x48
	OP_DATA_73 Op = 0x49
	OP_DATA_74 Op = 0x4a
	OP_DATA_75 Op = 0x4b

	OP_PUSHDATA1 Op = 0x4c
	OP_PUSHDATA2 Op = 0x4d
	OP_PUSHDATA4 Op = 0x4e
	OP_1NEGATE   Op = 0x4f
	OP_NOP       Op = 0x61

	OP_JUMP           Op = 0x63
	OP_JUMPIF         Op = 0x64
	OP_VERIFY         Op = 0x69
	OP_FAIL           Op = 0x6a
	OP_CHECKPREDICATE Op = 0xc0

	OP_TOALTSTACK   Op = 0x6b
	OP_FROMALTSTACK Op = 0x6c
	OP_2DROP        Op = 0x6d
	OP_2DUP         Op = 0x6e
	OP_3DUP         Op = 0x6f
	OP_2OVER        Op = 0x70
	OP_2ROT         Op = 0x71
	OP_2SWAP        Op = 0x72
	OP_IFDUP        Op = 0x73
	OP_DEPTH        Op = 0x74
	OP_DROP         Op = 0x75
	OP_DUP          Op = 0x76
	OP_NIP          Op = 0x77
	OP_OVER         Op = 0x78
	OP_PICK         Op = 0x79
	OP_ROLL         Op = 0x7a
	OP_ROT          Op = 0x7b
	OP_SWAP         Op = 0x7c
	OP_TUCK         Op = 0x7d

	OP_CAT         Op = 0x7e
	OP_SUBSTR      Op = 0x7f
	OP_LEFT        Op = 0x80
	OP_RIGHT       Op = 0x81
	OP_SIZE        Op = 0x82
	OP_CATPUSHDATA Op = 0x89

	OP_INVERT      Op = 0x83
	OP_AND         Op = 0x84
	OP_OR          Op = 0x85
	OP_XOR         Op = 0x86
	OP_EQUAL       Op = 0x87
	OP_EQUALVERIFY Op = 0x88

	OP_1ADD               Op = 0x8b
	OP_1SUB               Op = 0x8c
	OP_2MUL               Op = 0x8d
	OP_2DIV               Op = 0x8e
	OP_NEGATE             Op = 0x8f
	OP_ABS                Op = 0x90
	OP_NOT                Op = 0x91
	OP_0NOTEQUAL          Op = 0x92
	OP_ADD                Op = 0x93
	OP_SUB                Op = 0x94
	OP_MUL                Op = 0x95
	OP_DIV                Op = 0x96
	OP_MOD                Op = 0x97
	OP_LSHIFT             Op = 0x98
	OP_RSHIFT             Op = 0x99
	OP_BOOLAND            Op = 0x9a
	OP_BOOLOR             Op = 0x9b
	OP_NUMEQUAL           Op = 0x9c
	OP_NUMEQUALVERIFY     Op = 0x9d
	OP_NUMNOTEQUAL        Op = 0x9e
	OP_LESSTHAN           Op = 0x9f
	OP_GREATERTHAN        Op = 0xa0
	OP_LESSTHANOREQUAL    Op = 0xa1
	OP_GREATERTHANOREQUAL Op = 0xa2
	OP_MIN                Op = 0xa3
	OP_MAX                Op = 0xa4
	OP_WITHIN             Op = 0xa5

	OP_SHA256        Op = 0xa8
	OP_SHA3          Op = 0xaa
	OP_HASH160       Op = 0xab
	OP_CHECKSIG      Op = 0xac
	OP_CHECKMULTISIG Op = 0xad
	OP_CHECKDIGEST   Op = 0xae

	OP_DATASEG Op = 0xf0
	OP_CODESEG Op = 0xf1
)
