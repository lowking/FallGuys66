//go:build windows

package hotkeys

import (
	"golang.design/x/hotkey"
	"strings"
)

func GetModifier(key string) hotkey.Modifier {
	switch strings.ToLower(key) {
	case "win":
		return hotkey.Modifier(hotkey.ModWin)
	case "ctrl":
		return hotkey.Modifier(hotkey.ModCtrl)
	case "alt":
		return hotkey.Modifier(hotkey.ModAlt)
	default:
		return hotkey.Modifier(hotkey.ModShift)
	}
}

func GetModifierName(key hotkey.Modifier) string {
	switch key {
	case hotkey.Modifier(hotkey.ModWin):
		return "win"
	case hotkey.Modifier(hotkey.ModCtrl):
		return "ctrl"
	case hotkey.Modifier(hotkey.ModAlt):
		return "alt"
	default:
		return "shift"
	}
}
