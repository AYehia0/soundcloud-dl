package theme

import (
	"fmt"

	"github.com/fatih/color"
)

var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var Magenta = color.New(color.FgMagenta).SprintFunc()
var Yellow = color.New(color.FgHiYellow).SprintFunc()

func FormatTime(n int64) string {
	n /= 1000
	sec := n % 60
	n /= 60
	min := n % 60
	n = n / 60
	return fmt.Sprintf("%02d:%02d:%02d", n, min, sec)
}
