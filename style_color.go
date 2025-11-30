package tbl

import (
	"fmt"
	"strconv"
	"strings"
)

// colorMode indicates color format.
type colorMode uint8

const (
	colorModeDefault   colorMode = iota
	colorModeStandard            // 16 colors
	colorMode256                 // 256-color palette
	colorModeTrueColor           // 24-bit RGB
)

// Color represents font (text) color.
type Color struct {
	mode  colorMode
	value uint32
}

// BgColor represents background color.
// Type alias ensures correct field assignment while reusing Color logic.
type BgColor Color

// Style implements Freestyler for Color (applies to FontColor field).
func (c Color) Style(base CellStyle) CellStyle {
	base.FontColor = c
	return base
}

// Style implements Freestyler for BgColor (applies to BgColor field).
func (b BgColor) Style(base CellStyle) CellStyle {
	base.BgColor = b
	return base
}

// Standard font colors (16 functions)

func Black() Color   { return Color{mode: colorModeStandard, value: 30} }
func Red() Color     { return Color{mode: colorModeStandard, value: 31} }
func Green() Color   { return Color{mode: colorModeStandard, value: 32} }
func Yellow() Color  { return Color{mode: colorModeStandard, value: 33} }
func Blue() Color    { return Color{mode: colorModeStandard, value: 34} }
func Magenta() Color { return Color{mode: colorModeStandard, value: 35} }
func Cyan() Color    { return Color{mode: colorModeStandard, value: 36} }
func White() Color   { return Color{mode: colorModeStandard, value: 37} }

func BrightBlack() Color   { return Color{mode: colorModeStandard, value: 90} }
func BrightRed() Color     { return Color{mode: colorModeStandard, value: 91} }
func BrightGreen() Color   { return Color{mode: colorModeStandard, value: 92} }
func BrightYellow() Color  { return Color{mode: colorModeStandard, value: 93} }
func BrightBlue() Color    { return Color{mode: colorModeStandard, value: 94} }
func BrightMagenta() Color { return Color{mode: colorModeStandard, value: 95} }
func BrightCyan() Color    { return Color{mode: colorModeStandard, value: 96} }
func BrightWhite() Color   { return Color{mode: colorModeStandard, value: 97} }

// Standard background colors (16 functions)

func BgBlack() BgColor   { return BgColor(Black()) }
func BgRed() BgColor     { return BgColor(Red()) }
func BgGreen() BgColor   { return BgColor(Green()) }
func BgYellow() BgColor  { return BgColor(Yellow()) }
func BgBlue() BgColor    { return BgColor(Blue()) }
func BgMagenta() BgColor { return BgColor(Magenta()) }
func BgCyan() BgColor    { return BgColor(Cyan()) }
func BgWhite() BgColor   { return BgColor(White()) }

func BgBrightBlack() BgColor   { return BgColor(BrightBlack()) }
func BgBrightRed() BgColor     { return BgColor(BrightRed()) }
func BgBrightGreen() BgColor   { return BgColor(BrightGreen()) }
func BgBrightYellow() BgColor  { return BgColor(BrightYellow()) }
func BgBrightBlue() BgColor    { return BgColor(BrightBlue()) }
func BgBrightMagenta() BgColor { return BgColor(BrightMagenta()) }
func BgBrightCyan() BgColor    { return BgColor(BrightCyan()) }
func BgBrightWhite() BgColor   { return BgColor(BrightWhite()) }

// Color256 creates font color from 256-color palette (0-255).
func Color256(index uint8) Color {
	return Color{mode: colorMode256, value: uint32(index)}
}

// RGB creates font color from RGB values (0-255 each).
func RGB(r, g, b uint8) Color {
	return Color{
		mode:  colorModeTrueColor,
		value: uint32(r)<<16 | uint32(g)<<8 | uint32(b),
	}
}

// Hex creates font color from hex string (#RRGGBB or RRGGBB).
func Hex(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		panic(fmt.Sprintf("tbl: invalid hex color %q (expected #RRGGBB)", hex))
	}

	val, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		panic(fmt.Sprintf("tbl: invalid hex color %q: %v", hex, err))
	}

	return Color{mode: colorModeTrueColor, value: uint32(val)}
}

// BgColor256 creates background color from 256-color palette (0-255).
func BgColor256(index uint8) BgColor {
	return BgColor(Color256(index))
}

// BgRGB creates background color from RGB values (0-255 each).
func BgRGB(r, g, b uint8) BgColor {
	return BgColor(RGB(r, g, b))
}

// BgHex creates background color from hex string (#RRGGBB or RRGGBB).
func BgHex(hex string) BgColor {
	return BgColor(Hex(hex))
}

// ansiCode returns ANSI escape sequence for color.
// isBg: true for background, false for font color.
func ansiCode(c Color, isBg bool) string {
	switch c.mode {
	case colorModeDefault:
		return ""

	case colorModeStandard:
		code := c.value
		if isBg {
			code += 10 // background offset
		}
		return fmt.Sprintf("\x1b[%dm", code)

	case colorMode256:
		prefix := "38" // font color
		if isBg {
			prefix = "48" // background
		}
		return fmt.Sprintf("\x1b[%s;5;%dm", prefix, uint8(c.value))

	case colorModeTrueColor:
		prefix := "38" // font color
		if isBg {
			prefix = "48" // background
		}
		r := uint8(c.value >> 16)
		g := uint8(c.value >> 8)
		b := uint8(c.value)
		return fmt.Sprintf("\x1b[%s;2;%d;%d;%dm", prefix, r, g, b)

	default:
		return ""
	}
}
