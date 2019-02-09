package model

import (
	"fmt"
	"html/template"
	"regexp"
	"strconv"
)

type Theme struct {
	StartColor     ThemeColor `json:"start" yaml:"start,flow"`
	EndColor       ThemeColor `json:"end" yaml:"end,flow"`
	IconLink       string     `json:"ilink" yaml:"ilink"`
	IconSymLink    string     `json:"slink" yaml:"slink"`
	IconVanityLink string     `json:"vlink" yaml:"vlink"`

	ToHTMLColor     func(ThemeColor) string            `json:"-" yaml:"-"`
	ToHTMLIcon      func(string) template.HTML         `json:"-" yaml:"-"`
	ToHTMLSmartIcon func(string, string) template.HTML `json:"-" yaml:"-"`
}

type TemplateTheme struct {
	StartColor     string
	EndColor       string
	IconLink       template.HTML
	IconSymLink    template.HTML
	IconVanityLink template.HTML
}

type ThemeColor struct {
	R int64 `json:"r" yaml:"r"`
	G int64 `json:"g" yaml:"g"`
	B int64 `json:"b" yaml:"b"`
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
		StartColor:     ToHTMLColor(t.StartColor),
		EndColor:       ToHTMLColor(t.EndColor),
		IconLink:       ToHTMLSmartIcon(t.IconLink, "themeicon"),
		IconSymLink:    ToHTMLSmartIcon(t.IconSymLink, "themeicon-sym"),
		IconVanityLink: ToHTMLSmartIcon(t.IconVanityLink, "icon"),
	}
}

var defaultTheme = Theme{
	StartColor:     ThemeColorFromHex("#42aac9"),
	EndColor:       ThemeColorFromHex("#12657e"),
	IconLink:       "/resources/logo.png",
	IconVanityLink: "/resources/logo_invert.png",

	ToHTMLIcon:      ToHTMLIcon,
	ToHTMLSmartIcon: ToHTMLSmartIcon,
	ToHTMLColor:     ToHTMLColor,
}

func ToHTMLColor(c ThemeColor) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func ToHTMLIcon(iconLink string) template.HTML {
	return template.HTML(
		fmt.Sprintf("<img id='themeicon' src='%s'></img>", iconLink),
	)
}

func ToHTMLSmartIcon(iconLink, id string) template.HTML {
	return template.HTML(
		fmt.Sprintf("<img id='%s' src='%s'></img>", id, iconLink),
	)
}

func GetDefaultTheme() Theme {
	return defaultTheme
}
