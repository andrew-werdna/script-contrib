package filters

import (
	"fmt"
	"io"
	"strings"
)

// ScriptFilter is the type of function required to use `script.Pipe`'s Filter
// method
type ScriptFilter func(r io.Reader, w io.Writer) error

// ScriptFilterScan is the type of function required to use `script.Pipe`'s FilterScan
// method
type ScriptFilterScan func(s string, w io.Writer)

// CutFirstN implements the equivalent of `cut -c1-N`
func CutFirstN(n int64) ScriptFilter {
	return func(r io.Reader, w io.Writer) error {
		_, err := io.CopyN(w, r, n)
		return err
	}
}

// Columns is a package function that returns a function to be used with `Pipe.FilterScan`.
// Provided fn a function that determines what runes to split on, and a slice of int,
// it will return multiple, newDelim-separated columns. There is an open issue
// in [script](https://github.com/bitfield/script/issues/127) that discusses the
// desire for this kind of functionality. This is sort of the equivalent of
// `cut -d {result of fn()} -f {cols}` or `awk -F '\t' '{ print $5 "|" $6 }'`
func Columns(fn func(rune) bool, newDelim string, cols ...int) ScriptFilterScan {
	return func(line string, w io.Writer) {
		if len(cols) < 1 {
			return
		}

		var b strings.Builder
		columns := strings.FieldsFunc(line, fn)
		hasWritten := false
		maxSafeIndex := len(columns) - 1
		finalDesiredColumn := len(cols) - 1

		for i, v := range cols {
			isSafeToUseColumn := v <= maxSafeIndex && v > 0
			isFinalDesiredColumn := i == finalDesiredColumn

			skipColumn := !isSafeToUseColumn && !isFinalDesiredColumn
			writeLineAndReturn := !isSafeToUseColumn && isFinalDesiredColumn && hasWritten
			justReturn := !isSafeToUseColumn && isFinalDesiredColumn && !hasWritten
			writeWithDelim := isSafeToUseColumn && !isFinalDesiredColumn
			writeWithNewline := isSafeToUseColumn && isFinalDesiredColumn

			if skipColumn {
				continue
			} else if justReturn {
				return
			} else if writeLineAndReturn {
				_, _ = fmt.Fprintf(w, "%s\n", strings.TrimRight(b.String(), newDelim))
				return
			} else if writeWithDelim {
				b.WriteString(fmt.Sprintf("%s%s", columns[v-1], newDelim))
				hasWritten = true
			} else if writeWithNewline {
				b.WriteString(fmt.Sprintf("%s\n", columns[v-1]))
			}
		}
		_, _ = fmt.Fprint(w, strings.TrimRight(b.String(), newDelim))
	}
}
