// The refmt package contains extra string utility functions
package refmt

import (
	"fmt"
	"regexp"
	"strings"
)

// The Wrap object may be used to indent and/or wrap/rewrap strings.
type Style struct {
	IndentWidth int  // Number of indent characters to prepend
	IndentChar  rune // Character to fill indent space with
	MaxWidth    int  // Maximum width of line (not counting indent)
	SplitChar   rune // usually a newline
	WordSep     rune // usually a space

	// when a word exceeds maxWidth, do we allow a long line to
	// exceed maxwidth or split the word?
	BreakLongWords bool
}

// NewWrap returns a new Wrap object preconfigured for a 4 space indent and 76
// character lines.
func NewStyle() *Style {
	// return a new instance of Wrap with sane defaults
	return &Style{4, ' ', 76, '\n', ' ', false}
}

func (s Style) Wrap(input string) string {
	// trim the string
	output := strings.Trim(input, string(s.WordSep)+string(s.SplitChar))
	// replace single line splits with word separators. This unwraps lines.
	collapseNL := regexp.MustCompile(fmt.Sprintf("([^%[1]s])%[1]s([^%[1]s])", string(s.SplitChar)))
	output = collapseNL.ReplaceAllString(output, fmt.Sprintf("$1%s$2", string(s.WordSep)))

	seek := 0
	end := s.MaxWidth

	// Closure for wrapping line
	split := func(loc int) {
		output = output[:loc] + string(s.SplitChar) + output[loc+1:]
		seek = loc + 1
		end = seek + s.MaxWidth
	}

	for end < len(output) {
		// if there's a newline in our range, advance the range beyond the newline
		if e := strings.Index(output[seek:end], string(s.SplitChar)); e >= 0 {
			seek += e + 1
			end = seek + s.MaxWidth
		}
		i := strings.LastIndex(output[seek:end], string(s.WordSep))
		if i >= 0 { // normal case
			split(i + seek)
		} else { // case when word exceeds line width
			if s.BreakLongWords { // break the word ,or
				split(end)
			} else { // allow line to exceed maxWidth
				if i := strings.Index(output[seek:], string(s.WordSep)); i >= 0 {
					split(i + seek)
				} else {
					break
				}
			}
		}
	}
	return output
}

// Indent returns an indented version of input
func (s Style) Indent(input string) string {
	indentStr := strings.Repeat(string(s.IndentChar), s.IndentWidth)
	// Insert an indent at the beginning of the input:
	output := indentStr + input
	// Replace all SplitChar with SplitChar + indentStr, except for the last
	// char (no need to indent at the end)
	indentRe := regexp.MustCompile(fmt.Sprintf("%[1]s", string(s.SplitChar)))
	output = indentRe.ReplaceAllString(output[:len(output)-1], string(s.SplitChar)+indentStr) + output[len(output)-1:]
	return output
}

///////////////////

// Underline returns a string for an underlined heading
func Underline(ulinechar string, format string, v ...interface{}) string {
	result := fmt.Sprintf(format, v...)
	result = fmt.Sprintf("%s\n%s", result, strings.Repeat(string(ulinechar[0]), len(result)))
	return result
}
