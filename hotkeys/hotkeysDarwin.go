//go:build darwin && cgo

package hotkeys

import (
	"golang.design/x/hotkey"
	"strings"
)

func GetModifier(key string) hotkey.Modifier {
	switch strings.ToLower(key) {
	case "cmd":
		return hotkey.Modifier(hotkey.ModCmd)
	case "ctrl":
		return hotkey.Modifier(hotkey.ModCtrl)
	case "option":
		return hotkey.Modifier(hotkey.ModOption)
	default:
		return hotkey.Modifier(hotkey.ModShift)
	}
}

func GetModifierName(key hotkey.Modifier) string {
	switch key {
	case hotkey.Modifier(hotkey.ModCmd):
		return "cmd"
	case hotkey.Modifier(hotkey.ModCtrl):
		return "ctrl"
	case hotkey.Modifier(hotkey.ModOption):
		return "option"
	default:
		return "shift"
	}
}
