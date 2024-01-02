package webhook

import (
	"go.skymyer.dev/show-control/io"
)

var controlMap = map[string]io.Control{
	"CTR_BACK":   io.CTR_BACK,
	"CTR_PAGE_1": io.CTR_PAGE_1,
	"CTR_PAGE_2": io.CTR_PAGE_2,
	"CTR_PAGE_3": io.CTR_PAGE_3,
	"CTR_PAGE_4": io.CTR_PAGE_4,
}

var buttonMap = map[string]io.Button{
	"BTN_ARROW_UP":    io.BTN_ARROW_UP,
	"BTN_ARROW_DOWN":  io.BTN_ARROW_DOWN,
	"BTN_ARROW_LEFT":  io.BTN_ARROW_LEFT,
	"BTN_ARROW_RIGHT": io.BTN_ARROW_RIGHT,

	"BTN_CALL_1": io.BTN_CALL_1,
	"BTN_CALL_2": io.BTN_CALL_2,
	"BTN_CALL_3": io.BTN_CALL_3,
	"BTN_CALL_4": io.BTN_CALL_4,
	"BTN_CALL_5": io.BTN_CALL_5,
	"BTN_CALL_6": io.BTN_CALL_6,
	"BTN_CALL_7": io.BTN_CALL_7,

	"BTN_1_1": io.BTN_1_1,
	"BTN_1_2": io.BTN_1_2,
	"BTN_1_3": io.BTN_1_3,
	"BTN_1_4": io.BTN_1_4,
	"BTN_1_5": io.BTN_1_5,
	"BTN_1_6": io.BTN_1_6,
	"BTN_1_7": io.BTN_1_7,
	"BTN_1_8": io.BTN_1_8,

	"BTN_2_1": io.BTN_2_1,
	"BTN_2_2": io.BTN_2_2,
	"BTN_2_3": io.BTN_2_3,
	"BTN_2_4": io.BTN_2_4,
	"BTN_2_5": io.BTN_2_5,
	"BTN_2_6": io.BTN_2_6,
	"BTN_2_7": io.BTN_2_7,
	"BTN_2_8": io.BTN_2_8,

	"BTN_3_1": io.BTN_3_1,
	"BTN_3_2": io.BTN_3_2,
	"BTN_3_3": io.BTN_3_3,
	"BTN_3_4": io.BTN_3_4,
	"BTN_3_5": io.BTN_3_5,
	"BTN_3_6": io.BTN_3_6,
	"BTN_3_7": io.BTN_3_7,
	"BTN_3_8": io.BTN_3_8,

	"BTN_4_1": io.BTN_4_1,
	"BTN_4_2": io.BTN_4_2,
	"BTN_4_3": io.BTN_4_3,
	"BTN_4_4": io.BTN_4_4,
	"BTN_4_5": io.BTN_4_5,
	"BTN_4_6": io.BTN_4_6,
	"BTN_4_7": io.BTN_4_7,
	"BTN_4_8": io.BTN_4_8,

	"BTN_5_1": io.BTN_5_1,
	"BTN_5_2": io.BTN_5_2,
	"BTN_5_3": io.BTN_5_3,
	"BTN_5_4": io.BTN_5_4,
	"BTN_5_5": io.BTN_5_5,
	"BTN_5_6": io.BTN_5_6,
	"BTN_5_7": io.BTN_5_7,
	"BTN_5_8": io.BTN_5_8,

	"BTN_6_1": io.BTN_6_1,
	"BTN_6_2": io.BTN_6_2,
	"BTN_6_3": io.BTN_6_3,
	"BTN_6_4": io.BTN_6_4,
	"BTN_6_5": io.BTN_6_5,
	"BTN_6_6": io.BTN_6_6,
	"BTN_6_7": io.BTN_6_7,
	"BTN_6_8": io.BTN_6_8,

	"BTN_7_1": io.BTN_7_1,
	"BTN_7_2": io.BTN_7_2,
	"BTN_7_3": io.BTN_7_3,
	"BTN_7_4": io.BTN_7_4,
	"BTN_7_5": io.BTN_7_5,
	"BTN_7_6": io.BTN_7_6,
	"BTN_7_7": io.BTN_7_7,
	"BTN_7_8": io.BTN_7_8,

	"BTN_8_1": io.BTN_8_1,
	"BTN_8_2": io.BTN_8_2,
	"BTN_8_3": io.BTN_8_3,
	"BTN_8_4": io.BTN_8_4,
	"BTN_8_5": io.BTN_8_5,
	"BTN_8_6": io.BTN_8_6,
	"BTN_8_7": io.BTN_8_7,
	"BTN_8_8": io.BTN_8_8,
}
