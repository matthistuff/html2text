package html2text

import (
	"regexp"
)

var (
	spaceNewLineRe = regexp.MustCompile(` ?\n ?`)
	multiNewLineRe = regexp.MustCompile(`\n\n+`)
	doubleSpaceRe  = regexp.MustCompile(`  +`)
	spacingRe      = regexp.MustCompile(`[ \r\n\t]+`)
)
