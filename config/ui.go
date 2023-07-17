package config

import (
	"fyne.io/fyne/v2"
	"image/color"
	"os"
	"path/filepath"
)

var AppName = "糖豆人直播助手：大周定制版"
var AppSize = fyne.NewSize(1105, 860)
var LogoSize = float32(90)
var Padding = float32(10)
var ToolbarPaddingTop = Padding + 10
var ToolbarPaddingLeft = float32(120)
var RemarkText = `### PS, 💖 Power with love.`
var UserConfigDir string

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	configDir = filepath.Join(configDir, "fallguys66")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		panic(err)
	}
	UserConfigDir = configDir
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
