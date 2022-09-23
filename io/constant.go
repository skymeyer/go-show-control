package io

type Mode int

const (
	MODE_MAIN Mode = iota
	MODE_1
	MODE_2
	MODE_3
	MODE_4
)

type Page int

const (
	PAGE_0 Page = iota
	PAGE_1
	PAGE_2
	PAGE_3
	PAGE_4
)

type Action int

const (
	_ Action = iota
	ACTION_PRESS
	ACTION_RELEASE
)

type Control int

const (
	CTR_SELECT Control = iota
	CTR_PAGE_1
	CTR_PAGE_2
	CTR_PAGE_3
	CTR_PAGE_4
	CTR_BACK
)

type ButtonKind int

const (
	BTN_KIND_OFF ButtonKind = iota
	BTN_KIND_STOP

	BTN_KIND_DEFAULT
	BTN_KIND_DEFAULT_SELECT

	BTN_KIND_CALL
	BTN_KIND_CALL_SELECT
	BTN_KIND_ARROW
	BTN_KIND_ARROW_SELECT

	BTN_KIND_COLOR_BLUE
	BTN_KIND_COLOR_CYAN
	BTN_KIND_COLOR_GRAY
	BTN_KIND_COLOR_GREEN
	BTN_KIND_COLOR_MAGENTA
	BTN_KIND_COLOR_ORANGE
	BTN_KIND_COLOR_RED
	BTN_KIND_COLOR_WHITE
	BTN_KIND_COLOR_YELLOW

	BTN_KIND_DUAL_MAGENTA_CYAN
	BTN_KIND_DUAL_MAGENTA_BLUE
	BTN_KIND_DUAL_GREEN_BLUE
	BTN_KIND_DUAL_GREEN_YELLOW
	BTN_KIND_DUAL_ORANGE_YELLOW
	BTN_KIND_DUAL_ORANGE_RED

	BTN_KIND_WIDGET_SLIDER_RED
	BTN_KIND_WIDGET_SLIDER_GREEN
	BTN_KIND_WIDGET_SLIDER_BLUE
	BTN_KIND_WIDGET_SLIDER_WHITE
)

type Button int

const (
	_ Button = iota

	BTN_ARROW_UP
	BTN_ARROW_DOWN
	BTN_ARROW_LEFT
	BTN_ARROW_RIGHT

	BTN_CALL_1
	BTN_CALL_2
	BTN_CALL_3
	BTN_CALL_4
	BTN_CALL_5
	BTN_CALL_6
	BTN_CALL_7

	BTN_1_1
	BTN_1_2
	BTN_1_3
	BTN_1_4
	BTN_1_5
	BTN_1_6
	BTN_1_7
	BTN_1_8

	BTN_2_1
	BTN_2_2
	BTN_2_3
	BTN_2_4
	BTN_2_5
	BTN_2_6
	BTN_2_7
	BTN_2_8

	BTN_3_1
	BTN_3_2
	BTN_3_3
	BTN_3_4
	BTN_3_5
	BTN_3_6
	BTN_3_7
	BTN_3_8

	BTN_4_1
	BTN_4_2
	BTN_4_3
	BTN_4_4
	BTN_4_5
	BTN_4_6
	BTN_4_7
	BTN_4_8

	BTN_5_1
	BTN_5_2
	BTN_5_3
	BTN_5_4
	BTN_5_5
	BTN_5_6
	BTN_5_7
	BTN_5_8

	BTN_6_1
	BTN_6_2
	BTN_6_3
	BTN_6_4
	BTN_6_5
	BTN_6_6
	BTN_6_7
	BTN_6_8

	BTN_7_1
	BTN_7_2
	BTN_7_3
	BTN_7_4
	BTN_7_5
	BTN_7_6
	BTN_7_7
	BTN_7_8

	BTN_8_1
	BTN_8_2
	BTN_8_3
	BTN_8_4
	BTN_8_5
	BTN_8_6
	BTN_8_7
	BTN_8_8
)

var PageToControl = map[Page]Control{
	PAGE_1: CTR_PAGE_1,
	PAGE_2: CTR_PAGE_2,
	PAGE_3: CTR_PAGE_3,
	PAGE_4: CTR_PAGE_4,
}
