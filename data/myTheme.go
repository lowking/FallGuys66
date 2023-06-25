package data

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"strings"
)

type MyTheme struct {
	Regular, Bold, Italic, BoldItalic, Monospace fyne.Resource
}

func (t *MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return m.Monospace
	}
	if style.Bold {
		if style.Italic {
			return m.BoldItalic
		}
		return m.Bold
	}
	if style.Italic {
		return m.Italic
	}
	return m.Regular
}

func (m *MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (t *MyTheme) SetFonts(regularFontPath string, monoFontPath string) {
	t.Regular = theme.TextFont()
	t.Bold = theme.TextBoldFont()
	t.Italic = theme.TextItalicFont()
	t.BoldItalic = theme.TextBoldItalicFont()
	t.Monospace = theme.TextMonospaceFont()

	if regularFontPath != "" {
		t.Regular = loadCustomFont(regularFontPath, "Regular", t.Regular)
		t.Bold = loadCustomFont(regularFontPath, "Bold", t.Bold)
		t.Italic = loadCustomFont(regularFontPath, "Italic", t.Italic)
		t.BoldItalic = loadCustomFont(regularFontPath, "BoldItalic", t.BoldItalic)
	}
	if monoFontPath != "" {
		t.Monospace = loadCustomFont(monoFontPath, "Regular", t.Monospace)
	} else {
		t.Monospace = t.Regular
	}
}

func loadCustomFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	variantPath := strings.Replace(env, "Regular", variant, -1)

	res, err := fyne.LoadResourceFromPath(variantPath)
	if err != nil {
		fyne.LogError("Error loading specified font", err)
		return fallback
	}

	return res
}
