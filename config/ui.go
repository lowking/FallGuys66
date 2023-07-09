package config

import (
	"fyne.io/fyne/v2"
	"image/color"
)

var AppName = "糖豆人直播助手：大周定制版"
var AppSize = fyne.NewSize(1005, 860)
var LogoSize = float32(90)
var Padding = float32(10)
var ToolbarPaddingTop = Padding + 10
var ToolbarPaddingLeft = float32(120)
var RemarkText = `### PS, 💖 Power with love.`

var AccentColor = color.RGBA{
	R: 46,
	G: 108,
	B: 246,
	A: 255,
}
var ShadowColor = color.RGBA{
	R: 66,
	G: 66,
	B: 66,
	A: 255,
}
var VersionColor = color.RGBA{
	R: 43,
	G: 87,
	B: 188,
	A: 255,
}
