package novation

import (
	"go.skymyer.dev/show-control/io"
)

const (
	BTN_1_1 = 0x51
	BTN_1_2 = 0x52
	BTN_1_3 = 0x53
	BTN_1_4 = 0x54
	BTN_1_5 = 0x55
	BTN_1_6 = 0x56
	BTN_1_7 = 0x57
	BTN_1_8 = 0x58

	BTN_2_1 = 0x47
	BTN_2_2 = 0x48
	BTN_2_3 = 0x49
	BTN_2_4 = 0x4A
	BTN_2_5 = 0x4B
	BTN_2_6 = 0x4C
	BTN_2_7 = 0x4D
	BTN_2_8 = 0x4E

	BTN_3_1 = 0x3D
	BTN_3_2 = 0x3E
	BTN_3_3 = 0x3F
	BTN_3_4 = 0x40
	BTN_3_5 = 0x41
	BTN_3_6 = 0x42
	BTN_3_7 = 0x43
	BTN_3_8 = 0x44

	BTN_4_1 = 0x33
	BTN_4_2 = 0x34
	BTN_4_3 = 0x35
	BTN_4_4 = 0x36
	BTN_4_5 = 0x37
	BTN_4_6 = 0x38
	BTN_4_7 = 0x39
	BTN_4_8 = 0x3A

	BTN_5_1 = 0x29
	BTN_5_2 = 0x2A
	BTN_5_3 = 0x2B
	BTN_5_4 = 0x2C
	BTN_5_5 = 0x2D
	BTN_5_6 = 0x2E
	BTN_5_7 = 0x2F
	BTN_5_8 = 0x30

	BTN_6_1 = 0x1F
	BTN_6_2 = 0x20
	BTN_6_3 = 0x21
	BTN_6_4 = 0x22
	BTN_6_5 = 0x23
	BTN_6_6 = 0x24
	BTN_6_7 = 0x25
	BTN_6_8 = 0x26

	BTN_7_1 = 0x15
	BTN_7_2 = 0x16
	BTN_7_3 = 0x17
	BTN_7_4 = 0x18
	BTN_7_5 = 0x19
	BTN_7_6 = 0x1A
	BTN_7_7 = 0x1B
	BTN_7_8 = 0x1C

	BTN_8_1 = 0x0B
	BTN_8_2 = 0x0C
	BTN_8_3 = 0x0D
	BTN_8_4 = 0x0E
	BTN_8_5 = 0x0F
	BTN_8_6 = 0x10
	BTN_8_7 = 0x11
	BTN_8_8 = 0x12

	BTN_ARROW_UP    = 0x5B
	BTN_ARROW_DOWN  = 0x5C
	BTN_ARROW_LEFT  = 0x5D
	BTN_ARROW_RIGHT = 0x5E

	BTN_SESSION = 0x5F
	BTN_DRUMS   = 0x60
	BTN_KEYS    = 0x61
	BTN_USER    = 0x62

	BTN_CALL_1 = 0x59
	BTN_CALL_2 = 0x4F
	BTN_CALL_3 = 0x45
	BTN_CALL_4 = 0x3B
	BTN_CALL_5 = 0x31
	BTN_CALL_6 = 0x27
	BTN_CALL_7 = 0x1D

	BTN_STOP = 0x13
	BTN_LOGO = 0x63
)

var (
	buttonGrid = [][]io.Button{
		{io.BTN_1_1, io.BTN_1_2, io.BTN_1_3, io.BTN_1_4, io.BTN_1_5, io.BTN_1_6, io.BTN_1_7, io.BTN_1_8},
		{io.BTN_2_1, io.BTN_2_2, io.BTN_2_3, io.BTN_2_4, io.BTN_2_5, io.BTN_2_6, io.BTN_2_7, io.BTN_2_8},
		{io.BTN_3_1, io.BTN_3_2, io.BTN_3_3, io.BTN_3_4, io.BTN_3_5, io.BTN_3_6, io.BTN_3_7, io.BTN_3_8},
		{io.BTN_4_1, io.BTN_4_2, io.BTN_4_3, io.BTN_4_4, io.BTN_4_5, io.BTN_4_6, io.BTN_4_7, io.BTN_4_8},
		{io.BTN_5_1, io.BTN_5_2, io.BTN_5_3, io.BTN_5_4, io.BTN_5_5, io.BTN_5_6, io.BTN_5_7, io.BTN_5_8},
		{io.BTN_6_1, io.BTN_6_2, io.BTN_6_3, io.BTN_6_4, io.BTN_6_5, io.BTN_6_6, io.BTN_6_7, io.BTN_6_8},
		{io.BTN_7_1, io.BTN_7_2, io.BTN_7_3, io.BTN_7_4, io.BTN_7_5, io.BTN_7_6, io.BTN_7_7, io.BTN_7_8},
		{io.BTN_8_1, io.BTN_8_2, io.BTN_8_3, io.BTN_8_4, io.BTN_8_5, io.BTN_8_6, io.BTN_8_7, io.BTN_8_8},
	}

	controlMapReverse = map[io.Control]byte{}
	controlMap        = map[byte]io.Control{
		BTN_STOP: io.CTR_BACK,
		BTN_LOGO: io.CTR_SELECT,

		BTN_SESSION: io.CTR_PAGE_1,
		BTN_DRUMS:   io.CTR_PAGE_2,
		BTN_KEYS:    io.CTR_PAGE_3,
		BTN_USER:    io.CTR_PAGE_4,
	}

	buttonMapReverse = map[io.Button]byte{}
	buttonMap        = map[byte]io.Button{
		BTN_ARROW_UP:    io.BTN_ARROW_UP,
		BTN_ARROW_DOWN:  io.BTN_ARROW_DOWN,
		BTN_ARROW_LEFT:  io.BTN_ARROW_LEFT,
		BTN_ARROW_RIGHT: io.BTN_ARROW_RIGHT,

		BTN_CALL_1: io.BTN_CALL_1,
		BTN_CALL_2: io.BTN_CALL_2,
		BTN_CALL_3: io.BTN_CALL_3,
		BTN_CALL_4: io.BTN_CALL_4,
		BTN_CALL_5: io.BTN_CALL_5,
		BTN_CALL_6: io.BTN_CALL_6,
		BTN_CALL_7: io.BTN_CALL_7,

		BTN_1_1: io.BTN_1_1,
		BTN_1_2: io.BTN_1_2,
		BTN_1_3: io.BTN_1_3,
		BTN_1_4: io.BTN_1_4,
		BTN_1_5: io.BTN_1_5,
		BTN_1_6: io.BTN_1_6,
		BTN_1_7: io.BTN_1_7,
		BTN_1_8: io.BTN_1_8,

		BTN_2_1: io.BTN_2_1,
		BTN_2_2: io.BTN_2_2,
		BTN_2_3: io.BTN_2_3,
		BTN_2_4: io.BTN_2_4,
		BTN_2_5: io.BTN_2_5,
		BTN_2_6: io.BTN_2_6,
		BTN_2_7: io.BTN_2_7,
		BTN_2_8: io.BTN_2_8,

		BTN_3_1: io.BTN_3_1,
		BTN_3_2: io.BTN_3_2,
		BTN_3_3: io.BTN_3_3,
		BTN_3_4: io.BTN_3_4,
		BTN_3_5: io.BTN_3_5,
		BTN_3_6: io.BTN_3_6,
		BTN_3_7: io.BTN_3_7,
		BTN_3_8: io.BTN_3_8,

		BTN_4_1: io.BTN_4_1,
		BTN_4_2: io.BTN_4_2,
		BTN_4_3: io.BTN_4_3,
		BTN_4_4: io.BTN_4_4,
		BTN_4_5: io.BTN_4_5,
		BTN_4_6: io.BTN_4_6,
		BTN_4_7: io.BTN_4_7,
		BTN_4_8: io.BTN_4_8,

		BTN_5_1: io.BTN_5_1,
		BTN_5_2: io.BTN_5_2,
		BTN_5_3: io.BTN_5_3,
		BTN_5_4: io.BTN_5_4,
		BTN_5_5: io.BTN_5_5,
		BTN_5_6: io.BTN_5_6,
		BTN_5_7: io.BTN_5_7,
		BTN_5_8: io.BTN_5_8,

		BTN_6_1: io.BTN_6_1,
		BTN_6_2: io.BTN_6_2,
		BTN_6_3: io.BTN_6_3,
		BTN_6_4: io.BTN_6_4,
		BTN_6_5: io.BTN_6_5,
		BTN_6_6: io.BTN_6_6,
		BTN_6_7: io.BTN_6_7,
		BTN_6_8: io.BTN_6_8,

		BTN_7_1: io.BTN_7_1,
		BTN_7_2: io.BTN_7_2,
		BTN_7_3: io.BTN_7_3,
		BTN_7_4: io.BTN_7_4,
		BTN_7_5: io.BTN_7_5,
		BTN_7_6: io.BTN_7_6,
		BTN_7_7: io.BTN_7_7,
		BTN_7_8: io.BTN_7_8,

		BTN_8_1: io.BTN_8_1,
		BTN_8_2: io.BTN_8_2,
		BTN_8_3: io.BTN_8_3,
		BTN_8_4: io.BTN_8_4,
		BTN_8_5: io.BTN_8_5,
		BTN_8_6: io.BTN_8_6,
		BTN_8_7: io.BTN_8_7,
		BTN_8_8: io.BTN_8_8,
	}
)
