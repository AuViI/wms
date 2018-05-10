package model

import "regexp"

var ThemeRegex = regexp.MustCompile("(.*)&([0-9a-f]{6})&([0-9a-f]{6})&(.*)")
