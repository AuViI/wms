package model

import (
	"fmt"
	"html/template"
	"regexp"
	"strconv"
)

type Theme struct {
	StartColor ThemeColor
	EndColor   ThemeColor
	IconLink   string

	ToHTMLColor func(ThemeColor) string
	ToHTMLIcon  func(string) template.HTML
}

type TemplateTheme struct {
	StartColor string
	EndColor   string
	IconLink   template.HTML
}

type ThemeColor struct {
	R int64
	G int64
	B int64
}

var reHexColor = regexp.MustCompile("#([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})")

func ThemeColorFromHex(hexdesc string) ThemeColor {
	match := reHexColor.FindStringSubmatch(hexdesc)

	var color ThemeColor
	color.R, _ = strconv.ParseInt(match[1], 16, 64)
	color.G, _ = strconv.ParseInt(match[2], 16, 64)
	color.B, _ = strconv.ParseInt(match[3], 16, 64)
	return color
}

func (t Theme) Prepare() TemplateTheme {
	return TemplateTheme{
		StartColor: ToHTMLColor(t.StartColor),
		EndColor:   ToHTMLColor(t.EndColor),
		IconLink:   ToHTMLIcon(t.IconLink),
	}
}

var defaultTheme = Theme{
	StartColor: ThemeColorFromHex("#82a0d5"),
	EndColor:   ThemeColorFromHex("#4b6ca6"),
	IconLink:   "/resources/logo.png",

	ToHTMLIcon:  ToHTMLIcon,
	ToHTMLColor: ToHTMLColor,
}

func ToHTMLColor(c ThemeColor) string {
	return fmt.Sprintf("#%x%x%x", c.R, c.G, c.B)
}

func ToHTMLIcon(iconLink string) template.HTML {
	return template.HTML(
		fmt.Sprintf("<img id='themeicon' src='%s'></img>", iconLink),
	)
}

func GetDefaultTheme() Theme {
	return defaultTheme
}
